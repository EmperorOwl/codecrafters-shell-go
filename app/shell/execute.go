package shell

import (
	"errors"
	"fmt"
	"io"
	"os/exec"

	"github.com/codecrafters-io/shell-starter-go/app/builtins"
	"github.com/codecrafters-io/shell-starter-go/app/external"
	"github.com/codecrafters-io/shell-starter-go/app/parser"
)

type lineContext struct {
	stdout io.Writer
	stderr io.Writer
	line   string
}

type resolvedCommand struct {
	builtin  *builtins.Builtin
	external *external.ExternalProgram
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

func (s *Shell) resolveCommand(fields []string, stdout, stderr io.Writer) (resolvedCommand, string, bool) {
	if len(fields) == 0 {
		return resolvedCommand{}, "", false
	}
	if builtins.IsBuiltin(fields[0]) {
		return resolvedCommand{
			builtin: builtins.New(fields[0], fields[1:], stdout, stderr, s.completers, &s.jobs),
		}, "", true
	}
	prog, ok := external.New(fields, stdout, stderr)
	if !ok {
		return resolvedCommand{}, fields[0], false
	}
	return resolvedCommand{external: prog}, "", true
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

	cmd, notFound, ok := s.resolveCommand(fields, outputs.Stdout, outputs.Stderr)
	if !ok {
		ctx.printCommandNotFound(notFound)
		return false, nil
	}

	if cmd.builtin != nil {
		return s.executeBuiltin(cmd.builtin)
	}

	if background {
		return s.executeBackground(cmd.external, ctx)
	}
	return s.executeForeground(cmd.external)
}

func (s *Shell) executeBuiltin(cmd *builtins.Builtin) (bool, error) {
	exitShell, err := cmd.Run()
	if exitShell {
		return true, nil
	}
	if err != nil {
		return true, err
	}
	return false, nil
}

func (s *Shell) executeBackground(prog *external.ExternalProgram, ctx lineContext) (bool, error) {
	var jobNumber int
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

func (s *Shell) executeForeground(prog *external.ExternalProgram) (bool, error) {
	if err := nonExitError(prog.Run()); err != nil {
		return true, err
	}
	return false, nil
}
