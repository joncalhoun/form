# Easily create HTML forms with structs

The `form` package makes it easy to take a Go struct and turn it into an HTML form.

## Notes

The first use I want to support is roughly...

```go
// form.Builder could be a type or a function. Some of the arguments here
// would probably include something like an html/template that defines
// what an input tag's html should be, and I'll probably want to provide
// a few more common ones (eg bootstrap)
fb := form.Builder(...)

// The most common use of the builder will probably be inside of a
// template as a function.
tpl := template.New("").Parse(...).Funcs(template.FuncMap{
  "inputs_for": fb.Inputs,
})

// Then we'd have a template like...
tplStr := `
<form ...>
  <h3>Some section</h3>
  {{inputs_for .Customer}}
  <h3>Other section</h3>
  {{inputs_for .Address}}
</form>`
```

## Potential features

### Parsing forms

Long term this could also support parsing forms, but gorilla/schema does a great job of that already so I don't see any reason to at this time. It would likely be easier to just make the default input names line up with what gorilla/schema expects and provide examples for how to use the two together.

### Error rendering

I could also look into ways to handle errors and add messages to forms. This shouldn't be *too* hard to do.

```html
{{with err}} <!-- ?? -->
	<p class="text-red-dark py-2 text-sm">{{.}}</p>
{{end}}
```

### Checkboxes and other data types

Maybe allow for various templates for different types, but for now this is possible to do in the HTML templates so it isn't completely missing.
