package Handlers

import (
	"bytes"
	"git.3stadt.com/3stadt/presla/src/PresLaTemplates"
	"github.com/labstack/echo"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"
)

func (conf *Conf) Presentation(c echo.Context) error {
	data := map[string]interface{}{
		"Pres":  c.Param("pres"),
		"Title": strings.Replace(c.Param("pres"), "_", " ", -1),
	}

	presConf, err := conf.getConf(c.Param("pres"))
	if err != nil {
		presConf = &PresentationConf{}
	}

	if presConf.TemplatePath == "" && conf.TemplatePath == "" {
		return renderWithDefaultTemplate(data, c)
	}

	if presConf.TemplatePath == "" {
		presConf.TemplatePath = conf.TemplatePath
	}

	tpl, err := ioutil.ReadFile(presConf.TemplatePath + "/index.html")
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
