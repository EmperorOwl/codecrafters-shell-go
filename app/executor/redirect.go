package executor

import (
	"io"
	"os"

	"github.com/codecrafters-io/shell-starter-go/app/parser"
)

// commandOutputs holds writers for command stdout and stderr, with optional cleanup.
type commandOutputs struct {
	Stdout io.Writer
	Stderr io.Writer
	close  func()
}

func (o commandOutputs) Close() {
	if o.close != nil {
		o.close()
	}
}

func (e *Executor) withRedirect(redirect parser.Redirect, fn func(commandOutputs) error) error {
	resolved, err := e.openCommandOutputs(redirect)
	if err != nil {
		return err
	}
	defer resolved.Close()
	return fn(resolved)
}

func (e *Executor) openCommandOutputs(redirect parser.Redirect) (commandOutputs, error) {
	out, closeStdout, err := openRedirect(e.stdout, redirect.StdoutPath, redirect.StdoutAppend)
	if err != nil {
		return commandOutputs{}, err
	}

	errOut, closeStderr, err := openRedirect(e.stderr, redirect.StderrPath, redirect.StderrAppend)
	if err != nil {
		closeStdout()
		return commandOutputs{}, err
	}

	return commandOutputs{
		Stdout: out,
		Stderr: errOut,
		close: func() {
			closeStdout()
			closeStderr()
		},
	}, nil
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
