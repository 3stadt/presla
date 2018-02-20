package Handlers

import (
	"github.com/labstack/echo"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestInternalStaticFileNotFound(t *testing.T) {
	file := "/foo/bar.tiff"
	fs := afero.NewMemMapFs()

	conf, c, rec := getInternalStaticContext(fs, file)
	c.SetParamNames("*")
	c.SetParamValues(file)

	err := conf.InternalStatic(c)
	assert.NotNil(t, err)
	assert.Equal(t, 404, rec.Code)
}

func TestInternalStaticThemeCss(t *testing.T) {
	testInternalStaticX(t, "css/theme-presla.css", "text/css; charset=utf-8")
}

func TestInternalStaticRemarkJs(t *testing.T) {
	testInternalStaticX(t, "js/remark-latest.min.js", "application/javascript")
}

func TestInternalStaticRemarkLoader(t *testing.T) {
	testInternalStaticX(t, "js/remark-loader.js", "application/javascript")
}

func TestInternalStaticAce(t *testing.T) {
	testInternalStaticX(t, "js/ace/ace.js", "application/javascript")
}

func TestInternalStaticEditorLoader(t *testing.T) {
	testInternalStaticX(t, "js/editor-loader.js", "application/javascript")
}

func testInternalStaticX(t *testing.T, file string, mimeType string) {
	fs := afero.NewMemMapFs()

	conf, c, rec := getInternalStaticContext(fs, file)
	c.SetParamNames("*")
	c.SetParamValues(file)
	err := conf.InternalStatic(c)

	assert.Nil(t, err)
	assert.Equal(t, 200, rec.Code)
	assert.Equal(t, mimeType, rec.HeaderMap.Get("Content-Type"))
	assert.NotEmpty(t, rec.Body.String())
}

func getInternalStaticContext(fs afero.Fs, staticPath string) (*Conf, echo.Context, *httptest.ResponseRecorder) {

	conf := &Conf{
		MarkdownPath: "/markdown",
		Fs:           fs,
	}
	e := echo.New()
	e.GET("/static/internal/:pres/*", conf.InternalStatic)
	req, _ := http.NewRequest(echo.GET, "/static"+staticPath, nil)
	rec := httptest.NewRecorder()

	return conf, e.NewContext(req, rec).(echo.Context), rec
}
