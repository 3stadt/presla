package Handlers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo"
	"github.com/robertkrimen/otto"
	_ "github.com/robertkrimen/otto/underscore" // allows the use of underscorejs functions in executors
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"
)

// Exec executes code on your system according to the post parameters it gets.
// This function is the reason presla should only listen on localhost
func (conf *Conf) Exec(c echo.Context) error {
	if strings.ToLower(conf.LogFormat) == "json" {
		log.SetFormatter(&log.JSONFormatter{})
	}
	logLevel := func() log.Level {
		if strings.ToLower(conf.LogLevel) == "debug" {
			return log.DebugLevel
		}
		return log.WarnLevel
	}
	log.SetLevel(logLevel())

	// Bind the post Params to a struct
	code := new(Code)
	if err := c.Bind(code); err != nil {
		return err
	}

	// The created file is deleted at the end. In case the file exists on the system, rename it. Can't be ours.
	code.renameExistingFile()

	file, err := os.Create(code.Filename)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	defer os.Remove(code.Filename)

	log.WithFields(log.Fields{
		"file": file.Name(),
	}).Debug("created file")

	fmt.Fprintf(file, code.Payload)

	log.WithFields(log.Fields{
		"payload": code.Payload,
	}).Debug("wrote to file")

	file.Close()

	// Check for custom executors first, only load internal executor when no custom on exists
	js, err := conf.getExecutor(code.Executor)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	cmdCommands, err := code.getCmdCommands(string(js))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	err = code.execute(c, conf, cmdCommands)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return err
}

func (conf *Conf) getExecutor(executorName string) (js []byte, err error) {
	if conf.CustomExecutors != "" && fileExists(conf.CustomExecutors+"/"+executorName+".js") {
		log.WithFields(log.Fields{
			"executor": conf.CustomExecutors + "/" + executorName + ".js",
		}).Debug("using custom executor")
		return ioutil.ReadFile(conf.CustomExecutors + "/" + executorName + ".js")
	}
	log.WithFields(log.Fields{
		"executor": "executors/" + executorName + ".js",
	}).Debug("using internal executor")
	return Asset("executors/" + executorName + ".js")
}

func (code *Code) renameExistingFile() {
	if fileExists(code.Filename) {
		newFileName := code.Filename + "_" + fmt.Sprintf("%d", makeTimestamp())
		log.WithFields(log.Fields{
			"oldName": code.Filename,
			"newName": newFileName,
		}).Warning("file exists, renaming.")
		os.Rename(code.Filename, newFileName)
	}
}

func (code *Code) getCmdCommands(js string) (out []ottoOut, err error) {
	vm := otto.New()

	oc := &ottoConf{
		out: out,
		vm:  vm,
	}

	// make the filename/path available to the JS script
	vm.Set("codePath", code.Filename)

	// The current OS. See https://github.com/golang/go/blob/master/src/go/build/syslist.go
	vm.Set("os", runtime.GOOS)

	// The current architecture, see https://github.com/golang/go/blob/master/src/go/build/syslist.go
	vm.Set("arch", runtime.GOARCH)

	// make exec function available in JS, takes all arguments and creates the command parameters cmdName and cmdString.
	// can be called multiple times, all commands will be executed in order of calling.
	// Execution takes place on the system presla is running on!
	vm.Set("exec", oc.ottoExec)
	vm.Set("execQuiet", oc.ottoExecQuiet)

	// Check if the given program is installed.
	// Supports Windows and any OS that knows the "which" command
	vm.Set("isInstalled", oc.ottoCheckProgramInstalled)

	// Checks if a given image, e.g. php or php:7.2 is installed on the system
	vm.Set("isDockerImageInstalled", oc.ottoCheckDockerImageInstalled)
	vm.Set("pullDockerImage", oc.ottoPullDockerImage)

	vm.Set("sendStdOut", oc.ottoSendstdOut)

	vm.Set("sendStdErr", oc.ottoSendstdErr)

	// execute the javascript code from the executor with the exposed function and var from above
	_, err = vm.Run(string(js))

	// return the command, created by the executor
	return oc.out, err
}

