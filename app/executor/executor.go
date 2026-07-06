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

// New returns an executor wired to the given shell state.
func New(jobTable *jobs.JobTable, completionRegistry *completion.CompletionRegistry) *Executor {
	return &Executor{
		jobTable:           jobTable,
		completionRegistry: completionRegistry,
	}
}

// SetStdin configures stdin used for foreground command execution.
func (e *Executor) SetStdin(stdin io.Reader) {
	e.stdin = stdin
}

// ExecuteBuiltin runs a builtin command. The bool is true when the shell should exit.
func (e *Executor) ExecuteBuiltin(stdout, stderr io.Writer, fields []string, redirect parser.Redirect) (bool, error) {
	var exitShell bool
	err := e.withOutputs(stdout, stderr, redirect, func(outputs commandOutputs) error {
		var err error
		exitShell, err = e.runBuiltin(outputs.Stdout, outputs.Stderr, fields, nil)
		return err
	})
	return exitShell, err
}

// ExecuteExternalForeground runs an external command and waits for it to finish.
func (e *Executor) ExecuteExternalForeground(stdout, stderr io.Writer, fields []string, redirect parser.Redirect) error {
	return e.withOutputs(stdout, stderr, redirect, func(outputs commandOutputs) error {
		return nonExitError(e.runExternal(outputs.Stdout, outputs.Stderr, fields, e.stdin))
	})
}

// ExecuteExternalBackground starts an external command in the background.
// It returns the assigned job number and process ID.
func (e *Executor) ExecuteExternalBackground(stdout, stderr io.Writer, fields []string, redirect parser.Redirect, line string) (int, int, error) {
	var jobNumber int
	var pid int

	err := e.withOutputs(stdout, stderr, redirect, func(outputs commandOutputs) error {
		var err error
		jobNumber, pid, err = e.runExternalBackground(outputs.Stdout, outputs.Stderr, fields, line)
		return err
	})
	if err != nil {
		return 0, 0, err
	}
	return jobNumber, pid, nil
}
