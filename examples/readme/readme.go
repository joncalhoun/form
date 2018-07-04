package main

import (
	"html/template"
	"net/http"

	"github.com/joncalhoun/form"
)

var inputTpl = `
<label {{with .ID}}for="{{.}}"{{end}}>
	{{.Label}}
</label>
<input {{with .ID}}id="{{.}}"{{end}} type="{{.Type}}" name="{{.Name}}" placeholder="{{.Placeholder}}" {{with .Value}}value="{{.}}"{{end}}>
{{with .Footer}}
  <p>{{.}}</p>
{{end}}
`

// Address is an example type used to demonstrate the form package.
type Address struct {
	Street1 string `form:"label=Street;placeholder=123 Sample St"`
	Street2 string `form:"label=Street (cont);placeholder=Apt 123"`
	City    string
	State   string `form:"footer=Or your Province"`
	Zip     string `form:"label=Postal Code"`
	Country string
}

func main() {
	tpl := template.Must(template.New("").Parse(inputTpl))
	fb := form.Builder{
		InputTemplate: tpl,
	}

	pageTpl := template.Must(template.New("").Funcs(fb.FuncMap()).Parse(`
		<html>
		<body>
			<form>
				{{inputs_for .}}
			</form>
		</body>
		</html>`))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		pageTpl.Execute(w, Address{
			Street1: "123 Known St",
			Country: "United States",
		})
	})
	http.ListenAndServe(":3000", nil)
}
