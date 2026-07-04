package shell

import (
	"errors"
	"fmt"
	"io"
	"os/exec"

	"github.com/codecrafters-io/shell-starter-go/app/external"
	"github.com/codecrafters-io/shell-starter-go/app/parser"
)

type lineContext struct {
	stdout io.Writer
	stderr io.Writer
	line   string
}

func (c lineContext) printCommandNotFound(command string) {
	fmt.Fprintf(c.stdout, "%s\n", CommandNotFoundMessage(command))
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

func (s *Shell) ExecuteLine(line string, stdout, stderr io.Writer) (bool, error) {
	ctx := lineContext{stdout: stdout, stderr: stderr, line: line}
	tokens := parser.Tokenize(line)
	if segments := parser.SplitPipelineTokens(tokens); len(segments) >= 2 {
		return s.executePipeline(segments, ctx)
	}
	return s.executeCommand(tokens, ctx)
}

func (s *Shell) executeCommand(tokens []string, ctx lineContext) (bool, error) {
	fields, redirect := parser.ParseRedirect(tokens)
	fields, background := parser.StripBackground(fields)
	if len(fields) == 0 {
		return false, nil
	}

	outputs, err := openCommandOutputs(ctx.stdout, ctx.stderr, redirect)
	if err != nil {
		return true, err
	}
	defer outputs.Close()

	if handled, shouldExit := TryBuiltin(fields, outputs.Stdout, outputs.Stderr, s.completers, &s.jobs); handled {
		if shouldExit {
			return true, nil
		}
		return false, nil
	}

	if background {
		return s.executeBackground(fields, ctx, outputs)
	}
	return s.executeForeground(fields, ctx, outputs)
}

func (s *Shell) executeBackground(fields []string, ctx lineContext, outputs commandOutputs) (bool, error) {
	var jobNumber int
	prog, ok := external.New(fields, outputs.Stdout, outputs.Stderr)
	if !ok {
		ctx.printCommandNotFound(fields[0])
		return false, nil
	}

	pid, execErr := prog.RunInBackground(func() {
		s.jobs.MarkDone(jobNumber)
	})
	if execErr != nil {
		return true, execErr
	}

	jobNumber = s.jobs.Add(pid, ctx.line)
	fmt.Fprintf(ctx.stdout, "[%d] %d\n", jobNumber, pid)
	return false, nil
}

func (s *Shell) executeForeground(fields []string, ctx lineContext, outputs commandOutputs) (bool, error) {
	prog, ok := external.New(fields, outputs.Stdout, outputs.Stderr)
	if !ok {
		ctx.printCommandNotFound(fields[0])
		return false, nil
	}

	if err := nonExitError(prog.Run()); err != nil {
		return true, err
	}
	return false, nil
}
