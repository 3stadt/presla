package Handlers

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/spf13/afero"
	"mime"
	"net/http"
	"os"
	"path/filepath"
)

// Assets serves assets like js and css files via http.
// Depending on the configuration, the assets are loaded from bindata or the file system
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
		return fmt.Errorf("file not found: %s", path)
	}

	var content []byte
	content, err = afero.ReadFile(conf.Fs, path)
	if err != nil {
		c.NoContent(http.StatusInternalServerError)
		return err
	}

	return c.Blob(http.StatusOK, mimeType, content)
}
