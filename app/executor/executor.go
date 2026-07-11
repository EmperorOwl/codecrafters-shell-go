package executor

import (
	"io"

	"github.com/codecrafters-io/shell-starter-go/app/parser"
	"github.com/codecrafters-io/shell-starter-go/app/session"
)

// Executor runs parsed commands using injected I/O streams.
type Executor struct {
	stdin io.Reader
}

// New returns an executor wired to the given stdin stream.
func New(stdin io.Reader) *Executor {
	return &Executor{stdin: stdin}
}

// Outputs configures default stdout, stderr, and redirects for command execution.
type Outputs struct {
	Stdout   io.Writer
	Stderr   io.Writer
	Redirect parser.Redirect
}

// ExecuteBuiltin runs a builtin command. The bool is true when the shell should exit.
func (e *Executor) ExecuteBuiltin(outputs Outputs, state *session.State, fields []string) (bool, error) {
	var exitShell bool
	err := e.withOutputs(outputs, func(resolved commandOutputs) error {
		var err error
		exitShell, err = e.runBuiltin(resolved.Stdout, resolved.Stderr, state, fields, nil)
		return err
	})
	return exitShell, err
}

// ExecuteExternalForeground runs an external command and waits for it to finish.
func (e *Executor) ExecuteExternalForeground(outputs Outputs, fields []string) error {
	return e.withOutputs(outputs, func(resolved commandOutputs) error {
		return nonExitError(e.runExternal(resolved.Stdout, resolved.Stderr, fields, e.stdin))
	})
}

// ExecuteExternalBackground starts an external command in the background and returns its PID.
func (e *Executor) ExecuteExternalBackground(outputs Outputs, fields []string, onExit func()) (int, error) {
	var pid int

	err := e.withOutputs(outputs, func(resolved commandOutputs) error {
		var err error
		pid, err = e.runExternalBackground(resolved.Stdout, resolved.Stderr, fields, onExit)
		return err
	})
	if err != nil {
		return 0, err
	}
	return pid, nil
}

// ExecutePipeline runs a pipeline of commands connected by pipes.
func (e *Executor) ExecutePipeline(outputs Outputs, state *session.State, segments [][]string) error {
	if len(segments) < 2 {
		return nil
	}

	return e.withOutputs(outputs, func(resolved commandOutputs) error {
		return nonExitError(e.runPipeline(segments, resolved.Stdout, resolved.Stderr, state))
	})
}
