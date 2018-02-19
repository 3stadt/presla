package Handlers

import (
	"errors"
	"fmt"
	"github.com/labstack/echo"
	"github.com/spf13/afero"
	"mime"
	"net/http"
	"os"
	"path/filepath"
)

func (conf *Conf) Assets(c echo.Context) error {
	pres := c.Param("pres")
	file := c.Param("*")

	if pres == "" || file == "" {
		c.NoContent(http.StatusBadRequest)
	}

	mimeType := mime.TypeByExtension(filepath.Ext(file))

	path := conf.MarkdownPath + "/" + pres + "/" + file

	_, err := conf.Fs.Stat(path)
	if os.IsNotExist(err) {
		c.NoContent(http.StatusNotFound)
		return errors.New(fmt.Sprintf("file not found: %s", path))
	}

	var content []byte
	content, err = afero.ReadFile(conf.Fs, path)
	if err != nil {
		c.NoContent(http.StatusInternalServerError)
		return err
	}

	return c.Blob(http.StatusOK, mimeType, content)
}
