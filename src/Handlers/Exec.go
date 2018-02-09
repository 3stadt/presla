package Handlers

import (
	"github.com/labstack/echo"
	"net/http"
	"github.com/robertkrimen/otto"
	_ "github.com/robertkrimen/otto/underscore"
	"fmt"
	"os"
	"bufio"
	"os/exec"
	"encoding/json"
	"io/ioutil"
	"time"
	"errors"
	"sync"
)

func (conf *Conf) Exec(c echo.Context) error {

	// Bind the post Params to a struct
	code := new(Code)
	if err := c.Bind(code); err != nil {
		return err
	}

	// The created file is deleted at the end. In case the file exists on the system, rename it. Can't be ours.
	if fileExists(code.Filename) {
		os.Rename(code.Filename, code.Filename+"_"+fmt.Sprintf("%d", makeTimestamp()))
	}

	file, err := os.Create(code.Filename)
	if err != nil {
		return err
	}
	fmt.Fprintf(file, code.Payload)
	file.Close()

	// Check for custom executors first, only load internal executor when no custom on exists
	var js []byte
	if conf.CustomExecutors != "" && fileExists(conf.CustomExecutors+"/"+code.Executor+".js") {
		js, err = ioutil.ReadFile(conf.CustomExecutors + "/" + code.Executor + ".js")
		if err != nil {
			return err
		}
	} else {
		js, err = Asset("executors/" + code.Executor + ".js")
		if err != nil {
			return err
		}
	}

	cmdName, cmdArgs, err := code.getCmdString(string(js))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	err = code.execute(c, cmdName, cmdArgs)
	os.Remove(code.Filename)
	return err
}

func (code *Code) getCmdString(js string) (cmdName string, cmdArgs []string, err error) {
	vm := otto.New()

	// make the filename/path available to the JS script
	vm.Set("codePath", code.Filename)

	// make exec function available in JS, takes all arguments and creates the command parameters cmdName and cmdString
	vm.Set("exec", func(call otto.FunctionCall) otto.Value {
		cmdString := call.ArgumentList
		cmdName, err = cmdString[0].ToString()
		if err != nil {
			errVal, _ := otto.ToValue(err.Error())
			return errVal
		}
		for _, arg := range cmdString[1:] {
			val, err := arg.ToString()
			if err != nil {
				errVal, _ := otto.ToValue(err.Error())
				return errVal
			}
			cmdArgs = append(cmdArgs, val)
		}
		result, _ := vm.ToValue(true)

		return result
	})

	// execute the javascript code from the executor with the exposed function and var from above
	_, err = vm.Run(string(js))

	// return the command, created by the executor
	return cmdName, cmdArgs, err
}

func (code *Code) execute(c echo.Context, cmdName string, cmdArgs []string) error {
	if cmdName == "" {
		return errors.New("empty command name")
	}

	// Set Header for chunked response
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.Response().WriteHeader(http.StatusOK)

	// prepare the command
	cmd := exec.Command(cmdName, cmdArgs...)
	cmdReader, err := cmd.StdoutPipe()
	cmdErrReader, _ := cmd.StderrPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for Cmd", err)
		return err
	}

	// scanner is needed for continously getting the output
	errScanner := bufio.NewScanner(cmdErrReader)
	scanner := bufio.NewScanner(cmdReader)
	chanSend := make(chan CmdOutput, 10) // used for sending stdout and stderr to browser without overlapping
	var wg sync.WaitGroup // used to keep browser connection open until all messages are sent
	// go functions run until the command is finished, sending output to the browser

	go func() error {
		defer wg.Done()
		for {
			text := <-chanSend
			if err := json.NewEncoder(c.Response()).Encode(text); err != nil {
				fmt.Fprintln(os.Stderr, "Error reading error output for Cmd", err)
				return err
			}
			time.Sleep(100 * time.Millisecond)
			c.Response().Flush()
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
			chanSend <- text
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
			chanSend <- text
		}
		return nil
	}()

	// Start() runs the command and continues
	err = cmd.Start()
	if err != nil {
		return err
	}

	// wait until finished
	err = cmd.Wait()
	if err != nil {
		return err
	}
	wg.Wait()
	return nil
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
