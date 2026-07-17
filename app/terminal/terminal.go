package terminal

import (
	"bufio"
	"fmt"
	"io"
)

// Terminal handles shell I/O.
type Terminal struct {
	tabHandler     TabHandler
	historyHandler HistoryHandler
	stdout         io.Writer
	stderr         io.Writer
	rawTTY         *RawMode
	reader         *bufio.Reader
	rawMode        bool
}

// New returns a terminal wired to the given streams and input handlers.
func New(tabHandler TabHandler, historyHandler HistoryHandler, stdin io.Reader, stdout, stderr io.Writer) *Terminal {
	return &Terminal{
		tabHandler:     tabHandler,
		historyHandler: historyHandler,
		stdout:         stdout,
		stderr:         stderr,
		rawTTY:         NewRawMode(stdin),
		reader:         bufio.NewReader(stdin),
	}
}

// ReadLine reads a line after the prompt has already been displayed.
func (t *Terminal) ReadLine() (line string, eof bool, err error) {
	stdout := WrapWriter(t.stdout, t.rawMode)
	writePrompt(t.stdout, t.rawMode)
	return readLine(t.reader, stdout, t.rawMode, t.tabHandler, t.historyHandler)
}

// WriteLine writes a single line to stdout, including a trailing newline.
func (t *Terminal) WriteLine(text string) {
	stdout := WrapWriter(t.stdout, t.rawMode)
	fmt.Fprintln(stdout, text)
}

// Stdout returns the stdout writer for command output, with raw-mode wrapping when active.
func (t *Terminal) Stdout() io.Writer {
	return modeAwareWriter{term: t}
}

// Stderr returns the stderr writer for command output, with raw-mode wrapping when active.
func (t *Terminal) Stderr() io.Writer {
	return modeAwareWriter{term: t, stderr: true}
}

// PrepareRead re-enables raw mode before the next prompt.
func (t *Terminal) PrepareRead() bool {
	t.rawMode = t.rawTTY.PrepareRead()
	return t.rawMode
}

// Close restores the terminal to cooked mode.
func (t *Terminal) Close() error {
	return t.rawTTY.Close()
}
