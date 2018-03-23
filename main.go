package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"github.com/3stadt/presla/src/Handlers"
	"github.com/3stadt/presla/src/PreslaTemplates"
	"github.com/BurntSushi/toml"
	"github.com/blang/semver"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/mitchellh/go-homedir"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"os"
	"path/filepath"
	"strings"
)

// Config holds the main configuration. The config file presla.toml is loaded into this struct.
type Config struct {
	ConfigFile         string
	MarkdownPath       string
	FooterText         string
	ListenOn           string
	TemplatePath       string
	StaticFiles        string
	Presentations      []Handlers.PresentationConf
	CustomExecutors    string
	LogLevel           string
	LogFormat          string
	CheckUpdateOnStart bool
}

var conf = Config{CheckUpdateOnStart: true}
var logger = log.New()

var version = "vlatest"

/**
Set up HTTP routes, start server
*/
func main() {
	// Use Disk
	fs := afero.NewOsFs()

	// If user sets conf param, use that as starting point
	configPathFlag := ""
	flag.StringVar(&configPathFlag, "conf", "", "The path to the configuration file")
	flag.BoolVar(&conf.CheckUpdateOnStart, "auto-update", conf.CheckUpdateOnStart, "Whether to enable or disable update check on start")
	flag.Parse()

	// Search for the config file to use or create one in working dir
	configPath, err := getConfPath(configPathFlag, fs)
	checkErr(err)

	// Read in config from above into global conf var
	err = readMainConfig(configPath, fs)
	checkErr(err)

	conf.checkForUpdate()

	handler := &Handlers.Conf{
		ConfigFile:      conf.ConfigFile,
		MarkdownPath:    conf.MarkdownPath,
		FooterText:      conf.FooterText,
		TemplatePath:    conf.TemplatePath,
		StaticFiles:     conf.StaticFiles,
		Presentations:   conf.Presentations,
		CustomExecutors: conf.CustomExecutors,
		LogLevel:        conf.LogLevel,
		LogFormat:       conf.LogFormat,
		Fs:              afero.NewOsFs(),
		SyncedEditorPub: make(chan Handlers.SyncedEditor),
		SyncedEditorSub: make(map[int]*Handlers.SyncedEditorWriter),
	}

	// Used for synchronization between editors in presenter view and copied browser on beamer/monitor.
	// If there is a message from a browser on any open websocket, this message is pushed to all websockets on all browsers
	go func() {
		for {
			sep := <-handler.SyncedEditorPub // Wait for message from any websocket/browser
			for _, ses := range handler.SyncedEditorSub {
				ses.Writer(sep, ses.Ws) // write the message to all websockets/browsers.
			}
		}
	}()

	e := echo.New()
	e.Use(middleware.Recover())

	for _, c := range conf.Presentations {
		_, err := os.Stat(c.TemplatePath)
		if c.TemplatePath != "" && c.PresentationName != "" && err == nil {
			e.Renderer = PreslaTemplates.Custom(c.TemplatePath)
		}
	}
	e.GET("/static/internal/*", handler.InternalStatic)
	e.GET("/static/:pres/*", handler.Static)
	e.GET("/favicon.ico", handler.Favicon)
	e.GET("/svg/footer-text.svg", handler.Svg)
	e.GET("/md/:file", handler.Md)
	e.GET("/md/:pres/*", handler.Assets)
	e.GET("/editorsync", handler.EditorSync)
	e.POST("/exec", handler.Exec)
	e.GET("/:pres", handler.Presentation)
	e.GET("/", handler.Home)
	logger.Infof("Starting server at: %s", fmt.Sprintf("http://%s", conf.ListenOn))
	logger.Infof("=> Use Ctrl+c to quit Presla")
	e.Start(conf.ListenOn)
}

