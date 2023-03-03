package main

import (
	"html/template"
	"io"
	"net/http"
)

// using the Echo Template Interface
type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

// Form struct
type EmployeeInsertForm struct {
	Name     string
	Dni      string
	Email    string
	Address  string `form:"type:email"`
	Phone    string
	Password string `form:"type:password"`
}

// handlers
func (app *echoApp) insert(c echo.Context) error {
	return c.Render(http.StatusOK, "form", EmployeeInsertForm{})
}

func (app echoApp) insertPost(c echo.Context) error {
	name := c.FormValue("Name")
	dni := c.FormValue("Dni")
	email := c.FormValue("Email")
	address := c.FormValue("Address")
	phone := c.FormValue("Phone")
	password := c.FormValue("Password")

	return c.String(http.StatusOK, name)
}
