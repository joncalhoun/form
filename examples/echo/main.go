package main

import (
	"github.com/joncalhoun/form"
	"html/template"
)

// dependency injection
type echoApp struct{}

func main() {
	app := &echoApp{}

	tpl := template.Must(template.New("input").ParseFiles("public/partials/input.gohtml"))
	fb := form.Builder{
		InputTemplate: tpl,
	}
	t := &Template{
		templates: template.Must(template.New("form").Funcs(fb.FuncMap()).ParseFiles("public/views/form.gohtml")),
	}

	e := echo.New()
	e.Renderer = t

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}, error=${error}\n",
	}))

	e.GET("/insert", app.insert)
	e.POST("/insert", app.insertPost)

	e.Logger.Fatal(e.Start(":1323"))
}
