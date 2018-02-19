package main

import (
	"github.com/3stadt/presla/src/Handlers"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestConfIsCreated(t *testing.T) {
	logger.SetLevel(logrus.PanicLevel) // do not clutter test output
	fs := new(afero.MemMapFs)

	dir, err := os.Getwd()
	assert.Nil(t, err)

	defaultConfigPath := filepath.Clean(dir + "/presla.toml")
	assert.Nil(t, err)

	path, err := getConfPath("", fs)
	assert.Nil(t, err)
	assert.Equal(t, defaultConfigPath, path)

	exists, err := afero.Exists(fs, defaultConfigPath)
	assert.Nil(t, err)
	assert.True(t, exists)

	content, err := afero.ReadFile(fs, defaultConfigPath)
	assert.Nil(t, err)
	assert.NotEmpty(t, content)
}

func TestBasicConfContentIsCreated(t *testing.T) {
	dir, err := os.Getwd()
	assert.Nil(t, err)

	defaultConfigFile := filepath.Clean(dir + "/presla.toml")
	assert.Nil(t, err)

	defaultMarkdownPath := filepath.Clean(dir)
	assert.Nil(t, err)

	expected := Config{
		ConfigFile:      defaultConfigFile,
		MarkdownPath:    defaultMarkdownPath,
		FooterText:      "please edit presla.toml",
		ListenOn:        "localhost:8080",
		TemplatePath:    "",
		StaticFiles:     "",
		Presentations:   []Handlers.PresentationConf(nil),
		CustomExecutors: "",
		LogLevel:        "",
		LogFormat:       "",
	}

	fs := new(afero.MemMapFs)
	configPath, err := getConfPath("", fs)
	assert.Nil(t, err)

	err = readMainConfig(configPath, fs)
	assert.Nil(t, err)
	assert.Equal(t, expected, conf)
}
