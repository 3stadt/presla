package Handlers

import (
	"github.com/3stadt/presla/src/PreslaTemplates"
	"github.com/gorilla/websocket"
	"github.com/robertkrimen/otto"
	"github.com/spf13/afero"
)

type SyncedEditor struct {
	Update string
}

type SyncedEditorWriter struct {
	Ws     *websocket.Conn
	Writer func(e SyncedEditor, ws *websocket.Conn) error
}

type Code struct {
	EditorId int    `form:"editorId"`
	Executor string `form:"executor"`
	Filename string `form:"filename"`
	Payload  string `form:"payload"`
	CmdArgs  string `form:"cmdargs"`
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
	SyncedEditorPub chan SyncedEditor
	SyncedEditorSub map[int]*SyncedEditorWriter
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
