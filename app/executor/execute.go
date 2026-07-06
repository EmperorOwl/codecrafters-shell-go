package executor

import (
	"errors"
	"os/exec"

	"github.com/codecrafters-io/shell-starter-go/app/builtins"
	"github.com/codecrafters-io/shell-starter-go/app/external"
)

func (e *Executor) builtinContext(outputs CommandOutputs) *builtins.Context {
	return &builtins.Context{
		Stdout:     outputs.Stdout,
		Stderr:     outputs.Stderr,
		Jobs:       e.jobTable,
		Completion: e.completionRegistry,
	}
}

// ExecuteBuiltin runs a builtin command. The bool is true when the shell should exit.
func (e *Executor) ExecuteBuiltin(fields []string, outputs CommandOutputs) (bool, error) {
	return builtins.Run(fields[0], fields[1:], e.builtinContext(outputs))
}

// ExecuteExternalForeground runs an external command and waits for it to finish.
func (e *Executor) ExecuteExternalForeground(fields []string, outputs CommandOutputs) error {
	prog, ok := external.New(fields, outputs.Stdout, outputs.Stderr)
	if !ok {
		return nil
	}
	return nonExitError(prog.Run())
}

// ExecuteExternalBackground starts an external command in the background.
// It returns the assigned job number and process ID.
func (e *Executor) ExecuteExternalBackground(fields []string, outputs CommandOutputs, line string) (int, int, error) {
	prog, ok := external.New(fields, outputs.Stdout, outputs.Stderr)
	if !ok {
		return 0, 0, nil
	}

	var jobNumber int
	pid, err := prog.RunInBackground(func() {
		e.jobTable.MarkDone(jobNumber)
	})
	if err != nil {
		return 0, 0, err
	}

	jobNumber = e.jobTable.Add(pid, line)
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