func (conf *Config) checkForUpdate() {
	if version == "vlatest" { // version is changed on compile via ldflags, see makefile
		log.Info("using development version, update check deactivated")
		return
	}

	if conf.CheckUpdateOnStart == false {
		log.Info("Auto update check is disabled in config.")
		return
	}

	ver := version[1:]
	latest, found, err := selfupdate.DetectLatest("3stadt/presla")
	if err != nil {
		log.Error("error occurred while detecting version: ", err.Error())
		return
	}

	v, err := semver.Parse(ver)
	if err != nil {
		log.Error("could not parse current version: ", err.Error())
		return
	}

	if !found || latest.Version.LTE(v) {
		log.Info("using latest version")
		return
	}

	log.Warn("New version available")
	fmt.Println("Please note: Automatic update to a new version always uses the uncompressed binary.")
	fmt.Println("----------")
	fmt.Println(latest.ReleaseNotes)
	fmt.Println("----------")
	fmt.Print("Do you want to update to version ", latest.Version, "? (y/N): ")
	input, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil || strings.ToLower(strings.TrimSpace(input)) != "y" {
		fmt.Println("Skipping update")
		fmt.Printf("You can download the update manually at %s\n", latest.URL)
		return
	}

	log.Warn("Updating to latest version, please be patient...")

	ex, err := os.Executable()
	if err != nil {
		log.Error("error occurred while updating binary: ", err)
		return
	}

	if err := selfupdate.UpdateTo(latest.AssetURL, ex); err != nil {
		log.Error("error occurred while updating binary: ", err)
		return
	}
	log.Info("successfully updated to version ", latest.Version)
}

func checkErr(e error) {
	if e != nil {
		msg := e.Error()
		if e.Error() == "toml: cannot load TOML value of type map[string]interface {} into a Go slice" {
			msg = "You must use double brackets in your config file, like [[This]], instead of single brackets like [This]"
		}
		logger.Fatalf("A critical error occurred: %s", msg)
	}
}

func readMainConfig(configPath string, fs afero.Fs) error {
	tomlData, err := afero.ReadFile(fs, configPath)
	if err != nil {
		return err
	}

	_, err = toml.Decode(string(tomlData), &conf)
	if err != nil {
		if strings.Contains(err.Error(), "expected eight hexadecimal digits after") { // make error message from toml lexer readable for users
			msg := "ERROR: please check your config file for unescaped backslashes. E.g. on windows use 'C:\\\\Users\\\\' instead of 'C:\\Users\\'"
			return errors.New(msg)
		}
		return err
	}

	if conf.ConfigFile, err = filepath.Abs(configPath); err != nil {
		return err
	}

	if conf.MarkdownPath, err = filepath.Abs(conf.MarkdownPath); err != nil {
		return err
	}

	return nil
}

func getConfPath(configPath string, fs afero.Fs) (string, error) {

	// By preserving user input, default config will be created on user specified if it doesn't exist yet
	if configPath == "" {
		dir, err := os.Getwd()
		if err != nil {
			return "", err
		}
		configPath = dir + "/presla.toml"
	}

	configPath = filepath.Clean(configPath)

	locations := []string{configPath}

	home, err := homedir.Dir()

	// Only search other configs when we have a home directory
	if err == nil {
		locations = append(locations, filepath.Clean(home+"/.presla.toml"), filepath.Clean(home+"/.config/presla.toml"))
	} else {
		logger.Errorf("could not find home directory: %s", err.Error())
	}

	for _, location := range locations {
		// If file exists, return the path immediately
		if _, err := fs.Stat(location); err == nil {
			logger.Infof("using config file: %s", location)
			return location, nil
		}
		logger.Infof("no config file at %s", location)
	}

	defaultConfig := []byte(getDefaultConfig())
	err = afero.WriteFile(fs, configPath, defaultConfig, 0644)
	if err != nil {
		return "", err
	}
	logger.Infof("created and using config file with default values: %s", configPath)
	return configPath, nil
}

func getDefaultConfig() string {
	return `## The path to your markdown files.
## One markdown file holds one presentation
MarkdownPath="./"

## Whatever you want to show as text when including /svg/footer-text.svg
## By default shown in the lower right corner
FooterText="please edit presla.toml"

## The port to bind on. You should use localhost as host
ListenOn="localhost:8080"

## Optional: Path to your own template.
## Needs the index.html holding remarkjs, an info.md as starting point and footer-text.svg 
# TemplatePath="/home/user/Documents/presla-theme/templates"

## Optional: path to the templates static files
## Holds css, js, fonts and images used in your template
# StaticFiles="/home/user/Documents/presla-theme/static"

## Optional, define your own Executors for running code from the presentation
# CustomExecutors="/home/user/Documents/presla-executors"

## Can be set to "debug", defaults to "warning": Set log level to debug
# LogLevel="debug"

## Can be set to "json", defaults to "text": Set log format to debug
# LogFormat="json"

## Controls if presla looks for a newer version on start
# CheckUpdateOnStart=true

## Optional, can be used multiple times
## This way you can specify a template used for only one presentation
# [[Presentations]]
# PresentationName="my_presentation"
# TemplatePath="/home/user/Documents/presla-theme-my-presentation/static"
# StaticFiles="/home/user/Documents/presla-theme-my-presentation/templates"
`
}
