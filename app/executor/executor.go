package executor

import (
	"io"

	"github.com/codecrafters-io/shell-starter-go/app/parser"
	"github.com/codecrafters-io/shell-starter-go/app/session"
)

// Executor runs parsed commands using injected I/O streams.
type Executor struct {
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
}

// New returns an executor wired to the given I/O streams.
func New(stdin io.Reader, stdout, stderr io.Writer) *Executor {
	return &Executor{stdin: stdin, stdout: stdout, stderr: stderr}
}

// ExecuteBuiltin runs a builtin command. The bool is true when the shell should exit.
func (e *Executor) ExecuteBuiltin(redirect parser.Redirect, sess *session.Session, fields []string) (bool, error) {
	var exitShell bool
	err := e.withRedirect(redirect, func(resolved commandOutputs) error {
		var err error
		exitShell, err = e.runBuiltin(resolved.Stdout, resolved.Stderr, sess, fields, nil)
		return err
	})
	return exitShell, err
}

// ExecuteExternalForeground runs an external command and waits for it to finish.
func (e *Executor) ExecuteExternalForeground(redirect parser.Redirect, fields []string) error {
	return e.withRedirect(redirect, func(resolved commandOutputs) error {
		return nonExitError(e.runExternal(resolved.Stdout, resolved.Stderr, fields, e.stdin))
	})
}

// ExecuteExternalBackground starts an external command in the background and returns its PID.
// onStarted runs synchronously after the process starts; onExit runs after it exits.
func (e *Executor) ExecuteExternalBackground(redirect parser.Redirect, fields []string, onStarted func(int), onExit func()) (int, error) {
	var pid int

	err := e.withRedirect(redirect, func(resolved commandOutputs) error {
		var err error
		pid, err = e.runExternalBackground(resolved.Stdout, resolved.Stderr, fields, onStarted, onExit)
		return err
	})
	if err != nil {
		return 0, err
	}
	return pid, nil
}

// ExecutePipeline runs a pipeline of commands connected by pipes.
func (e *Executor) ExecutePipeline(redirect parser.Redirect, sess *session.Session, segments [][]string) error {
	if len(segments) < 2 {
		return nil
	}

	return e.withRedirect(redirect, func(resolved commandOutputs) error {
		return nonExitError(e.runPipeline(segments, resolved.Stdout, resolved.Stderr, sess))
	})
}
