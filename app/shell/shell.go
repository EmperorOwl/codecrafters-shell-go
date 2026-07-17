package shell

import (
	"fmt"
	"io"
	"os"

	"github.com/codecrafters-io/shell-starter-go/app/builtins"
	"github.com/codecrafters-io/shell-starter-go/app/completer"
	"github.com/codecrafters-io/shell-starter-go/app/executor"
	"github.com/codecrafters-io/shell-starter-go/app/external"
	"github.com/codecrafters-io/shell-starter-go/app/jobs"
	"github.com/codecrafters-io/shell-starter-go/app/parser"
	"github.com/codecrafters-io/shell-starter-go/app/session"
	"github.com/codecrafters-io/shell-starter-go/app/terminal"
	"github.com/codecrafters-io/shell-starter-go/app/variables"
)

// Shell is the top-level orchestrator for the interactive shell.
type Shell struct {
	terminal  *terminal.Terminal
	executor  *executor.Executor
	completer *completer.Completer
	session   *session.Session
}

// New wires shell dependencies and returns a ready-to-run shell.
func New(stdin io.Reader, stdout, stderr io.Writer) *Shell {
	sess := newSession()

	s := &Shell{
		completer: completer.New(sess),
		session:   sess,
	}
	s.terminal = terminal.New(s, s, stdin, stdout, stderr)
	s.executor = executor.New(stdin, s.terminal.Stdout(), s.terminal.Stderr())
	return s
}

// CommandNotFoundMessage formats the standard command-not-found error.
func CommandNotFoundMessage(command string) string {
	return command + ": command not found"
}

// Run executes the read-eval loop until exit or EOF.
func (s *Shell) Run() error {
	defer s.terminal.Close()
	defer func() {
		if s.session.Histfile != "" {
			_ = s.session.History.WriteToFile(s.session.Histfile)
		}
	}()

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
	s.session.History.Add(line)

	parsed := parser.ParseLine(line)
	parsed = expandParsedLine(parsed, s.session.Variables)
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

	if builtins.IsBuiltin(fields[0]) {
		exitShell, err := s.executor.ExecuteBuiltin(parsed.Redirect, s.session, fields)
		if exitShell {
			return true, nil
		}
		return false, err
	}

	if parsed.Background {
		return s.executeBackgroundCommand(parsed.Redirect, fields, line)
	}

	if err := s.executor.ExecuteExternalForeground(parsed.Redirect, fields); err != nil {
		return false, err
	}
	return false, nil
}

func (s *Shell) executeBackgroundCommand(redirect parser.Redirect, fields []string, line string) (bool, error) {
	var jobNumber int
	pid, err := s.executor.ExecuteExternalBackground(redirect, fields,
		func(pid int) {
			jobNumber = s.session.Jobs.Add(pid, line)
			s.terminal.WriteLine(fmt.Sprintf("[%d] %d", jobNumber, pid))
		},
		func() {
			s.session.Jobs.MarkDone(jobNumber)
		},
	)
	if err != nil {
		return false, err
	}
	if pid == 0 {
		return false, nil
	}
	return false, nil
}

func (s *Shell) executePipeline(parsed parser.Line) (bool, error) {
	notFound, ok := validatePipelineSegments(parsed.Commands)
	if !ok {
		if notFound != "" {
			s.terminal.WriteLine(CommandNotFoundMessage(notFound))
		} else {
			s.terminal.WriteLine("syntax error near unexpected token '|'")
		}
		return false, nil
	}

	if err := s.executor.ExecutePipeline(parsed.Redirect, s.session, parsed.Commands); err != nil {
		return false, err
	}
	return false, nil
}

// HandleTab is the TabHandler entry point called by terminal on each Tab press.
func (s *Shell) HandleTab(state *terminal.TabState, buffer string) terminal.TabResult {
	return s.completer.HandleTab(state, buffer)
}

// HistoryPrevious returns a previous command for up-arrow recall.
func (s *Shell) HistoryPrevious(stepsBack int) (string, bool) {
	return s.session.History.Previous(stepsBack)
}

func newSession() *session.Session {
	sess := session.NewSession()
	sess.Histfile = os.Getenv("HISTFILE")
	_ = sess.History.AppendFromFile(sess.Histfile)
	return sess
}

// writeReapedJobs prints any finished background jobs before the next prompt.
func (s *Shell) writeReapedJobs() {
	for _, line := range jobs.FormatLines(s.session.Jobs.ReapDone()) {
		s.terminal.WriteLine(line)
	}
}

func expandParsedLine(parsed parser.Line, store *variables.Store) parser.Line {
	for i, fields := range parsed.Commands {
		parsed.Commands[i] = variables.ExpandFields(store, fields)
	}
	parsed.Redirect.StdoutPath = variables.ExpandField(store, parsed.Redirect.StdoutPath)
	parsed.Redirect.StderrPath = variables.ExpandField(store, parsed.Redirect.StderrPath)
	return parsed
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
		if len(fields) == 0 {
			return "", false
		}
		notFound, ok := commandFound(fields)
		if !ok {
			return notFound, false
		}
	}
	return "", true
}
