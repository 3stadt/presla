package Handlers

import (
	"github.com/labstack/echo"
	"mime"
	"net/http"
	"path/filepath"
)

func (conf *Conf) InternalStatic(c echo.Context) error {
	path := "static/" + c.Param("*")

	content, err := Asset(path)

	if err != nil {
		return c.String(http.StatusInternalServerError, "Asset not found: "+path)
	}

	mimeType := mime.TypeByExtension(filepath.Ext(path))
	return c.Blob(http.StatusOK, mimeType, content)
}
