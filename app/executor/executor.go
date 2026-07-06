package executor

import (
	"io"

	"github.com/codecrafters-io/shell-starter-go/app/completion"
	"github.com/codecrafters-io/shell-starter-go/app/jobs"
	"github.com/codecrafters-io/shell-starter-go/app/parser"
)

// Executor runs parsed commands using injected shell state.
type Executor struct {
	jobTable           *jobs.JobTable
	completionRegistry *completion.CompletionRegistry
	stdin              io.Reader
}

// New returns an executor wired to the given shell state and stdin stream.
func New(jobTable *jobs.JobTable, completionRegistry *completion.CompletionRegistry, stdin io.Reader) *Executor {
	return &Executor{
		jobTable:           jobTable,
		completionRegistry: completionRegistry,
		stdin:              stdin,
	}
}

// Outputs configures default stdout, stderr, and redirects for command execution.
type Outputs struct {
	Stdout   io.Writer
	Stderr   io.Writer
	Redirect parser.Redirect
}

// ExecuteBuiltin runs a builtin command. The bool is true when the shell should exit.
func (e *Executor) ExecuteBuiltin(outputs Outputs, fields []string) (bool, error) {
	var exitShell bool
	err := e.withOutputs(outputs, func(resolved commandOutputs) error {
		var err error
		exitShell, err = e.runBuiltin(resolved.Stdout, resolved.Stderr, fields, nil)
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

// ExecuteExternalBackground starts an external command in the background.
// It returns the assigned job number and process ID.
func (e *Executor) ExecuteExternalBackground(outputs Outputs, fields []string, line string) (int, int, error) {
	var jobNumber int
	var pid int

	err := e.withOutputs(outputs, func(resolved commandOutputs) error {
		var err error
		jobNumber, pid, err = e.runExternalBackground(resolved.Stdout, resolved.Stderr, fields, line)
		return err
	})
	if err != nil {
		return 0, 0, err
	}
	return jobNumber, pid, nil
}

// ExecutePipeline runs a pipeline of commands connected by pipes.
func (e *Executor) ExecutePipeline(outputs Outputs, segments [][]string) error {
	if len(segments) < 2 {
		return nil
	}

	return e.withOutputs(outputs, func(resolved commandOutputs) error {
		return nonExitError(e.runPipeline(segments, resolved.Stdout, resolved.Stderr))
	})
}
