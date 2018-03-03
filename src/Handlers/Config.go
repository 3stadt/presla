package Handlers

import (
	"github.com/labstack/echo"
	"net/http"
	"github.com/3stadt/presla/src/PreslaTemplates"
	"html/template"
	"bytes"
)

func (conf *Conf) SaveConfig(c echo.Context) error {
	return nil
}

func (conf *Conf) Config(c echo.Context) error {
	data := map[string]interface{}{
		"Title": "info",
	}

	return renderWithFormTemplate(data, c)
}

func renderWithFormTemplate(data map[string]interface{}, c echo.Context) error {
	tpl, err := Asset("templates/config.html")
	if err != nil {
		return err
	}
	parsedTemplate, err := template.New("default").Parse(string(tpl))
	if err != nil {
		return err
	}
	t := &PreslaTemplates.DefaultTemplate{
		Template: parsedTemplate,
	}
	buf := new(bytes.Buffer)
	err = t.Render(buf, "default", data, c)
	if err != nil {
		return err
	}
	return c.Blob(http.StatusOK, "text/html", buf.Bytes())
}
