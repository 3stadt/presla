package Handlers

import (
	"github.com/labstack/echo"
	"net/http"
)

// Favicon returns the default presla favicon to the browser
func (conf *Conf) Favicon(c echo.Context) error {
	content, err := Asset("static/favicon.ico")

	if err != nil {
		return c.String(http.StatusInternalServerError, "")
	}

	return c.Blob(http.StatusOK, "image/x-icon", content)
}
