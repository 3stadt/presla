package Handlers

import (
	"github.com/3stadt/presla/src/PreslaTemplates"
	"github.com/robertkrimen/otto"
	"github.com/spf13/afero"
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
	ConfigFile      string
	MarkdownPath    string
	FooterText      string
	TemplatePath    string
	StaticFiles     string
	Presentations   []PresentationConf
	DefaultTemplate *PreslaTemplates.DefaultTemplate
	CustomExecutors string
	LogLevel        string
	LogFormat       string
	Fs              afero.Fs
}

type CmdCommand struct {
	cmdName string
	cmdArgs []string
	quiet   bool
}

type CmdOutput struct {
	StdOut string `json:"stdout"`
	StdErr string `json:"stderr"`
}

type ottoOut struct {
	cmd    CmdCommand
	stdOut string
	stdErr string
}

type ottoConf struct {
	out []ottoOut
	vm  *otto.Otto
}
