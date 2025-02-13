package tpl_test

import (
	"html/template"
	"testing"

	"github.com/XDoubleU/essentia/pkg/tpl"
	"github.com/stretchr/testify/assert"
)

func TestRenderWithPanic(t *testing.T) {
	template := template.New("")
	assert.Panics(t, func() { tpl.RenderWithPanic(template, nil, "", nil) })
}
