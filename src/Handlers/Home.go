package Handlers

import (
	"github.com/labstack/echo"
)

// Home serves the default presentation built into presla
func (conf *Conf) Home(c echo.Context) error {
	data := map[string]interface{}{
		"Pres":  "info",
		"Title": "info",
	}

	return renderWithDefaultTemplate(data, c)
}
