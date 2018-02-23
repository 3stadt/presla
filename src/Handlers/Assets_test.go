package Handlers

import (
	"github.com/labstack/echo"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPresentationNotFound(t *testing.T) {
	fs := afero.NewMemMapFs()
	fs.MkdirAll("/markdown", 0755)

	conf := &Conf{
		MarkdownPath: "/markdown",
		Fs:           fs,
	}

	xAssetNotFound(t, conf)
}

func TestAssetFileNotFound(t *testing.T) {
	fs := afero.NewMemMapFs()
	fs.MkdirAll("/markdown/foobar", 0755)

	conf := &Conf{
		MarkdownPath: "/markdown",
		Fs:           fs,
	}

	xAssetNotFound(t, conf)
}

func TestAssetFileFound(t *testing.T) {
	fs := afero.NewMemMapFs()
	fs.MkdirAll("/markdown/foobar", 0755)
	gif := []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x72\xf7\x74\xb3\xb0\x4c\x64\x64\x60\x64\x60\xf8\xcf\xa0\xc3\xc0\xc0\xc0\x00\x66\x33\x31\x58\x03\x02\x00\x00\xff\xff\xc7\x52\xbc\x6d\x1a\x00\x00\x00")
	afero.WriteFile(fs, "/markdown/foobar/bar.gif", gif, 0644)

	conf := &Conf{
		MarkdownPath: "/markdown",
		Fs:           fs,
	}

	e := echo.New()
	e.GET("/md/:pres/*", conf.Assets)
	req, _ := http.NewRequest(echo.GET, "/md/foobar/bar.gif", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec).(echo.Context)
	c.SetParamNames("pres", "*")
	c.SetParamValues("foobar", "bar.gif")

	err := conf.Assets(c)

	assert.Nil(t, err)
	assert.Equal(t, 200, rec.Code)
	assert.Equal(t, gif, rec.Body.Bytes())
	assert.Equal(t, "image/gif", rec.HeaderMap.Get("Content-Type"))
}

func xAssetNotFound(t *testing.T, conf *Conf) {
	e := echo.New()
	e.GET("/md/:pres/*", conf.Assets)
	req, _ := http.NewRequest(echo.GET, "/md/foobar/bar.gif", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec).(echo.Context)
	c.SetParamNames("pres", "*")
	c.SetParamValues("foobar", "bar.gif")

	err := conf.Assets(c)

	assert.NotNil(t, err)
	assert.Equal(t, 404, rec.Code)
	assert.Empty(t, rec.Body.String())
}
