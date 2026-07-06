package shell

import (
	"fmt"
	"io"

	"github.com/codecrafters-io/shell-starter-go/app/builtins"
	"github.com/codecrafters-io/shell-starter-go/app/completer"
	"github.com/codecrafters-io/shell-starter-go/app/executor"
	"github.com/codecrafters-io/shell-starter-go/app/external"
	"github.com/codecrafters-io/shell-starter-go/app/jobs"
	"github.com/codecrafters-io/shell-starter-go/app/parser"
	"github.com/codecrafters-io/shell-starter-go/app/repl"
	"github.com/codecrafters-io/shell-starter-go/app/terminal"
)

// Shell is the top-level orchestrator for the interactive shell.
type Shell struct {
	terminal  *terminal.Terminal
	executor  *executor.Executor
	completer *completer.Completer
	state     *repl.State
}

// New wires shell dependencies and returns a ready-to-run shell.
func New(stdin io.Reader, stdout, stderr io.Writer) *Shell {
	state := repl.NewState()

	s := &Shell{
		executor:  executor.New(stdin),
		completer: completer.New(state),
		state:     state,
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

		s.writeReapedJobs()

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
	parsed := parser.ParseLine(line)
	if parsed.Pipeline {
		return s.executePipeline(parsed)
	}
	return s.executeCommand(parsed, line)
}

func (s *Shell) executeCommand(parsed parser.Line, line string) (bool, error) {
	fields := parsed.Commands[0]
	if len(fields) == 0 {
		return false, nil
	}

	notFound, ok := commandFound(fields)
	if !ok {
		s.terminal.WriteLine(CommandNotFoundMessage(notFound))
		return false, nil
	}

	outputs := executor.Outputs{
		Stdout:   s.terminal.Stdout(),
		Stderr:   s.terminal.Stderr(),
		Redirect: parsed.Redirect,
	}

	if builtins.IsBuiltin(fields[0]) {
		exitShell, err := s.executor.ExecuteBuiltin(outputs, s.state, fields)
		if exitShell {
			return true, nil
		}
		return err != nil, err
	}

	if parsed.Background {
		return s.executeBackgroundCommand(outputs, fields, line)
	}

	if err := s.executor.ExecuteExternalForeground(outputs, fields); err != nil {
		return true, err
	}
	return false, nil
}

func (s *Shell) executeBackgroundCommand(outputs executor.Outputs, fields []string, line string) (bool, error) {
	var jobNumber int
	pid, err := s.executor.ExecuteExternalBackground(outputs, fields, func() {
		s.state.Jobs.MarkDone(jobNumber)
	})
	if err != nil {
		return true, err
	}
	if pid > 0 {
		jobNumber = s.state.Jobs.Add(pid, line)
		s.terminal.WriteLine(fmt.Sprintf("[%d] %d", jobNumber, pid))
	}
	return false, nil
}

func (s *Shell) executePipeline(parsed parser.Line) (bool, error) {
	notFound, ok := validatePipelineSegments(parsed.Commands)
	if !ok {
		if notFound != "" {
			s.terminal.WriteLine(CommandNotFoundMessage(notFound))
		}
		return false, nil
	}

	outputs := executor.Outputs{
		Stdout:   s.terminal.Stdout(),
		Stderr:   s.terminal.Stderr(),
		Redirect: parsed.Redirect,
	}

	if err := s.executor.ExecutePipeline(outputs, s.state, parsed.Commands); err != nil {
		return true, err
	}
	return false, nil
}

// HandleTab is the TabHandler entry point called by terminal on each Tab press.
func (s *Shell) HandleTab(state *terminal.TabState, buffer string) terminal.TabResult {
	return s.completer.HandleTab(state, buffer)
}

// writeReapedJobs prints any finished background jobs before the next prompt.
func (s *Shell) writeReapedJobs() {
	for _, line := range jobs.FormatLines(s.state.Jobs.ReapDone()) {
		s.terminal.WriteLine(line)
	}
}

func commandFound(fields []string) (notFound string, ok bool) {
	if len(fields) == 0 {
		return "", false
	}
	name := fields[0]
	if builtins.IsBuiltin(name) {
		return "", true
	}
	if _, found := external.FindExecutableInPath(name); found {
		return "", true
	}
	return name, false
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
