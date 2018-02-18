package PreslaTemplates

import (
	"github.com/labstack/echo"
	"html/template"
	"io"
)

type DefaultTemplate struct {
	Template *template.Template
}

func (t *DefaultTemplate) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.Template.ExecuteTemplate(w, name, data)
}
