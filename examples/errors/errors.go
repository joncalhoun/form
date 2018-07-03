// Produces a form like:
//
package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/schema"
	"github.com/joncalhoun/form"
)

var inputTpl = `
<div class="mb-4">
	<label class="block text-grey-darker text-sm font-bold mb-2" {{with .ID}}for="{{.}}"{{end}}>
		{{.Label}}
	</label>
	<input class="shadow appearance-none border rounded w-full py-2 px-3 text-grey-darker leading-tight {{if errors}}border-red{{end}}" {{with .ID}}id="{{.}}"{{end}} type="{{.Type}}" name="{{.Name}}" placeholder="{{.Placeholder}}" {{with .Value}}value="{{.}}"{{end}}>
	{{range errors}}
		<p class="text-red pt-2 text-xs italic">{{.}}</p>
	{{end}}
	{{with .Footer}}
		<p class="text-grey pt-2 text-xs italic">{{.}}</p>
	{{end}}
</div>`

func main() {
	tpl := template.Must(template.New("").Funcs(form.FuncMap()).Parse(inputTpl))
	fb := form.Builder{
		InputTemplate: tpl,
	}

	pageTpl := template.Must(template.New("").Funcs(fb.FuncMap()).Parse(`
<html>
<head>
	<link href="https://cdn.jsdelivr.net/npm/tailwindcss/dist/tailwind.min.css" rel="stylesheet">
</head>
<body class="bg-grey-lighter">
	<div class="w-full max-w-md mx-auto my-8">
		<form class="bg-white shadow-md rounded px-8 pt-6 pb-8 mb-4" action="/" method="post">
			{{inputs_and_errors_for .Form .Errors}}
			<div class="flex items-center justify-between">
				<button class="bg-blue hover:bg-blue-dark text-white font-bold py-2 px-4 rounded" type="submit">
					Submit
				</button>
			</div>
		</form>
		<p class="text-center text-grey text-xs">
			&copy; 2018 Acme Corp. All rights reserved.
		</p>
	</div>
</body>
</html>
	`))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			w.Header().Set("Content-Type", "text/html")
			data := struct {
				Form   nestedForm
				Errors []error
			}{
				Form: nestedForm{
					Name:    "Michael Scott",
					Email:   "michael@dunder.com",
					Address: nil,
				},
				Errors: []error{
					fieldError{
						Field: "Email",
						Issue: "is already taken",
					},
					fieldError{
						Field: "Address.Street1",
						Issue: "is required",
					},
					fieldError{
						Field: "Address.City",
						Issue: "is required",
					},
					fieldError{
						Field: "Address.State",
						Issue: "must be a US state",
					},
					fieldError{
						Field: "Address.Zip",
						Issue: "must be 5 digits",
					},
					fieldError{
						Field: "Address.Zip",
						Issue: "is required",
					},
				},
			}
			err := pageTpl.Execute(w, data)
			if err != nil {
				panic(err)
			}
			return
		case http.MethodPost:
		default:
			http.NotFound(w, r)
			return
		}

		// You can also process these forms using the gorilla/schema package.
		r.ParseForm()
		dec := schema.NewDecoder()
		dec.IgnoreUnknownKeys(true)
		var form nestedForm
		err := dec.Decode(&form, r.PostForm)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		b, _ := json.Marshal(form)
		w.Write(b)
	})
	http.ListenAndServe(":3000", nil)
}

type nestedForm struct {
	Name    string
	Email   string
	Address *address
}

type address struct {
	Street1 string
	Street2 string
	City    string
	State   string
	Zip     string `form:"label=Postal Code"`
}

type fieldError struct {
	Field string
	Issue string
}

func (fe fieldError) Error() string {
	return fmt.Sprintf("%v: %v", fe.Field, fe.Issue)
}

// You can implement this however you want - this is just how I'm doing it.
func (fe fieldError) FieldError() (field, err string) {
	return fe.Field, fe.Issue
}
