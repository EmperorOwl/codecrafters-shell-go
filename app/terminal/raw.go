package terminal

import (
	"io"
	"os"

	"golang.org/x/term"
)

// Session manages raw terminal mode for interactive input.
type Session struct {
	stdin    *os.File
	rawMode  bool
	oldState *term.State
}

// NewSession enables raw mode when stdin is an interactive terminal.
func NewSession(stdin io.Reader) *Session {
	f, ok := stdin.(*os.File)
	if !ok || !term.IsTerminal(int(f.Fd())) {
		return &Session{}
	}

	oldState, err := term.MakeRaw(int(f.Fd()))
	if err != nil {
		return &Session{stdin: f}
	}

	return &Session{
		stdin:    f,
		rawMode:  true,
		oldState: oldState,
	}
}

// Close restores the terminal to its original state.
func (s *Session) Close() error {
	if !s.rawMode || s.stdin == nil || s.oldState == nil {
		return nil
	}
	s.rawMode = false
	return term.Restore(int(s.stdin.Fd()), s.oldState)
}

// PrepareRead re-enables raw mode before reading the next prompt.
//
// External programs inherit stdin and may restore cooked (line-buffered) mode
// when they exit—especially on Windows. In cooked mode Tab is handled by the
// console instead of our completion logic, so tab appears to insert whitespace.
// Re-applying raw mode before each read keeps completion working after exec.
func (s *Session) PrepareRead() bool {
	if !s.rawMode || s.stdin == nil {
		return false
	}
	if _, err := term.MakeRaw(int(s.stdin.Fd())); err != nil {
		s.rawMode = false
		return false
	}
	return true
}

// RawMode reports whether interactive raw-mode input is active.
func (s *Session) RawMode() bool {
	return s.rawMode
}
