package executor

import (
	"errors"
	"io"
	"os/exec"

	"github.com/codecrafters-io/shell-starter-go/app/builtins"
	"github.com/codecrafters-io/shell-starter-go/app/external"
)

func (e *Executor) runBuiltin(stdout, stderr io.Writer, fields []string, stdin io.Reader) (bool, error) {
	ctx := &builtins.Context{
		Stdout:     stdout,
		Stderr:     stderr,
		Jobs:       e.jobTable,
		Completion: e.completionRegistry,
	}
	if stdin != nil {
		return runDrainingStdin(fields[0], fields[1:], ctx, stdin)
	}
	return builtins.Run(fields[0], fields[1:], ctx)
}

func (e *Executor) runExternal(stdout, stderr io.Writer, fields []string, stdin io.Reader) error {
	prog, ok := external.New(fields, stdout, stderr)
	if !ok {
		return nil
	}
	prog.Stdin = stdin
	return prog.Run()
}

func (e *Executor) runExternalBackground(stdout, stderr io.Writer, fields []string, line string) (int, int, error) {
	prog, ok := external.New(fields, stdout, stderr)
	if !ok {
		return 0, 0, nil
	}
	prog.Stdin = e.stdin

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

// runDrainingStdin runs a builtin while discarding pipeline stdin in the background.
// Builtins do not read stdin, so a middle pipeline stage would leave the upstream
// pipe unread and deadlock once the buffer fills. Draining stdin unblocks writers.
func runDrainingStdin(name string, args []string, ctx *builtins.Context, stdin io.Reader) (bool, error) {
	drainDone := make(chan struct{})
	go func() {
		_, _ = io.Copy(io.Discard, stdin)
		close(drainDone)
	}()

	exitShell, err := builtins.Run(name, args, ctx)
	<-drainDone
	return exitShell, err
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
