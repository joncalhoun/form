// Produces a form like:
// https://www.dropbox.com/s/72z88osbcwik26n/Screenshot%202018-07-01%2014.13.31.png?dl=0&raw=1

package main

import (
	"html/template"
	"net/http"

	"github.com/joncalhoun/form"
)

var inputTpl = `
<div class="form-group">
	<label {{with .ID}}for="{{.}}"{{end}}>
		{{.Label}}
	</label>
	{{if eq .Type "textarea"}}
		<textarea class="form-control" {{with .ID}}id="{{.}}"{{end}} name="{{.Name}}" rows="3" placeholder="{{.Placeholder}}">{{with .Value}}{{.}}{{end}}</textarea>
	{{else}}
		<input type="{{.Type}}" class="form-control" {{with .ID}}id="{{.}}"{{end}} name="{{.Name}}" placeholder="{{.Placeholder}}" {{with .Value}}value="{{.}}"{{end}}>
	{{end}}
	{{with .Footer}}
		<small class="form-text text-muted">
			{{.}}
		</small>
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
	<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0/css/bootstrap.min.css" integrity="sha384-Gn5384xqQ1aoWXA+058RXPxPg6fy4IWvTNh0E263XmFcJlSAwiGgFAW/dAiS6JXm" crossorigin="anonymous">
	<script src="https://code.jquery.com/jquery-3.2.1.slim.min.js" integrity="sha384-KJ3o2DKtIkvYIK3UENzmM7KCkRr/rE9/Qpg6aAZGJwFDMVNA/GpGFF93hXpG5KkN" crossorigin="anonymous"></script>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.12.9/umd/popper.min.js" integrity="sha384-ApNbgh9B+Y1QKtv3Rn7W3mgPxhU9K/ScQsAP7hUibX39j7fakFPskvXusvfa0b4Q" crossorigin="anonymous"></script>
	<script src="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0/js/bootstrap.min.js" integrity="sha384-JZR6Spejh4U02d8jOt6vLEHfe/JQGiRRSQQxSfFWpi1MquVdAyjUar5+76PVCmYl" crossorigin="anonymous"></script>
</head>
<body>
	<div class="container">
		<!-- Button trigger modal -->
		<button type="button" class="btn btn-primary" data-toggle="modal" data-target="#exampleModal">
			Sign up
		</button>

		<!-- Modal -->
		<div class="modal fade" id="exampleModal" tabindex="-1" role="dialog">
			<div class="modal-dialog" role="document">
				<div class="modal-content">
					<div class="modal-header">
						<h5 class="modal-title" id="exampleModalLabel">Sign Up</h5>
						<button type="button" class="close" data-dismiss="modal" aria-label="Close">
							<span aria-hidden="true">&times;</span>
						</button>
					</div>
					<div class="modal-body">
						<form action="/do-stuff" method="post">
							{{inputs_for .}}
						</form>
					</div>
					<div class="modal-footer">
						<button type="button" class="btn btn-secondary" data-dismiss="modal">Close</button>
						<button type="button" class="btn btn-primary">Sign Up</button>
					</div>
				</div>
			</div>
		</div>
	</div>
</body>
</html>
	`))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		pageTpl.Execute(w, userForm{
			Name:    "Jon Calhoun",
			Email:   "jon@calhoun.io",
			Bio:     "I like to write Go code!",
			Ignored: "info here should be ignored entirely because of the \"-\" form tag value",
		})
	})
	http.ListenAndServe(":3000", nil)
}

type userForm struct {
	Name  string
	Email string `form:"type=email"`

	Password     string `form:"type=password"`
	Confirmation string `form:"display=Password Confirmation;type=password"`

	Bio string `form:"type=textarea;placeholder=Tell us a bit about yourself"`

	Ignored string `form:"-"`
}
