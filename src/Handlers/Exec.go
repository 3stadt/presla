package Handlers

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labstack/echo"
	"github.com/robertkrimen/otto"
	_ "github.com/robertkrimen/otto/underscore"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

func (conf *Conf) Exec(c echo.Context) error {
	if strings.ToLower(conf.LogFormat) == "json" {
		log.SetFormatter(&log.JSONFormatter{})
	}
	logLevel := func() log.Level {
		if strings.ToLower(conf.LogLevel) == "debug" {
			return log.DebugLevel
		} else {
			return log.WarnLevel
		}
	}
	log.SetLevel(logLevel())

	// Bind the post Params to a struct
	code := new(Code)
	if err := c.Bind(code); err != nil {
		return err
	}

	// The created file is deleted at the end. In case the file exists on the system, rename it. Can't be ours.
	if fileExists(code.Filename) {
		newFileName := code.Filename + "_" + fmt.Sprintf("%d", makeTimestamp())
		log.WithFields(log.Fields{
			"oldName": code.Filename,
			"newName": newFileName,
		}).Warning("file exists, renaming.")
		os.Rename(code.Filename, newFileName)
	}

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
	var js []byte
	if conf.CustomExecutors != "" && fileExists(conf.CustomExecutors+"/"+code.Executor+".js") {
		log.WithFields(log.Fields{
			"executor": conf.CustomExecutors + "/" + code.Executor + ".js",
		}).Debug("using custom executor")
		js, err = ioutil.ReadFile(conf.CustomExecutors + "/" + code.Executor + ".js")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	} else {
		log.WithFields(log.Fields{
			"executor": "executors/" + code.Executor + ".js",
		}).Debug("using internal executor")
		js, err = Asset("executors/" + code.Executor + ".js")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	cmdCommands, err := code.getCmdString(string(js))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	err = code.execute(c, cmdCommands)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return err
}

func (code *Code) getCmdString(js string) (commands []CmdCommand, err error) {
	vm := otto.New()

	// make the filename/path available to the JS script
	vm.Set("codePath", code.Filename)

	// make exec function available in JS, takes all arguments and creates the command parameters cmdName and cmdString.
	// can be called multiple times, all commands will be executed in order
	vm.Set("exec", func(call otto.FunctionCall) otto.Value {
		cmd := CmdCommand{}
		cmdString := call.ArgumentList
		CurCmdName, err := cmdString[0].ToString()
		if err != nil {
			errVal, _ := otto.ToValue(err.Error())
			return errVal
		}
		cmd.cmdName = CurCmdName
		cmd.cmdArgs = []string{}
		for _, arg := range cmdString[1:] {
			val, err := arg.ToString()
			if err != nil {
				errVal, _ := otto.ToValue(err.Error())
				return errVal
			}
			cmd.cmdArgs = append(cmd.cmdArgs, val)
		}
		result, _ := vm.ToValue(true)
		commands = append(commands, cmd)
		return result
	})

	// execute the javascript code from the executor with the exposed function and var from above
	_, err = vm.Run(string(js))

	// return the command, created by the executor
	return commands, err
}

func (code *Code) execute(c echo.Context, commands []CmdCommand) (err error) {
	// Execute each command in order
	for _, command := range commands {
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

		// scanner is needed for continously getting the output
		errScanner := bufio.NewScanner(cmdErrReader)
		scanner := bufio.NewScanner(cmdReader)
		chanSend := make(chan *CmdOutput, 10) // used for sending stdout and stderr to browser without overlapping
		var wg sync.WaitGroup                 // used to keep browser connection open until all messages are sent

		// go functions run until the command is finished, sending output to the browser
		go func() error {
			for {
				text := <-chanSend
				if err := json.NewEncoder(c.Response()).Encode(text); err != nil {
					log.Warningf("error encoding output for Cmd: %#s", err.Error())
					return err
				}
				c.Response().Flush()
				time.Sleep(100 * time.Millisecond) // Used to prevent spamming the browser with responses
				log.WithFields(log.Fields{
					"stdout": text.StdOut,
					"stderr": text.StdErr,
				}).Debug("sent output to browser")
				wg.Done()
			}
			return nil
		}()

		// Capture error output and send it
		go func() error {
			for errScanner.Scan() {
				wg.Add(1)
				text := CmdOutput{
					StdErr: errScanner.Text(),
				}
				// send to browser
				chanSend <- &text
			}
			return nil
		}()

		// Capture Stdout and send it
		go func() error {
			for scanner.Scan() {
				wg.Add(1)
				text := CmdOutput{
					StdOut: scanner.Text(),
				}
				// send to browser
				chanSend <- &text
			}
			return nil
		}()

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
