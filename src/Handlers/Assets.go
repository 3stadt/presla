package Handlers

import (
	"github.com/labstack/echo"
	"io/ioutil"
	"mime"
	"net/http"
	"path/filepath"
)

func (conf *Conf) Assets(c echo.Context) error {
	pres := c.Param("pres")
	file := c.Param("*")

	mimeType := mime.TypeByExtension(filepath.Ext(file))

	path := conf.MarkdownPath + "/" + pres + "/" + file

	var content []byte
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	return c.Blob(http.StatusOK, mimeType, content)
}
