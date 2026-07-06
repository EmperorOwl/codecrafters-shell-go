package terminal

import (
	"io"
	"os"

	"golang.org/x/term"
)

// RawMode manages TTY raw mode for interactive input.
type RawMode struct {
	stdin    *os.File
	active   bool
	oldState *term.State
}

// NewRawMode enables raw mode when stdin is an interactive terminal.
func NewRawMode(stdin io.Reader) *RawMode {
	f, ok := stdin.(*os.File)
	if !ok || !term.IsTerminal(int(f.Fd())) {
		return &RawMode{}
	}

	oldState, err := term.MakeRaw(int(f.Fd()))
	if err != nil {
		return &RawMode{stdin: f}
	}

	return &RawMode{
		stdin:  f,
		active: true,
		oldState: oldState,
	}
}

// Close restores the terminal to its original state.
func (r *RawMode) Close() error {
	if !r.active || r.stdin == nil || r.oldState == nil {
		return nil
	}
	r.active = false
	return term.Restore(int(r.stdin.Fd()), r.oldState)
}

// PrepareRead re-enables raw mode before reading the next prompt.
//
// External programs inherit stdin and may restore cooked (line-buffered) mode
// when they exit—especially on Windows. In cooked mode Tab is handled by the
// console instead of our completion logic, so tab appears to insert whitespace.
// Re-applying raw mode before each read keeps completion working after exec.
func (r *RawMode) PrepareRead() bool {
	if !r.active || r.stdin == nil {
		return false
	}
	if _, err := term.MakeRaw(int(r.stdin.Fd())); err != nil {
		r.active = false
		return false
	}
	return true
}

// Active reports whether interactive raw-mode input is active.
func (r *RawMode) Active() bool {
	return r.active
}
