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

func (s *Shell) Run(in io.Reader, out io.Writer) error {
	reader := bufio.NewReader(in)
	for {
		WritePrompt(out)

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

		fields, stdoutRedirect := parser.ParseRedirect(parser.Tokenize(line))
		if len(fields) == 0 {
			if err == io.EOF {
				return nil
			}
			continue
		}

		command := fields[0]
		stdout, closeStdout, redirectErr := openStdout(out, stdoutRedirect)
		if redirectErr != nil {
			return redirectErr
		}

		if handled, shouldExit := TryBuiltin(fields, stdout); handled {
			closeStdout()
			if shouldExit {
				return nil
			}
			if err == io.EOF {
				return nil
			}
			continue
		}

		if executed, execErr := ExecuteExternalProgram(fields, stdout); executed {
			closeStdout()
			var exitErr *exec.ExitError
			if execErr != nil && !errors.As(execErr, &exitErr) {
				return execErr
			}
			if err == io.EOF {
				return nil
			}
			continue
		}

		closeStdout()
		fmt.Fprintf(out, "%s\n", CommandNotFoundMessage(command))

		if err == io.EOF {
			return nil
		}
	}
}

func openStdout(out io.Writer, path string) (io.Writer, func(), error) {
	if path == "" {
		return out, func() {}, nil
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return nil, func() {}, err
	}

	return file, func() { file.Close() }, nil
}
