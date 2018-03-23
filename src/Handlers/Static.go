package Handlers

import (
	"github.com/labstack/echo"
	"io/ioutil"
	"mime"
	"net/http"
	"path/filepath"
)

// Static loads files like js and css from the static folder on the disk
func (conf *Conf) Static(c echo.Context) error {
	file := c.Param("*")
	pres := c.Param("pres")

	presConf, err := conf.getConf(pres)
	if err != nil {
		presConf = &PresentationConf{}
	}

	if presConf.StaticFiles == "" {
		presConf.StaticFiles = conf.StaticFiles
	}

	file = presConf.StaticFiles + "/" + pres + "/" + file

	mimeType := mime.TypeByExtension(filepath.Ext(file))

	content, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	return c.Blob(http.StatusOK, mimeType, content)
}
