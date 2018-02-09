package Handlers

import "git.3stadt.com/3stadt/presla/src/PresLaTemplates"

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
}
