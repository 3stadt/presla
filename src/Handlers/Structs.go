package Handlers

import (
	"github.com/3stadt/presla/src/PresLaTemplates"
)

type Code struct {
	Executor string `form:"executor"`
	Filename string `form:"filename"`
	Payload  string `form:"payload"`
}

type PresentationConf struct {
	PresentationName string
	TemplatePath     string
	StaticFiles      string
}

type Conf struct {
	MarkdownPath    string
	FooterText      string
	TemplatePath    string
	StaticFiles     string
	Presentations   []PresentationConf
	DefaultTemplate *PresLaTemplates.DefaultTemplate
	CustomExecutors string
	LogLevel        string
	LogFormat       string
}

type CmdCommand struct {
	cmdName string
	cmdArgs []string
}

type CmdOutput struct {
	StdOut string `json:"stdout"`
	StdErr string `json:"stderr"`
}
