package shell

import (
	"errors"
	"fmt"
	"io"
	"os/exec"

	"github.com/codecrafters-io/shell-starter-go/app/parser"
)

type lineContext struct {
	stdout io.Writer
	stderr io.Writer
	line   string
	eof    bool
}

func (c lineContext) stopAfter(err error) (bool, error) {
	if err != nil {
		return true, err
	}
	if c.eof {
		return true, nil
	}
	return false, nil
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

func (s *Shell) ExecuteLine(line string, eof bool, stdout, stderr io.Writer) (bool, error) {
	ctx := lineContext{stdout: stdout, stderr: stderr, line: line, eof: eof}
	tokens := parser.Tokenize(line)
	if segments := parser.SplitPipelineTokens(tokens); len(segments) == 2 {
		return s.executePipeline(segments, ctx)
	}
	return s.executeCommand(tokens, ctx)
}

func parsePipelineSegments(segments [][]string) ([2][]string, parser.Redirect) {
	fields0, _ := parser.ParseRedirect(segments[0])
	fields0, _ = parser.StripBackground(fields0)
	fields1, redirect := parser.ParseRedirect(segments[1])
	fields1, _ = parser.StripBackground(fields1)
	return [2][]string{fields0, fields1}, redirect
}

func (s *Shell) executePipeline(segments [][]string, ctx lineContext) (bool, error) {
	commands, redirect := parsePipelineSegments(segments)
	outputs, err := openCommandOutputs(ctx.stdout, ctx.stderr, redirect)
	if err != nil {
		return true, err
	}
	defer outputs.Close()

	executed, notFound, execErr := ExecutePipeline(commands, outputs.Stdout, outputs.Stderr)
	if !executed {
		ctx.printCommandNotFound(notFound)
	} else if err := nonExitError(execErr); err != nil {
		return true, err
	}
	return ctx.stopAfter(nil)
}

func (s *Shell) executeCommand(tokens []string, ctx lineContext) (bool, error) {
	fields, redirect := parser.ParseRedirect(tokens)
	fields, background := parser.StripBackground(fields)
	if len(fields) == 0 {
		return ctx.stopAfter(nil)
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
		return ctx.stopAfter(nil)
	}

	if background {
		return s.executeBackground(fields, ctx, outputs)
	}
	return s.executeForeground(fields, ctx, outputs)
}

func (s *Shell) executeBackground(fields []string, ctx lineContext, outputs commandOutputs) (bool, error) {
	executed, pid, cmd, execErr := StartExternalProgram(fields, outputs.Stdout, outputs.Stderr)
	if !executed {
		ctx.printCommandNotFound(fields[0])
		return ctx.stopAfter(nil)
	}
	if execErr != nil {
		return true, execErr
	}

	jobNumber := s.jobs.Add(pid, ctx.line)
	startBackgroundWait(cmd, func() {
		s.jobs.MarkDone(jobNumber)
	})
	fmt.Fprintf(ctx.stdout, "[%d] %d\n", jobNumber, pid)
	return ctx.stopAfter(nil)
}

func (s *Shell) executeForeground(fields []string, ctx lineContext, outputs commandOutputs) (bool, error) {
	executed, execErr := ExecuteExternalProgram(fields, outputs.Stdout, outputs.Stderr)
	if executed {
		if err := nonExitError(execErr); err != nil {
			return true, err
		}
		return ctx.stopAfter(nil)
	}

	ctx.printCommandNotFound(fields[0])
	return ctx.stopAfter(nil)
}
