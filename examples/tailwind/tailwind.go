// Produces a form like:
// https://www.dropbox.com/s/72z88osbcwik26n/Screenshot%202018-07-01%2014.13.31.png?dl=0&raw=1

package main

import (
	"encoding/json"
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
	<input class="shadow appearance-none border rounded w-full py-2 px-3 text-grey-darker leading-tight" {{with .ID}}id="{{.}}"{{end}} type="{{.Type}}" name="{{.Name}}" placeholder="{{.Placeholder}}" {{with .Value}}value="{{.}}"{{end}}>
	{{with .Footer}}
		<p class="text-grey pt-2 text-xs italic">{{.}}</p>
	{{end}}
</div>`

func main() {
	tpl := template.Must(template.New("").Parse(inputTpl))
	fb := form.Builder{
		InputTemplate: tpl,
	}

	pageTpl := template.Must(template.New("").Funcs(fb.FuncMap()).Parse(`
<html>
<head>
	<link href="https://cdn.jsdelivr.net/npm/tailwindcss/dist/tailwind.min.css" rel="stylesheet">
</head>
<body class="bg-grey-lighter">
	<div class="w-full max-w-xs mx-auto my-8">
		<form class="bg-white shadow-md rounded px-8 pt-6 pb-8 mb-4" action="/" method="post">
			{{inputs_for .}}
			<div class="flex items-center justify-between">
				<button class="bg-blue hover:bg-blue-dark text-white font-bold py-2 px-4 rounded" type="submit">
					Sign In
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
			pageTpl.Execute(w, loginForm{})
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
		var form loginForm
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

type loginForm struct {
	Email    string `form:"type=email;;label=Email Address;id=email;placeholder=bob@example.com"`
	Password string `form:"type=password;id=password;footer=Keep it secret!"`
}
