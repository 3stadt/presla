package Handlers

import (
	"github.com/3stadt/presla/src/PreslaTemplates"
	"github.com/gorilla/websocket"
	"github.com/robertkrimen/otto"
	"github.com/spf13/afero"
)

// SyncedEditor is used in edtor synchronisation to update editors
type SyncedEditor struct {
	Update string
}

// SyncedEditorWriter is used in edtor synchronisation to write information from a SyncedEditor to a websocket
type SyncedEditorWriter struct {
	Ws     *websocket.Conn
	Writer func(e SyncedEditor, ws *websocket.Conn) error
}

// Code is the struct the form post data for an execution request is unmarshaled to
type Code struct {
	EditorID int    `form:"editorId"`
	Executor string `form:"executor"`
	Filename string `form:"filename"`
	Payload  string `form:"payload"`
	CmdArgs  string `form:"cmdargs"`
}

// PresentationConf holds the config for a specific presentation
type PresentationConf struct {
	PresentationName string
	TemplatePath     string
	StaticFiles      string
}

// Conf is a redefinition of the main config, enhanced with a filesystem and the editor syncing
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

// CmdCommand holds all information needed to execute a command via `exec.Command()`
type CmdCommand struct {
	cmdName string
	cmdArgs []string
	quiet   bool
}

// CmdOutput holds one line of cmd output from a command
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
