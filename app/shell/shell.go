package shell

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/codecrafters-io/shell-starter-go/app/parser"
	shellpath "github.com/codecrafters-io/shell-starter-go/app/path"
	"github.com/codecrafters-io/shell-starter-go/app/shellio"
	"golang.org/x/term"
)

type Shell struct{}

func New() *Shell {
	return &Shell{}
}

func CommandNotFoundMessage(command string) string {
	return command + ": command not found"
}

func (s *Shell) Run(shellStdin io.Reader, shellStdout, shellStderr io.Writer) error {
	stdinFile, rawMode := shellio.TerminalStdin(shellStdin)
	if rawMode {
		oldState, err := term.MakeRaw(int(stdinFile.Fd()))
		if err != nil {
			rawMode = false
		} else {
			defer term.Restore(int(stdinFile.Fd()), oldState)
		}
	}

	reader := bufio.NewReader(shellStdin)
	for {
		line, eof, err := shellio.ReadLine(reader, shellStdout, rawMode, BuiltinNames(), shellpath.FindAllExecutablesInPath())
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
		stdout = wrapTerminalWriter(stdout, rawMode && redirect.StdoutPath == "")
		stderr, closeStderr, redirectErr := openRedirect(shellStderr, redirect.StderrPath, redirect.StderrAppend)
		if redirectErr != nil {
			closeStdout()
			return redirectErr
		}
		stderr = wrapTerminalWriter(stderr, rawMode && redirect.StderrPath == "")
		closeRedirects := func() {
			closeStdout()
			closeStderr()
		}

		if handled, shouldExit := TryBuiltin(fields, stdout, stderr); handled {
			closeRedirects()
			if shouldExit {
				return nil
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
		fmt.Fprintf(wrapTerminalWriter(shellStdout, rawMode), "%s\n", CommandNotFoundMessage(command))

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
