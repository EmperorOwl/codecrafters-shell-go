package shell

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/parser"
)

type Shell struct{}

const Prompt = "$ "

func New() *Shell {
	return &Shell{}
}

func WritePrompt(w io.Writer) {
	io.WriteString(w, Prompt)
}

func CommandNotFoundMessage(command string) string {
	return command + ": command not found"
}

func (s *Shell) Run(shellStdin io.Reader, shellStdout, shellStderr io.Writer) error {
	reader := bufio.NewReader(shellStdin)
	for {
		WritePrompt(shellStdout)

		line, err := reader.ReadString('\n')
		if err == io.EOF {
			if strings.TrimSpace(line) == "" {
				return nil
			}
		} else if err != nil {
			return err
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		fields, redirect := parser.ParseRedirect(parser.Tokenize(line))
		if len(fields) == 0 {
			if err == io.EOF {
				return nil
			}
			continue
		}

		command := fields[0]
		stdout, closeStdout, redirectErr := openRedirect(shellStdout, redirect.StdoutPath, redirect.StdoutAppend)
		if redirectErr != nil {
			return redirectErr
		}
		stderr, closeStderr, redirectErr := openRedirect(shellStderr, redirect.StderrPath, redirect.StderrAppend)
		if redirectErr != nil {
			closeStdout()
			return redirectErr
		}
		closeRedirects := func() {
			closeStdout()
			closeStderr()
		}

		if handled, shouldExit := TryBuiltin(fields, stdout, stderr); handled {
			closeRedirects()
			if shouldExit {
				return nil
			}
			if err == io.EOF {
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
			if err == io.EOF {
				return nil
			}
			continue
		}

		closeRedirects()
		fmt.Fprintf(shellStdout, "%s\n", CommandNotFoundMessage(command))

		if err == io.EOF {
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
