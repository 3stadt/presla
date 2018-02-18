package PreslaTemplates

import (
	"github.com/labstack/echo"
	"html/template"
	"io"
)

type customTemplate struct {
	templates *template.Template
}

func Custom(location string) *customTemplate {
	return &customTemplate{
		templates: template.Must(template.ParseGlob(location)),
	}
}

func (t *customTemplate) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