func (code *Code) execute(c echo.Context, conf *Conf, commands []ottoOut) (err error) {

	// Execute each command in order
	for _, out := range commands {
		command := out.cmd
		if out.stdErr != "" || out.stdOut != "" {
			update, err := json.Marshal(map[string]interface{}{
				"type":     "logupdate",
				"editorId": code.EditorID,
				"stdout":   out.stdOut,
				"stderr":   out.stdErr,
				"clear":    false,
			})
			if err != nil {
				log.Error(err.Error())
			}
			conf.SyncedEditorPub <- SyncedEditor{
				Update: string(update),
			}
			continue
		}
		log.WithFields(log.Fields{
			"binary":     command.cmdName,
			"parameters": command.cmdArgs,
		}).Debug("executing command")
		if command.cmdName == "" {
			log.Warningf("Empty command name: %#v", command)
			continue
		}

		// Set Header for chunked response
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		// prepare the command
		command.cmdArgs = append(command.cmdArgs, strings.Split(code.CmdArgs, " ")...)
		cmd := exec.Command(command.cmdName, command.cmdArgs...)

		cmdReader, err := cmd.StdoutPipe()
		if err != nil {
			log.Warningf("error creating StdoutPipe for Cmd: %s", err.Error())
			continue
		}
		cmdErrReader, err := cmd.StderrPipe()
		if err != nil {
			log.Warningf("error creating StderrPipe for Cmd: %s", err.Error())
			continue
		}

		// scanners are needed for continuously getting the output
		errScanner := bufio.NewScanner(cmdErrReader)
		outScanner := bufio.NewScanner(cmdReader)
		chanSend := make(chan *CmdOutput, 10) // used for sending stdout and stderr to browser without overlapping
		var wg sync.WaitGroup                 // used to keep browser connection open until all messages are sent

		// go functions run until the command is finished, sending output to the browser
		go code.responseLoop(chanSend, conf, &wg)

		// Capture error output and send it
		go captureScanner(command, chanSend, errScanner, &wg, true)

		// Capture Stdout and send it
		go captureScanner(command, chanSend, outScanner, &wg, false)

		// Start() runs the command and continues
		err = cmd.Start()
		if err != nil {
			log.Errorf("error while executing command: %s", err.Error())
		}

		// wait until finished
		err = cmd.Wait()
		if err != nil {
			log.Errorf("error while executing command: %s", err.Error())
		}
		wg.Wait()
	}
	return err
}

func captureScanner(command CmdCommand, chanSend chan *CmdOutput, scanner *bufio.Scanner, wg *sync.WaitGroup, isErrScanner bool) {
	for scanner.Scan() {
		if command.quiet {
			continue
		}
		wg.Add(1)
		text := CmdOutput{}
		if isErrScanner {
			text.StdErr = scanner.Text()
		} else {
			text.StdOut = scanner.Text()
		}
		// send to browser
		chanSend <- &text
	}
}

func (code *Code) responseLoop(chanSend chan *CmdOutput, conf *Conf, wg *sync.WaitGroup) {
	for {
		text := <-chanSend
		update, err := json.Marshal(map[string]interface{}{
			"type":     "logupdate",
			"editorId": code.EditorID,
			"stdout":   text.StdOut,
			"stderr":   text.StdErr,
			"clear":    false,
		})
		if err != nil {
			log.Error(err.Error())
		}
		conf.SyncedEditorPub <- SyncedEditor{
			Update: string(update),
		}
		log.WithFields(log.Fields{
			"stdout": text.StdOut,
			"stderr": text.StdErr,
		}).Debug("sent output to websockets")
		wg.Done()
	}
}

func fileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
