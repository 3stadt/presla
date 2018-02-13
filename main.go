package main

import (
	"flag"
	"fmt"
	"git.3stadt.com/3stadt/presla/src/Handlers"
	"git.3stadt.com/3stadt/presla/src/PresLaTemplates"
	"github.com/BurntSushi/toml"
	"github.com/fatih/color"
	"github.com/labstack/echo"
	"github.com/mitchellh/go-homedir"
	"io/ioutil"
	"os"
	"strings"
)

type Config struct {
	MarkdownPath    string
	FooterText      string
	ListenOn        string
	TemplatePath    string
	StaticFiles     string
	Presentations   []Handlers.PresentationConf
	CustomExecutors string
	LogLevel        string
	LogFormat       string
}

var conf Config

/**
Get config values from `config.toml`.
Will create the file with default values if it doesn't exist
*/
func init() {

	configPath, err := getConfPath()
	checkErr(err)

	tomlData, err := ioutil.ReadFile(configPath)
	checkErr(err)

	_, err = toml.Decode(string(tomlData), &conf)
	checkErr(err)

	conf.MarkdownPath = strings.TrimSuffix(conf.MarkdownPath, "/")
}

/**
Set up HTTP routes, start server
*/
func main() {
	handler := &Handlers.Conf{
		MarkdownPath:    conf.MarkdownPath,
		FooterText:      conf.FooterText,
		TemplatePath:    conf.TemplatePath,
		StaticFiles:     conf.StaticFiles,
		Presentations:   conf.Presentations,
		CustomExecutors: conf.CustomExecutors,
		LogLevel:        conf.LogLevel,
		LogFormat:       conf.LogFormat,
	}

	e := echo.New()

	for _, c := range conf.Presentations {
		_, err := os.Stat(c.TemplatePath)
		if c.TemplatePath != "" && c.PresentationName != "" && err == nil {
			e.Renderer = PresLaTemplates.Custom(c.TemplatePath)
		}
	}
	e.GET("/static/internal/*", handler.InternalStatic)
	e.GET("/static/:pres/*", handler.Static)
	e.GET("/favicon.ico", handler.Favicon)
	e.GET("/svg/footer-text.svg", handler.Svg)
	e.GET("/md/:file", handler.Md)
	e.GET("/md/:pres/*", handler.Assets)
	e.POST("/exec", handler.Exec)
	e.GET("/:pres", handler.Presentation)
	e.GET("/", handler.Home)
	fmt.Println()
	color.Green("Starting server at: " + color.HiBlueString("http://"+conf.ListenOn))
	color.Green("=> Use Ctrl+c to quit PresLa")
	e.Start(conf.ListenOn)
}

func checkErr(e error) {
	errorColor := color.New(color.FgWhite, color.BgRed, color.Bold)
	if e != nil {
		errorColor.Println("A critical error occured:")
		if e.Error() == "toml: cannot load TOML value of type map[string]interface {} into a Go slice" {
			errorColor.Println("You must use double brackets in your config file, like [[This]], instead of single brackets like [This]")
		} else {
			errorColor.Println(e.Error())
		}
		os.Exit(1)
	}
}

func getConfPath() (string, error) {

	var configPath string
	flag.StringVar(&configPath, "conf", "", "The path to the configuration file")
	flag.Parse()

	// By preserving user input, default config will be created on user specified if it doesn't exist yet
	if configPath == "" {
		configPath = "presla.toml"
	}

	locations := []string{configPath}

	home, err := homedir.Dir()

	// Only search other configs when we have a home directory
	if err == nil {
		locations = append(locations, home+"/.presla.toml", home+"/.config/presla.toml")
	} else {
		fmt.Println("Error searching for the config in Home directory: ")
		fmt.Println(err.Error())
	}

	for _, location := range locations {
		// If file exists, return the path immediatly
		if _, err := os.Stat(location); err == nil {
			color.Green("Using config file: " + location)
			return location, nil
		} else {
			color.Yellow("No config file at " + location + " ...")
		}
	}

	defaultConfig := []byte(getDefaultConfig())
	err = ioutil.WriteFile(configPath, defaultConfig, 0644)
	if err != nil {
		return "", err
	}
	fmt.Println("Created config file with default values: " + configPath)
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

## Optional, can be used multiple times
## This way you can specify a template used for only one presentation
# [[Presentations]]
# PresentationName="my_presentation"
# TemplatePath="/home/user/Documents/presla-theme-my-presentation/static"
# StaticFiles="/home/user/Documents/presla-theme-my-presentation/templates"
`
}
