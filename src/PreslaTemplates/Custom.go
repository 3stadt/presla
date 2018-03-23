package PreslaTemplates

import (
	"github.com/labstack/echo"
	"html/template"
	"io"
)

// CustomTemplate is used for template rendering
type CustomTemplate struct {
	Templates *template.Template
}

// Custom is used for template rendering
func Custom(location string) *CustomTemplate {
	return &CustomTemplate{
		Templates: template.Must(template.ParseGlob(location)),
	}
}

// Render is used for template rendering
func (t *CustomTemplate) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.Templates.ExecuteTemplate(w, name, data)
}
