package Handlers

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// make sure all elements are in default template
func TestDefaultTemplate(t *testing.T) {
	e := echo.New()
	req, _ := http.NewRequest(echo.GET, "/md/info.md", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec).(echo.Context)

	data := map[string]interface{}{
		"Pres":  "PRES",
		"Title": "TITLE",
	}

	err := renderWithDefaultTemplate(data, c)
	assert.Nil(t, err)

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(rec.Body.String()))
	assert.Nil(t, err)

	title := doc.Find("title")
	assert.Equal(t, 1, title.Length())
	assert.Equal(t, "[TITLE]", title.First().Text())

	varSetter := doc.Find("script[type='text/javascript']:not([src])")
	assert.Equal(t, 1, varSetter.Length())
	assert.Equal(t, "let pres = \"PRES\";", strings.TrimSpace(varSetter.First().Text()))

	css := doc.Find("link[rel='stylesheet'][href='/static/internal/css/theme-presla.css']")
	assert.Equal(t, 1, css.Length())

	textarea := doc.Find("textarea[id='source']")
	assert.Equal(t, 1, textarea.Length())

	remark := doc.Find("script[src='/static/internal/js/remark-latest.min.js'][type='text/javascript']")
	assert.Equal(t, 1, remark.Length())

	remarkLoader := doc.Find("script[src='/static/internal/js/remark-loader.js'][type='text/javascript']")
	assert.Equal(t, 1, remarkLoader.Length())

	ace := doc.Find("script[src='/static/internal/js/ace/ace.js'][type='text/javascript'][charset='utf-8']")
	assert.Equal(t, 1, ace.Length())

	editorLoader := doc.Find("script[src='/static/internal/js/editor-loader.js'][type='text/javascript']")
	assert.Equal(t, 1, editorLoader.Length())
}
