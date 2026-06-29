package shell

import "io"

// terminalWriter translates LF to CRLF so output starts at column 0 in raw mode.
type terminalWriter struct {
	w io.Writer
}

func (t terminalWriter) Write(p []byte) (int, error) {
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

func redrawLine(w io.Writer, line string) {
	io.WriteString(w, "\r\033[K")
	io.WriteString(w, Prompt)
	io.WriteString(w, line)
}

func wrapTerminalWriter(w io.Writer, rawMode bool) io.Writer {
	if rawMode {
		return terminalWriter{w: w}
	}
	return w
}
