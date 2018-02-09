package Handlers

import (
	"bytes"
	"github.com/labstack/echo"
	"git.3stadt.com/3stadt/presla/src/PresLaTemplates"
	"html/template"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
)

func (conf *Conf) Md(c echo.Context) error {
	file := c.Param("file")

	if file == "info.md" {
		return conf.showInfo(c)
	}

	file = conf.MarkdownPath + "/" + file

	var content []byte
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	return c.Blob(http.StatusOK, "text/markdown; charset=UTF-8", content)
}

func (conf *Conf) showInfo(c echo.Context) error {
	var presentations []string

	files := make([]string, 1)
	files, err := filepath.Glob(conf.MarkdownPath + "/*.md")
	if err != nil {
		files[0] = "Error loading presentations: " + err.Error()
	}
	for _, file := range files {
		presentations = append(presentations, strings.TrimSuffix(filepath.Base(file), ".md"))
	}
	data := map[string]interface{}{
		"Presentations": presentations,
	}

	tpl, err := Asset("templates/info.md")
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
	return c.Blob(http.StatusOK, "text/markdown; charset=UTF-8", buf.Bytes())
}
