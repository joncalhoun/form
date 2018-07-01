# Form

Easily create HTML forms with Go structs.

[![go report card](https://goreportcard.com/badge/github.com/joncalhoun/form "go report card")](https://goreportcard.com/report/github.com/joncalhoun/form)
[![Build Status](https://travis-ci.org/joncalhoun/form.svg?branch=master)](https://travis-ci.org/joncalhoun/form)
[![MIT license](http://img.shields.io/badge/license-MIT-brightgreen.svg)](http://opensource.org/licenses/MIT)
[![GoDoc](https://godoc.org/github.com/joncalhoun/form?status.svg)](https://godoc.org/github.com/joncalhoun/form)

## Overview

The `form` package makes it easy to take a Go struct and turn it into an HTML form using whatever HTML format you want. Below is an example, along with the output. This entire example can be found in the [examples/readme](examples/readme) directory.

**Source Code**

```go
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

	pageTpl := template.Must(template.New("").Funcs(template.FuncMap{
		"inputs_for": fb.Inputs,
	}).Parse(`
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
```

**Relevant HTML** trimmed for brevity


```html
<form>
  <label >
    Street
  </label>
  <input  type="text" name="Street1" placeholder="123 Sample St" value="123 Known St">

  <label >
    Street (cont)
  </label>
  <input  type="text" name="Street2" placeholder="Apt 123" >

  <label >
    City
  </label>
  <input  type="text" name="City" placeholder="City" >

  <label >
    State
  </label>
  <input  type="text" name="State" placeholder="State" >
  <p>Or your Province</p>

  <label >
    Postal Code
  </label>
  <input  type="text" name="Zip" placeholder="Postal Code" >

  <label >
    Country
  </label>
  <input  type="text" name="Country" placeholder="Country" value="United States">
</form>
```

## How it works

The `form.Builder` type provides a single method - `Inputs` - which will parse the provided struct to determine which fields it contains, any values set for each field, and any struct tags provided for the form package. Once that information is parsed it will execute the provided `InputTemplate` field in the builder for each field in the struct, **including nested fields**.

Most of the time you will probably want to just make this helper available to your html templates via the `template.Funcs()` functions and the `template.FuncMap` type, as I did in the example above.

## I don't recommend tagging domain types

It is also worth mentioning that I don't really recommend adding `form` struct tags to your domain types, and I typically create types specifically used to generate forms. Eg:

```go
// This is my domain type
type User struct {
  ID           int
  Name         string
  Email        string
  PasswordHash string
}

// Somewhere else I'll create my html-specific type:
type signupForm struct {
  Name         string `form:"..."`
  Email        string `form:"type=email"`
  Password     string `form:"type=password"`
  Confirmation string `form:"type=password;label=Password Confirmation"`
}
```

## Parsing submitted forms

If you also need to parse forms created by this package, I recommend using the [gorilla/schema](https://github.com/gorilla/schema) package. This package *should* generate input names compliant with the `gorilla/schema` package by default, so as long as you don't change the names it should be pretty trivial to decode.

There is an example of this in the [examples/tailwind](examples/tailwind) directory.


## This may have bugs

This is a very early iteration of the package, and while it appears to be working for my needs chances are it doesn't cover every use case. If you do find one that isn't covered, try to provide a PR with a breaking test.


## Notes

This section is mostly for myself to jot down notes, but feel free to read away.

### Potential features

#### Parsing forms

Long term this could also support parsing forms, but gorilla/schema does a great job of that already so I don't see any reason to at this time. It would likely be easier to just make the default input names line up with what gorilla/schema expects and provide examples for how to use the two together.

#### Error rendering

I could also look into ways to handle errors and add messages to forms. This shouldn't be *too* hard to do. It would probably be something like an optional argument passed into the form builder and then we process it looking for implementations of an interface like:

```go
for _, err := range errors {
  if fe, ok := err.(interface{
    Field() string
    Message() string
  }); ok {
    map[fe.Field()] = fe.Message()
  }
}
```

Then we could pass it as an `Error` field each time we render:

```html
<div>
	<label>
		{{.Label}}
	</label>
	<input name="{{.Name}}" placeholder="{{.Placeholder}}" {{with .Value}}value="{{.}}"{{end}} class="{{with .Errors}}border-red{{end}}">
  {{range .Errors}}
    <p class="text-red-dark py-2 text-sm">{{.}}</p>
  {{end}}
	{{with .Footer}}
		<p class="text-grey pt-2 text-xs italic">{{.}}</p>
	{{end}}
</div>
```


#### Checkboxes and other data types

Maybe allow for various templates for different types, but for now this is possible to do in the HTML templates so it isn't completely missing.

#### Headers on nested structs

Let's say we have this type:

```go
type Nested struct {
  Name string
  Email string
  Address Address
}

type Address struct {
  Street1 string
  Street2 string
  // ...
}
```

It might make sense to make an optional way to add headers in the form when the nested Address portion is rendered, so the form looks like:

```
Name:    [    ]
Email:   [    ]

<Address Header Here>

Street1: [    ]
Street2: [    ]
...
```

This *should* be pretty easy to do with struct tags on the `Address Address` line.
