package terminal

import (
	"bufio"
	"io"
	"os"
	"slices"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/completion"
	"golang.org/x/term"
)

// skipNextLF is set after CR so the LF from a Windows CRLF Enter is discarded
// on the next readLineRaw call instead of submitting an empty line.
var skipNextLF bool

// Stdin reports whether stdin is an interactive terminal.
// Raw-mode input is only enabled when both are true.
func Stdin(r io.Reader) (*os.File, bool) {
	f, ok := r.(*os.File)
	if !ok {
		return nil, false
	}
	return f, term.IsTerminal(int(f.Fd()))
}

func ReadLine(reader *bufio.Reader, w io.Writer, rawMode bool, builtins, executables, files []string) (line string, eof bool, err error) {
	if rawMode {
		return readLineRaw(reader, w, builtins, executables, files)
	}

	writePrompt(w, false)
	text, err := reader.ReadString('\n')
	if err == io.EOF {
		return strings.TrimSpace(text), true, nil
	}
	if err != nil {
		return "", false, err
	}
	return strings.TrimSpace(text), false, nil
}

func readLineRaw(reader *bufio.Reader, w io.Writer, builtins, executables, files []string) (string, bool, error) {
	writePrompt(w, true)

	var buffer []byte
	var pendingListings []string

	for {
		b, err := reader.ReadByte()
		if err == io.EOF {
			if len(buffer) == 0 {
				return "", true, nil
			}
			return string(buffer), true, nil
		}
		if err != nil {
			return "", false, err
		}

		switch b {
		case '\t': // Tab — autocomplete the current command prefix
			newBuffer, listings := completion.ApplyTab(builtins, executables, files, string(buffer))
			switch {
			case len(listings) > 0:
				if slices.Equal(pendingListings, listings) {
					pendingListings = nil
					writeListings(w, listings)
					redrawLine(w, string(buffer))
				} else {
					pendingListings = listings
					ringBell(w)
				}
			case newBuffer != string(buffer):
				pendingListings = nil
				buffer = []byte(newBuffer)
				redrawLine(w, newBuffer)
			default:
				pendingListings = nil
				ringBell(w)
			}
		case '\r': // Enter on Windows
			skipNextLF = true
			writeCRLF(w)
			return string(buffer), false, nil
		case '\n': // Enter on Unix; skip LF when it follows CR on Windows
			if skipNextLF {
				skipNextLF = false
				continue
			}
			writeCRLF(w)
			return string(buffer), false, nil
		case 127, 8: // Backspace (DEL on Unix, BS elsewhere)
			if len(buffer) > 0 {
				buffer = buffer[:len(buffer)-1]
				writeBackspace(w)
			}
		default: // Echo printable input
			buffer = append(buffer, b)
			w.Write([]byte{b})
		}
	}
}
