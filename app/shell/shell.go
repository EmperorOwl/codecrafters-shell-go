package shell

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/codecrafters-io/shell-starter-go/app/completion"
	"github.com/codecrafters-io/shell-starter-go/app/files"
	"github.com/codecrafters-io/shell-starter-go/app/jobs"
	"github.com/codecrafters-io/shell-starter-go/app/parser"
	shellpath "github.com/codecrafters-io/shell-starter-go/app/path"
	"github.com/codecrafters-io/shell-starter-go/app/terminal"
)

type Shell struct {
	nextJobID int
	jobs      []jobs.Job
}

func New() *Shell {
	return &Shell{}
}

func CommandNotFoundMessage(command string) string {
	return command + ": command not found"
}

func (s *Shell) Run(shellStdin io.Reader, shellStdout, shellStderr io.Writer) error {
	session := terminal.NewSession(shellStdin)
	defer session.Close()

	reader := bufio.NewReader(shellStdin)
	registeredCompleters := map[string]string{}
	for {
		rawMode := session.PrepareRead()

		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		listFiles := func(dir string) []string {
			return files.ListInDir(cwd, dir)
		}
		completerFuncs := completion.BuildCompleterFuncs(registeredCompleters)
		line, eof, err := terminal.ReadLine(reader, shellStdout, rawMode, BuiltinNames(), shellpath.FindAllExecutablesInPath(), listFiles, completerFuncs)
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

		fields, redirect := parser.ParseRedirect(parser.Tokenize(line))
		fields, background := parser.StripBackground(fields)
		if len(fields) == 0 {
			if eof {
				return nil
			}
			continue
		}

		command := fields[0]
		stdout, closeStdout, redirectErr := openRedirect(shellStdout, redirect.StdoutPath, redirect.StdoutAppend)
		if redirectErr != nil {
			return redirectErr
		}
		stdout = terminal.WrapWriter(stdout, rawMode && redirect.StdoutPath == "")
		stderr, closeStderr, redirectErr := openRedirect(shellStderr, redirect.StderrPath, redirect.StderrAppend)
		if redirectErr != nil {
			closeStdout()
			return redirectErr
		}
		stderr = terminal.WrapWriter(stderr, rawMode && redirect.StderrPath == "")
		closeRedirects := func() {
			closeStdout()
			closeStderr()
		}

		if handled, shouldExit := TryBuiltin(fields, stdout, stderr, registeredCompleters, s.jobs); handled {
			closeRedirects()
			if shouldExit {
				return nil
			}
			if eof {
				return nil
			}
			continue
		}

		if background {
			executed, pid, execErr := StartExternalProgram(fields, stdout, stderr)
			closeRedirects()
			if !executed {
				fmt.Fprintf(terminal.WrapWriter(shellStdout, rawMode), "%s\n", CommandNotFoundMessage(command))
			} else {
				if execErr != nil {
					return execErr
				}
				jobNumber := jobs.AddJob(&s.jobs, &s.nextJobID, pid, line)
				fmt.Fprintf(terminal.WrapWriter(shellStdout, rawMode), "[%d] %d\n", jobNumber, pid)
			}
			if eof {
				return nil
			}
			continue
		}

		if executed, execErr := ExecuteExternalProgram(fields, stdout, stderr); executed {
			closeRedirects()
			var exitErr *exec.ExitError
			if execErr != nil && !errors.As(execErr, &exitErr) {
				return execErr
			}
			if eof {
				return nil
			}
			continue
		}

		closeRedirects()
		fmt.Fprintf(terminal.WrapWriter(shellStdout, rawMode), "%s\n", CommandNotFoundMessage(command))

		if eof {
			return nil
		}
	}
}

func openRedirect(defaultWriter io.Writer, path string, shouldAppend bool) (io.Writer, func(), error) {
	if path == "" {
		return defaultWriter, func() {}, nil
	}

	flags := os.O_CREATE | os.O_WRONLY
	if !shouldAppend {
		flags |= os.O_TRUNC
	}

	file, err := os.OpenFile(path, flags, 0644)
	if err != nil {
		return nil, func() {}, err
	}

	if shouldAppend {
		// Seek to end instead of O_APPEND: on Windows, O_APPEND files do not
		// receive writes when passed to exec.Cmd.Stderr.
		if _, err := file.Seek(0, io.SeekEnd); err != nil {
			file.Close()
			return nil, func() {}, err
		}
	}

	return file, func() { file.Close() }, nil
}
