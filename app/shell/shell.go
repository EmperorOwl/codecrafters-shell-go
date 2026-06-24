package shell

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type Shell struct{}

func New() *Shell {
	return &Shell{}
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

		command := strings.Fields(line)[0]
		fmt.Fprintf(out, "%s\n", CommandNotFoundMessage(command))

		if err == io.EOF {
			return nil
		}
	}
}
