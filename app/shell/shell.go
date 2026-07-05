package shell

import (
	"fmt"
	"io"

	"github.com/codecrafters-io/shell-starter-go/app/builtins"
	"github.com/codecrafters-io/shell-starter-go/app/completion"
	"github.com/codecrafters-io/shell-starter-go/app/executor"
	"github.com/codecrafters-io/shell-starter-go/app/external"
	"github.com/codecrafters-io/shell-starter-go/app/jobs"
	"github.com/codecrafters-io/shell-starter-go/app/parser"
	"github.com/codecrafters-io/shell-starter-go/app/terminal"
)

// Shell is the top-level orchestrator for the interactive shell.
type Shell struct {
	terminal           *terminal.Terminal
	executor           *executor.Executor
	jobManager         *jobs.JobManager
	completionRegistry *completion.CompletionRegistry
}

// New wires shell dependencies and returns a ready-to-run shell.
func New(stdin io.Reader, stdout, stderr io.Writer) *Shell {
	jobManager := &jobs.JobManager{}
	completionRegistry := completion.NewCompletionRegistry()
	ex := executor.New(jobManager, completionRegistry)

	s := &Shell{
		executor:           ex,
		jobManager:         jobManager,
		completionRegistry: completionRegistry,
	}
	s.terminal = terminal.New(s, stdin, stdout, stderr)
	return s
}

// CommandNotFoundMessage formats the standard command-not-found error.
func CommandNotFoundMessage(command string) string {
	return command + ": command not found"
}

// Run executes the read-eval loop until exit or EOF.
func (s *Shell) Run() error {
	defer s.terminal.Close()

	for {
		s.terminal.PrepareRead()

		for _, line := range jobs.FormatLines(s.jobManager.ReapDone()) {
			s.terminal.WriteLine(line)
		}

		line, eof, err := s.terminal.ReadLine()
		if err != nil {
			return err
		}
		if eof && line == "" {
			return nil
		}
		if line == "" {
			if eof {
				return nil
			}
			continue
		}

		stop, err := s.ExecuteLine(line)
		if err != nil {
			return err
		}
		if stop || eof {
			return nil
		}
	}
}

// ExecuteLine parses and runs a single input line.
func (s *Shell) ExecuteLine(line string) (bool, error) {
	tokens := parser.Tokenize(line)
	if segments := parser.SplitPipelineTokens(tokens); len(segments) >= 2 {
		return s.executePipeline(segments)
	}
	return s.executeCommand(tokens, line)
}

func (s *Shell) executeCommand(tokens []string, line string) (bool, error) {
	fields, redirect := parser.ParseRedirect(tokens)
	fields, background := parser.StripBackground(fields)
	if len(fields) == 0 {
		return false, nil
	}

	outputs, err := s.executor.OpenCommandOutputs(s.terminal.Stdout(), s.terminal.Stderr(), redirect)
	if err != nil {
		return true, err
	}
	defer outputs.Close()

	notFound, ok := commandFound(fields)
	if !ok {
		s.terminal.WriteLine(CommandNotFoundMessage(notFound))
		return false, nil
	}

	if builtins.IsBuiltin(fields[0]) {
		exitShell, err := s.executor.ExecuteBuiltin(fields, outputs)
		if exitShell {
			return true, nil
		}
		return err != nil, err
	}

	if background {
		jobNumber, pid, err := s.executor.ExecuteExternalBackground(fields, outputs, line)
		if err != nil {
			return true, err
		}
		if jobNumber > 0 {
			s.terminal.WriteLine(fmt.Sprintf("[%d] %d", jobNumber, pid))
		}
		return false, nil
	}

	if err := s.executor.ExecuteExternalForeground(fields, outputs); err != nil {
		return true, err
	}
	return false, nil
}

func (s *Shell) executePipeline(segments [][]string) (bool, error) {
	commands, redirect := executor.ParsePipelineSegments(segments)

	notFound, ok := validatePipelineSegments(commands)
	if !ok {
		if notFound != "" {
			s.terminal.WriteLine(CommandNotFoundMessage(notFound))
		}
		return false, nil
	}

	outputs, err := s.executor.OpenCommandOutputs(s.terminal.Stdout(), s.terminal.Stderr(), redirect)
	if err != nil {
		return true, err
	}
	defer outputs.Close()

	if err := s.executor.ExecutePipeline(commands, outputs); err != nil {
		return true, err
	}
	return false, nil
}

func commandFound(fields []string) (notFound string, ok bool) {
	if len(fields) == 0 {
		return "", false
	}
	if builtins.IsBuiltin(fields[0]) {
		return "", true
	}
	if _, found := external.FindExecutableInPath(fields[0]); found {
		return "", true
	}
	return fields[0], false
}

func validatePipelineSegments(segments [][]string) (notFound string, ok bool) {
	for _, fields := range segments {
		notFound, ok := commandFound(fields)
		if !ok {
			return notFound, false
		}
	}
	return "", true
}
