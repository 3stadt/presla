package Handlers

import (
	"github.com/labstack/echo"
)

func (conf *Conf) Home(c echo.Context) error {
	data := map[string]interface{}{
		"Pres":  "info",
		"Title": "info",
	}

	return renderWithDefaultTemplate(data, c)
}
