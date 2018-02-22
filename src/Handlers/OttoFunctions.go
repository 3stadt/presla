package Handlers

import (
	"github.com/robertkrimen/otto"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"runtime"
)

func (oc *ottoConf) ottoPullDockerImage(call otto.FunctionCall) otto.Value {
	msg := call.ArgumentList
	image, err := msg[0].ToString()
	if err != nil {
		log.Error(err.Error())
		return otto.FalseValue()
	}
	cmd := exec.Command("docker", "pull", image)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Start()
	if err != nil {
		log.Error(err.Error())
		return otto.FalseValue()
	}
	err = cmd.Wait()
	if err != nil {
		log.Error(err.Error())
		return otto.FalseValue()
	}
	if cmd.ProcessState.Success() {
		return otto.TrueValue()
	}
	return otto.FalseValue()
}

func (oc *ottoConf) ottoCheckDockerImageInstalled(call otto.FunctionCall) otto.Value {
	msg := call.ArgumentList
	image, err := msg[0].ToString()
	if err != nil {
		log.Error(err.Error())
		return otto.FalseValue()
	}
	out, err := exec.Command("docker", "images", "-q", image).Output()
	if err != nil {
		panic(err) // kills the current execution pipeline
	}
	if len(string(out)) > 0 {
		return otto.TrueValue()
	}
	return otto.FalseValue()
}

func (oc *ottoConf) ottoSendstdErr(call otto.FunctionCall) otto.Value {
	msg := call.ArgumentList
	message, err := msg[0].ToString()
	if err != nil {
		log.Error(err.Error())
		return otto.FalseValue()
	}
	out := ottoOut{
		stdErr: message,
	}
	oc.out = append(oc.out, out)
	return otto.TrueValue()
}

func (oc *ottoConf) ottoSendstdOut(call otto.FunctionCall) otto.Value {
	msg := call.ArgumentList
	message, err := msg[0].ToString()
	if err != nil {
		log.Error(err.Error())
		return otto.FalseValue()
	}
	out := ottoOut{
		stdOut: message,
	}
	oc.out = append(oc.out, out)
	return otto.TrueValue()
}

func (oc *ottoConf) ottoCheckProgramInstalled(call otto.FunctionCall) otto.Value {
	cmdString := call.ArgumentList
	cmdName, err := cmdString[0].ToString()
	if err != nil {
		log.Error(err.Error())
		return otto.FalseValue()
	}
	bin := "which"
	if runtime.GOOS == "windows" {
		bin = "where"
	}
	cmd := exec.Command(bin, cmdName)
	cmd.Run()
	isInstalled, err := otto.ToValue(cmd.ProcessState.Success())
	if err != nil {
		log.Error(err.Error())
		return otto.FalseValue()
	}
	return isInstalled
}

func (oc *ottoConf) ottoExecQuiet(call otto.FunctionCall) otto.Value {
	cmd, err := getCommand(call)
	if err != nil {
		log.Error(err.Error())
		return otto.FalseValue()
	}
	cmd.quiet = true
	oc.out = append(oc.out, ottoOut{
		cmd: cmd,
	})
	return otto.TrueValue()
}

func (oc *ottoConf) ottoExec(call otto.FunctionCall) otto.Value {
	cmd, err := getCommand(call)
	if err != nil {
		log.Error(err.Error())
		return otto.FalseValue()
	}
	cmd.quiet = false
	oc.out = append(oc.out, ottoOut{
		cmd: cmd,
	})
	return otto.TrueValue()
}

func getCommand(call otto.FunctionCall) (CmdCommand, error) {
	cmd := CmdCommand{}
	cmdString := call.ArgumentList
	CurCmdName, err := cmdString[0].ToString()
	if err != nil {
		return CmdCommand{}, err
	}
	cmd.cmdName = CurCmdName
	cmd.cmdArgs = []string{}
	for _, arg := range cmdString[1:] {
		val, err := arg.ToString()
		if err != nil {
			return CmdCommand{}, err
		}
		cmd.cmdArgs = append(cmd.cmdArgs, val)
	}
	return cmd, nil
}
