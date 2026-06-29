package shell

import (
	"bufio"
	"io"
	"os"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/completion"
	"golang.org/x/term"
)

// skipNextLF is set after CR so the LF from a Windows CRLF Enter is discarded
// on the next readLineRaw call instead of submitting an empty line.
var skipNextLF bool

const Prompt = "$ "

func WritePrompt(w io.Writer) {
	io.WriteString(w, Prompt)
}

// terminalStdin reports whether stdin is an interactive terminal.
// Raw-mode input is only enabled when both are true.
func terminalStdin(r io.Reader) (*os.File, bool) {
	f, ok := r.(*os.File)
	if !ok {
		return nil, false
	}
	return f, term.IsTerminal(int(f.Fd()))
}

func readLine(reader *bufio.Reader, w io.Writer, rawMode bool) (line string, eof bool, err error) {
	if rawMode {
		return readLineRaw(reader, w)
	}

	WritePrompt(w)
	text, err := reader.ReadString('\n')
	if err == io.EOF {
		return strings.TrimSpace(text), true, nil
	}
	if err != nil {
		return "", false, err
	}
	return strings.TrimSpace(text), false, nil
}

func readLineRaw(reader *bufio.Reader, w io.Writer) (string, bool, error) {
	io.WriteString(w, "\r$ ")

	var buffer []byte

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
			newBuffer, listings := completion.ApplyTab(BuiltinNames(), string(buffer))
			if len(listings) > 0 {
				for _, match := range listings {
					io.WriteString(w, "\r\n")
					io.WriteString(w, match)
				}
				io.WriteString(w, "\r\n")
				redrawLine(w, string(buffer))
			} else {
				buffer = []byte(newBuffer)
				redrawLine(w, newBuffer)
			}
		case '\r': // Enter on Windows
			skipNextLF = true
			io.WriteString(w, "\r\n")
			return string(buffer), false, nil
		case '\n': // Enter on Unix; skip LF when it follows CR on Windows
			if skipNextLF {
				skipNextLF = false
				continue
			}
			io.WriteString(w, "\r\n")
			return string(buffer), false, nil
		case 127, 8: // Backspace (DEL on Unix, BS elsewhere)
			if len(buffer) > 0 {
				buffer = buffer[:len(buffer)-1]
				io.WriteString(w, "\b \b")
			}
		default: // Echo printable input
			buffer = append(buffer, b)
			w.Write([]byte{b})
		}
	}
}
