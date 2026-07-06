package executor

import (
	"errors"
	"io"
	"os/exec"

	"github.com/codecrafters-io/shell-starter-go/app/builtins"
	"github.com/codecrafters-io/shell-starter-go/app/external"
	"github.com/codecrafters-io/shell-starter-go/app/parser"
)

func (e *Executor) builtinContext(outputs commandOutputs) *builtins.Context {
	return &builtins.Context{
		Stdout:     outputs.Stdout,
		Stderr:     outputs.Stderr,
		Jobs:       e.jobTable,
		Completion: e.completionRegistry,
	}
}

// ExecuteBuiltin runs a builtin command. The bool is true when the shell should exit.
func (e *Executor) ExecuteBuiltin(stdout, stderr io.Writer, fields []string, redirect parser.Redirect) (bool, error) {
	var exitShell bool
	err := e.withOutputs(stdout, stderr, redirect, func(outputs commandOutputs) error {
		var err error
		exitShell, err = builtins.Run(fields[0], fields[1:], e.builtinContext(outputs))
		return err
	})
	return exitShell, err
}

// ExecuteExternalForeground runs an external command and waits for it to finish.
func (e *Executor) ExecuteExternalForeground(stdout, stderr io.Writer, fields []string, redirect parser.Redirect) error {
	return e.withOutputs(stdout, stderr, redirect, func(outputs commandOutputs) error {
		prog, ok := external.New(fields, outputs.Stdout, outputs.Stderr)
		if !ok {
			return nil
		}
		prog.Stdin = e.stdin
		return nonExitError(prog.Run())
	})
}

// ExecuteExternalBackground starts an external command in the background.
// It returns the assigned job number and process ID.
func (e *Executor) ExecuteExternalBackground(stdout, stderr io.Writer, fields []string, redirect parser.Redirect, line string) (int, int, error) {
	var jobNumber int
	var pid int

	err := e.withOutputs(stdout, stderr, redirect, func(outputs commandOutputs) error {
		prog, ok := external.New(fields, outputs.Stdout, outputs.Stderr)
		if !ok {
			return nil
		}
		prog.Stdin = e.stdin

		var err error
		pid, err = prog.RunInBackground(func() {
			e.jobTable.MarkDone(jobNumber)
		})
		if err != nil {
			return err
		}

		jobNumber = e.jobTable.Add(pid, line)
		return nil
	})
	if err != nil {
		return 0, 0, err
	}
	return jobNumber, pid, nil
}

func nonExitError(err error) error {
	if err == nil {
		return nil
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return nil
	}
	return err
}
