// Copyright (C) 2024  Morgan S Hein
//
// This program is subject to the terms
// of the GNU Affero General Public License, version 3.
// If a copy of the AGPL was not distributed with this file, You
// can obtain one at https://www.gnu.org/licenses/.

package sh

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/mvdan/sh/shell"
)

type Executor struct {
	ctx               context.Context
	cmd               string
	args              []string
	dir               string
	env               []string
	printFinalCommand bool // if enabled, before passing the command on, print the command out to stdout
	stdOut            io.Writer
	stdErr            io.Writer
	stdIn             io.Reader
}

func Cmd(input ...string) *Executor {
	return CmdWithCtx(context.Background(), input...)
}

func CmdWithCtx(ctx context.Context, cmd ...string) *Executor {
	actualCmd := ""
	var args []string
	if len(cmd) == 1 {
		actualCmd = cmd[0]
	}
	if len(cmd) > 1 {
		actualCmd = cmd[0]
		args = cmd[1:]
	}
	return &Executor{
		ctx:  ctx,
		cmd:  actualCmd,
		args: args,
		env:  os.Environ(),
	}
}

// EnvMapToEnv converts a map of environment variables to a slice of strings
func EnvMapToEnv(env map[string]string) []string {
	var envs []string
	for k, v := range env {
		envs = append(envs, k+"="+v)
	}
	return envs
}

// Dir sets the working directory for the command
func (e *Executor) Dir(dir string) *Executor {
	e.dir = dir
	return e
}

// SetPrintFinalCommand sets the printFinalCommand flag
func (e *Executor) SetPrintFinalCommand(printFinalCommand bool) *Executor {
	e.printFinalCommand = printFinalCommand
	return e
}

// SetArgs sets the arguments for the command
func (e *Executor) SetArgs(args ...string) *Executor {
	e.args = args
	return e
}

// Stdin sets the stdin for the command
func (e *Executor) Stdin(in io.Reader) *Executor {
	e.stdIn = in
	return e
}

// SetEnv sets the environment variables for the command
func (e *Executor) SetEnv(env []string) *Executor {
	e.env = env
	return e
}

func (e *Executor) AddEnv(env []string) *Executor {
	e.env = append(e.env, env...)
	return e
}

// StdOut runs the command, and returns the stdout as a string
func (e *Executor) StdOut() (string, error) {
	var out bytes.Buffer
	e.stdOut = &out
	err := e.Run()
	return out.String(), err
}

// String runs the command, and returns the combined stdout and stderr as a string
func (e *Executor) String() (string, error) {
	out := &bytes.Buffer{}
	e.stdOut = out
	e.stdErr = out
	err := e.Run()
	return out.String(), err
}

// RunWithWriters executes the command and writes the output to the provided writers
// If stdOut or stdErr are nil, they are not used.
func (e *Executor) RunWithWriters(stdOut, errOut io.Writer) error {
	if stdOut == nil {
		stdOut = os.Stdout
	}
	if errOut == nil {
		e.stdErr = errOut
	}
	return e.Run()
}

// RunAndStream runs the command and streams the output to os.stdOut and os.StdErr
func (e *Executor) RunAndStream() error {
	e.stdOut = os.Stdout
	e.stdErr = os.Stderr
	return e.Run()
}

// Run runs the command
func (e *Executor) Run() error {
	// check if there are any arguments, if not and there are spaces in the command, perform
	// argparsing on the input and set the command and args
	if len(e.args) == 0 && strings.Contains(e.cmd, " ") {
		parts, err := shell.Fields(e.cmd, nil)
		if err != nil {
			fmt.Printf("error parsing command: %s\n", err)
		}
		if err == nil {
			e.cmd = parts[0]
			e.args = parts[1:]
		}
	}
	// if we've set some args, but the command has spaces, we need to parse the command and args and combine
	if len(e.args) > 0 && strings.Contains(e.cmd, " ") {
		parts, err := shell.Fields(e.cmd, nil)
		if err != nil {
			fmt.Printf("error parsing command: %s\n", err)
		}
		if err == nil {
			e.cmd = parts[0]
			e.args = append(parts[1:], e.args...)
		}
	}
	if e.printFinalCommand {
		fmt.Printf("Running command: %s %s\n", e.cmd, strings.Join(e.args, " "))
	}

	c := exec.CommandContext(e.ctx, e.cmd, e.args...)
	// get absolute path to the dir
	if e.dir != "" {
		absPath, err := filepath.Abs(e.dir)
		if err != nil {
			return err
		}
		c.Dir = absPath
	}
	c.Env = e.env
	c.Stdout = e.stdOut
	c.Stderr = e.stdOut
	c.Stdin = e.stdIn

	err := c.Run()
	return err
}
