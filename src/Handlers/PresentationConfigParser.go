package Handlers

import "errors"

var presConf = map[string]PresentationConf{}

func (conf *Conf) parseConf() {
	for _, pc := range conf.Presentations {
		if pc.PresentationName == "" {
			continue
		}
		presConf[pc.PresentationName] = pc
	}
}

func (conf *Conf) getConf(presentationName string) (*PresentationConf, error) {
	if len(presConf) == 0 {
		conf.parseConf()
	}
	if c, ok := presConf[presentationName]; ok {
		return &c, nil
	}
	return &PresentationConf{}, errors.New("No special config for: " + presentationName)
}
