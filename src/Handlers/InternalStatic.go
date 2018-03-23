package Handlers

import (
	"fmt"
	"github.com/labstack/echo"
	"mime"
	"net/http"
	"path/filepath"
)

// InternalStatic is used to fetch assets explicit from the bindata file
func (conf *Conf) InternalStatic(c echo.Context) error {
	path := "static/" + c.Param("*")

	content, err := Asset(path)

	if err != nil {
		c.NoContent(http.StatusNotFound)
		return fmt.Errorf("asset not found: %s", path)
	}

	mimeType := mime.TypeByExtension(filepath.Ext(path))
	return c.Blob(http.StatusOK, mimeType, content)
}
