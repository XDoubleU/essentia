// Package tpl provides tools for templates.
package tpl

import (
	"html/template"
	"io"
)

// RenderWithPanic tries to render a template and panics when this returns an error.
func RenderWithPanic(tpl *template.Template, wr io.Writer, name string, data any) {
	err := tpl.ExecuteTemplate(wr, name, data)
	if err != nil {
		panic(err)
	}
}
