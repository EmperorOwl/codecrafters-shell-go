package terminal

import "io"

// lfWriter translates LF to CRLF for output written while raw mode is active.
//
// In cooked mode the terminal driver turns \n into \r\n automatically, so each
// new line starts at column 0. Raw mode disables that processing: \n moves the
// cursor down but leaves the column unchanged. The next line then prints where
// the previous one ended, which looks like leading spaces. Writing \r\n 
// returns to column 0 first.
type lfWriter struct {
	w io.Writer
}

func (t lfWriter) Write(p []byte) (int, error) {
	n := len(p)
	for i, b := range p {
		if b == '\n' && (i == 0 || p[i-1] != '\r') {
			if _, err := t.w.Write([]byte{'\r', '\n'}); err != nil {
				return 0, err
			}
			continue
		}
		if _, err := t.w.Write([]byte{b}); err != nil {
			return 0, err
		}
	}
	return n, nil
}

// WrapWriter returns a writer that translates LF to CRLF when rawMode is true.
func WrapWriter(w io.Writer, rawMode bool) io.Writer {
	if rawMode {
		return lfWriter{w: w}
	}
	return w
}

// modeAwareWriter delegates each Write to WrapWriter using the terminal's
// current raw-mode flag so command output stays aligned after PrepareRead.
type modeAwareWriter struct {
	term   *Terminal
	stderr bool
}

func (m modeAwareWriter) Write(p []byte) (int, error) {
	base := m.term.stdout
	if m.stderr {
		base = m.term.stderr
	}
	return WrapWriter(base, m.term.rawMode).Write(p)
}
