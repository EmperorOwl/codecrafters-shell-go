package terminal

import (
	"bufio"
	"io"

	"github.com/codecrafters-io/shell-starter-go/app/shell"
)

// Terminal handles shell I/O; Shell handles command logic.
type Terminal struct {
	shell   *shell.Shell
	stdin   io.Reader
	stdout  io.Writer
	stderr  io.Writer
	session *Session
	reader  *bufio.Reader
}

func New(sh *shell.Shell, stdin io.Reader, stdout, stderr io.Writer) *Terminal {
	return &Terminal{
		shell:   sh,
		stdin:   stdin,
		stdout:  stdout,
		stderr:  stderr,
		session: NewSession(stdin),
		reader:  bufio.NewReader(stdin),
	}
}

func (t *Terminal) Run() error {
	defer t.session.Close()

	for {
		rawMode := t.session.PrepareRead()
		stdout := WrapWriter(t.stdout, rawMode)
		stderr := WrapWriter(t.stderr, rawMode)

		t.shell.PrintReapedJobs(stdout)

		writePrompt(t.stdout, rawMode)
		line, eof, err := ReadLine(t.reader, stdout, rawMode, t.shell)
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

		stop, err := t.shell.ExecuteLine(line, stdout, stderr)
		if err != nil {
			return err
		}
		if stop || eof {
			return nil
		}
	}
}
