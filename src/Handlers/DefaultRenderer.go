package Handlers

import (
	"bytes"
	"github.com/labstack/echo"
	"git.3stadt.com/3stadt/presla/src/PresLaTemplates"
	"html/template"
	"net/http"
)

func renderWithDefaultTemplate(data map[string]interface{}, c echo.Context) error {
	tpl, err := Asset("templates/index.html")
	if err != nil {
		return err
	}
	parsedTemplate, err := template.New("default").Parse(string(tpl))
	if err != nil {
		return err
	}
	t := &PresLaTemplates.DefaultTemplate{
		Template: parsedTemplate,
	}
	buf := new(bytes.Buffer)
	err = t.Render(buf, "default", data, c)
	if err != nil {
		return err
	}
	return c.Blob(http.StatusOK, "text/html", buf.Bytes())
}
