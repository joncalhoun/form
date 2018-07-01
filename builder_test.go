package form

import (
	"html/template"
	"reflect"
	"strings"
	"testing"
)

func TestBuilder_Inputs(t *testing.T) {
	tpl := template.Must(template.New("").Parse(strings.TrimSpace(`
		<label>{{.Label}}</label><input type="{{.Type}}" name="{{.Name}}" placeholder="{{.Placeholder}}"{{with .Value}} value="{{.}}"{{end}}>
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
			}{
				Name: "Michael Scott",
			},
			want: template.HTML(strings.Join([]string{
				strings.TrimSpace(`
					<label>Name</label><input type="text" name="Name" placeholder="Name" value="Michael Scott">`),
				strings.TrimSpace(`
					<label>Email</label><input type="email" name="Email" placeholder="bob@example.com">`),
			}, "")),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b := &Builder{
				InputTemplate: tc.tpl,
			}
			if got := b.Inputs(tc.arg); !reflect.DeepEqual(got, tc.want) {
				t.Errorf("Builder.Inputs() = %v, want %v", got, tc.want)
			}
		})
	}
}
