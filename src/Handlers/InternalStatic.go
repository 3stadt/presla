package Handlers

import (
	"errors"
	"fmt"
	"github.com/labstack/echo"
	"mime"
	"net/http"
	"path/filepath"
)

func (conf *Conf) InternalStatic(c echo.Context) error {
	path := "static/" + c.Param("*")

	content, err := Asset(path)

	if err != nil {
		c.NoContent(http.StatusNotFound)
		return errors.New(fmt.Sprintf("asset not found: %s", path))
	}

	mimeType := mime.TypeByExtension(filepath.Ext(path))
	return c.Blob(http.StatusOK, mimeType, content)
}
