package PreslaTemplates

import (
	"github.com/labstack/echo"
	"html/template"
	"io"
)

// DefaultTemplate is used for template rendering
type DefaultTemplate struct {
	Template *template.Template
}

// Render is used for template rendering
func (t *DefaultTemplate) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.Template.ExecuteTemplate(w, name, data)
}
