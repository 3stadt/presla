package Handlers

import (
	"bytes"
	"github.com/3stadt/presla/src/PresLaTemplates"
	"github.com/labstack/echo"
	"html/template"
	"io/ioutil"
	"net/http"
)

func (conf *Conf) Svg(c echo.Context) error {
	text := conf.FooterText
	if text == "" {
		text = "@shopware"
	}
	data := map[string]interface{}{
		"Text": text,
	}

	presConf, err := conf.getConf(c.Param("pres"))
	if err != nil {
		presConf = &PresentationConf{}
	}

	var tpl []byte

	if presConf.TemplatePath == "" && conf.TemplatePath == "" {
		tpl, err = Asset("templates/footer-text.svg")
		if err != nil {
			return err
		}
	}

	if tpl == nil {
		if presConf.TemplatePath == "" {
			presConf.TemplatePath = conf.TemplatePath
		}
		tpl, err = ioutil.ReadFile(presConf.TemplatePath + "/footer-text.svg")
		if err != nil {
			return err
		}
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

	return c.Blob(http.StatusOK, "image/svg+xml", buf.Bytes())
}
