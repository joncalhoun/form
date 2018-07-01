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
// Inputs only accepts structs. This may change later, but that is all I
// needed for my use case so it is what it does.
func (b *Builder) Inputs(v interface{}) template.HTML {
	fields := fields(v)
	var html template.HTML
	for _, field := range fields {
		var sb strings.Builder
		b.InputTemplate.Execute(&sb, field)
		html = html + template.HTML(sb.String())
	}
	return html
}
