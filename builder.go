package form

import (
	"html/template"
	"strings"
)

// Builder is used to build HTML forms/inputs
// for Go structs
type Builder struct {
	inputTpl *template.Template
}

// Inputs returns the HTML inputs for all the exported
// fields of the provided interface.
func (b *Builder) Inputs(v interface{}) template.HTML {
	fields := fields(v)
	var html template.HTML
	for _, field := range fields {
		var sb strings.Builder
		b.inputTpl.Execute(&sb, field)
		html = html + template.HTML(sb.String())
	}
	return html
}

// NewBuilder creates a Builder.
func NewBuilder(tpl *template.Template) *Builder {
	return &Builder{
		inputTpl: tpl,
	}
}
