// Package form is used to generate HTML forms. Most notably, the
// form.Builder makes it very easy to take a struct with set values and
// generate the input tags, labels, etc for each field in the struct.
//
// See the examples directory for a more comprehensive idea of what can be\
// accomplished with this package.
package form

import (
	"html/template"
	"strings"
)

// Builder is used to build HTML forms/inputs for Go structs. Basic
// usage looks something like this:
//
//   tpl := template.Must(template.New("").Parse(`
//   	 <input type="{{.Type}}" name="{{.Name}}"" {{with .Value}}value="{{.}}"{{end}}>
//   `))
//   fb := Builder{InputTemplate: tpl}
//   html := fb.Inputs(struct{
//	   Name string `form:"name=full-name"`
//     Email string `form:"type=email"`
//   }{"Michael Scott", "michael@dundermifflin.com"})
//   // Outputs:
//   //   <input type="text" name="full-name" value="Michael Scott">
//   //   <input type="email" name="Email" value="michael@dundermifflin.com">
//
// This is a VERY limited example, but should demonstrate the basic
// idea. The Builder uses a single template to and will call it with all the
// information about each individual field and return the resulting HTML.
//
// The most common use for this is to provide a helper function for your HTML
// templates. Eg something like:
//
//   fb := Builder{...}
//   tpl, err := template.New("").Funcs(template.FuncMap{
//     "inputs_for": fb.Inputs,
//   })
//
//   // Then later in a template:
//   <form>
//     {{inputs_for .SomeStruct}}
//   </form>
//
// For a much more thorough set of examples, see the examples directory.
// There is even an example illustrating how the gorilla/schema package can
// be used to parse forms that are created by the Builder.
type Builder struct {
	InputTemplate *template.Template
}

// Inputs will parse the provided struct into fields and then execute the
// Builder.InputTemplate with each field. The returned HTML is simply
// all of these results appended one after another.
//
// Inputs only accepts structs for the first argument. This may change
// later, but that is all I needed for my use case so it is what it does.
// If you need support for something else like maps let me know.
//
// Inputs' second argument - errs - will be used to render errors for
// individual fields. This is done by looking for errors that implement
// the fieldError interface:
//
//   type fieldError interface {
//    	FieldError() (field, err string)
//   }
//
// Where the first return value is expected to be the field with an error
// and the second is the actual error message to be displayed. This is then
// used to provide an `errors` template function that will return a slice
// of errors (if there are any) for the current field your InputTemplate
// is rendering. See examples/errors/errors.go for an example of this in
// action.
//
// This interface is not exported and you can pass other errors into Inputs
// but they currently won't be used.
func (b *Builder) Inputs(v interface{}, errs ...error) (template.HTML, error) {
	tpl, err := b.InputTemplate.Clone()
	if err != nil {
		return "", err
	}
	fields := fields(v)
	errors := errors(errs)
	var html template.HTML
	for _, field := range fields {
		var sb strings.Builder
		tpl.Funcs(template.FuncMap{
			"errors": func() []string {
				if errs, ok := errors[field.Name]; ok {
					return errs
				}
				return nil
			},
		})
		err := tpl.Execute(&sb, field)
		if err != nil {
			return "", err
		}
		html = html + template.HTML(sb.String())
	}
	return html, nil
}

// FuncMap returns a template.FuncMap that defines both the inputs_for and
// inputs_and_errors_for functions for usage in the template package. The
// latter is provided via a closure because variadic parameters and the
// template package don't play very nicely and this just simplifies things
// a lot for end users of the form package.
func (b *Builder) FuncMap() template.FuncMap {
	return template.FuncMap{
		"inputs_for": b.Inputs,
		"inputs_and_errors_for": func(v interface{}, errs []error) (template.HTML, error) {
			return b.Inputs(v, errs...)
		},
	}
}

// FuncMap is present to make it a little easier to build the InputTemplate
// field of the Builder type. In order to parse a template that uses the
// `errors` function, you need to have that template defined when the
// template is parsed. We clearly don't know whether a field has an error
// or not until it is parsed via the Inputs method call, so this basically
// just provides a stubbed out errors function that returns nil so the template
// compiles correctly.
//
// See examples/errors/errors.go for a clear example of this being used.
func FuncMap() template.FuncMap {
	return template.FuncMap{
		"errors": ErrorsStub,
	}
}

// ErrorsStub is a stubbed out function that simply returns nil. It is present
// to make it a little easier to build the InputTemplate field of the Builder
// type, since your template will likely use the errors function in the
// template before it can truly be defined. You probably just want to use
// the provided FuncMap helper, but this can be useful when you need to
// build your own template.FuncMap.
//
// See examples/errors/errors.go for a clear example of the FuncMap function
// being used, and see FuncMap for an example of how ErrorsStub can be used.
func ErrorsStub() []string {
	return nil
}

// fieldError is an interface defining an error that represents something
// wrong with a particular struct field. The name should correspond to the
// name value used when building the HTML form, which is currently a period
// separated list of all fields that lead up to the particular field.
// Eg, in the following struct the Mouse field would have a key of Cat.Mouse:
//
//   type Dog struct {
//     Cat: struct{
//       Mouse string
//     }
//   }
//
// The top level Dog struct name is not used because this is unnecessary,
// but any other nested struct names are necessary to properly determine
// the field.
//
// It should also be noted that if you provide a custom field name, that
// name should also be used in fieldError implementations.
type fieldError interface {
	FieldError() (field, err string)
}

// errors will build a map where each key is the field name, and each
// value is a slice of strings representing errors with that field.
//
// It works by looking for errors that implement the following interface:
//
//   interface {
//     FieldError() (string, string)
//   }
//
// Where the first string returned is expected to be the field name, and
// the second return value is expected to be an error with that field.
// Any errors that implement this interface are then used to build the
// slice of errors for the field, meaning you can provide multiple
// errors for the same field and all will be utilized.
func errors(errs []error) map[string][]string {
	ret := make(map[string][]string)
	for _, err := range errs {
		fe, ok := err.(fieldError)
		if !ok {
			continue
		}
		field, fieldErr := fe.FieldError()
		ret[field] = append(ret[field], fieldErr)
	}
	return ret
}
