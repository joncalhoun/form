package form

import (
	"errors"
	"fmt"
	"html/template"
	"reflect"
	"strings"
	"testing"
)

func TestBuilder_Inputs(t *testing.T) {
	tpl := template.Must(template.New("").Parse(strings.TrimSpace(`
		<label>{{.Label}}</label><input type="{{.Type}}" name="{{.Name}}" placeholder="{{.Placeholder}}"{{with .Value}} value="{{.}}"{{end}}{{with .Class}} class="{{.}}"{{end}}>
	`)))
	tests := []struct {
		name string
		tpl  *template.Template
		arg  interface{}
		want template.HTML
	}{
		{
			name: "label and input",
			tpl:  tpl,
			arg: struct {
				Name  string
				Email string `form:"type=email;placeholder=bob@example.com"`
				City  string `form:"class=custom-city-class"`
			}{
				Name: "Michael Scott",
				City: "New York",
			},
			want: template.HTML(strings.Join([]string{
				strings.TrimSpace(`
					<label>Name</label><input type="text" name="Name" placeholder="Name" value="Michael Scott">`),
				strings.TrimSpace(`
					<label>Email</label><input type="email" name="Email" placeholder="bob@example.com">`),
				strings.TrimSpace(`
					<label>City</label><input type="text" name="City" placeholder="City" value="New York" class="custom-city-class">`),
			}, "")),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b := &Builder{
				InputTemplate: tc.tpl,
			}
			got, err := b.Inputs(tc.arg)
			if err != nil {
				t.Errorf("Builder.Inputs() err = %v, want %v", err, nil)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("Builder.Inputs() = %v, want %v", got, tc.want)
			}
		})
	}
}

type testFieldError struct {
	field, err string
}

func (e testFieldError) Error() string {
	return fmt.Sprintf("invalid field: %v", e.field)
}

func (e testFieldError) FieldError() (field, err string) {
	return e.field, e.err
}

func TestBuilder_Inputs_errors(t *testing.T) {
	// Sanity check on our test type first
	tfe := testFieldError{
		field: "field",
		err:   "err",
	}
	var fe fieldError
	if !errors.As(tfe, &fe) {
		t.Fatalf("As(testFieldError, fieldError) = false")
	}
	if !errors.As(fmt.Errorf("wrapped: %w", tfe), &fe) {
		t.Fatalf("As(wrapped, fieldError) = false")
	}

	tpl := template.Must(template.New("").Funcs(FuncMap()).Parse(strings.TrimSpace(`
		<label>{{.Label}}</label>{{range errors}}<p>{{.}}</p>{{end}}
	`)))
	tests := []struct {
		name   string
		tpl    *template.Template
		arg    interface{}
		errors []error
		want   template.HTML
	}{
		{
			name: "label and input",
			tpl:  tpl,
			arg: struct {
				Name  string
				Email string `form:"type=email;placeholder=bob@example.com"`
			}{
				Name: "Michael Scott",
			},
			errors: []error{
				fmt.Errorf("wrapped: %w", testFieldError{
					field: "Name",
					err:   "is required",
				}),
				fmt.Errorf("first: %w", fmt.Errorf("second: %w", testFieldError{
					field: "Email",
					err:   "is taken",
				})),
			},
			want: template.HTML(strings.Join([]string{
				strings.TrimSpace(`
					<label>Name</label><p>is required</p>`),
				strings.TrimSpace(`
            <label>Email</label><p>is taken</p>`),
			}, "")),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b := &Builder{
				InputTemplate: tc.tpl,
			}
			got, err := b.Inputs(tc.arg, tc.errors...)
			if err != nil {
				t.Errorf("Builder.Inputs() err = %v, want %v", err, nil)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("Builder.Inputs() = %v, want %v", got, tc.want)
			}
		})
	}
}
