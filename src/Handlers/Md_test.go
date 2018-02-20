package Handlers

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMdFromFsNotFound(t *testing.T) {
	mdPath := "/foo/bar/md"
	filename := "foobar.md"

	fs := afero.NewMemMapFs()

	conf := &Conf{
		MarkdownPath: mdPath,
		Fs:           fs,
	}

	e := echo.New()
	req, _ := http.NewRequest(echo.GET, "/md/foobar.md", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec).(echo.Context)
	c.SetParamNames("file")
	c.SetParamValues(filename)

	err := conf.Md(c)
	assert.NotNil(t, err)

	assert.Equal(t, 404, rec.Code)
}

func TestMdFromFs(t *testing.T) {
	mdPath := "/foo/bar/md"
	filename := "foobar.md"
	content := []byte(`# Foo
---
# Bar`)
	fs := afero.NewMemMapFs()
	err := fs.MkdirAll(mdPath, 0755)
	assert.Nil(t, err)
	afero.WriteFile(fs, fmt.Sprintf("%s/%s", mdPath, filename), content, 644)
	conf := &Conf{
		MarkdownPath: mdPath,
		Fs:           fs,
	}
	e := echo.New()
	req, _ := http.NewRequest(echo.GET, "/md/foobar.md", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec).(echo.Context)
	c.SetParamNames("file")
	c.SetParamValues(filename)

	err = conf.Md(c)
	assert.Nil(t, err)

	assert.Equal(t, 200, rec.Code)
	assert.Equal(t, string(content), rec.Body.String())
}
