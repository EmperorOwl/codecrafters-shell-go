package terminal

import (
	"bufio"
	"io"
	"strings"
)

// skipNextLF is set after CR so the LF from a Windows CRLF Enter is discarded
// on the next readLine call instead of submitting an empty line.
var skipNextLF bool

func readLine(reader *bufio.Reader, w io.Writer, rawMode bool, tabHandler TabHandler, historyHandler HistoryHandler) (line string, eof bool, err error) {
	if rawMode {
		return readLineRaw(reader, w, tabHandler, historyHandler)
	}

	text, err := reader.ReadString('\n')
	if err == io.EOF {
		return strings.TrimSpace(text), true, nil
	}
	if err != nil {
		return "", false, err
	}
	return strings.TrimSpace(text), false, nil
}

func readLineRaw(reader *bufio.Reader, w io.Writer, tabHandler TabHandler, historyHandler HistoryHandler) (string, bool, error) {
	var buffer []byte
	var tabState TabState
	var historyState historyBrowseState

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
		case '\x1b':
			if handled, newBuffer, ok := handleEscapeSequence(reader, w, historyHandler, &historyState, buffer); ok {
				if handled {
					buffer = newBuffer
				}
				continue
			}
		case '\t':
			if tabHandler == nil {
				ringBell(w)
				continue
			}
			result := tabHandler.HandleTab(&tabState, string(buffer))
			if result.RingBell {
				ringBell(w)
			}
			if len(result.ListingsToShow) > 0 {
				writeListings(w, result.ListingsToShow)
				redrawLine(w, result.Buffer)
			} else if result.Buffer != string(buffer) {
				buffer = []byte(result.Buffer)
				redrawLine(w, result.Buffer)
			}
		case '\r':
			skipNextLF = true
			writeCRLF(w)
			return string(buffer), false, nil
		case '\n':
			if skipNextLF {
				skipNextLF = false
				continue
			}
			writeCRLF(w)
			return string(buffer), false, nil
		case 127, 8:
			if len(buffer) > 0 {
				historyState.reset()
				buffer = buffer[:len(buffer)-1]
				writeBackspace(w)
			}
		default:
			historyState.reset()
			buffer = append(buffer, b)
			w.Write([]byte{b})
		}
	}
}

func handleEscapeSequence(reader *bufio.Reader, w io.Writer, historyHandler HistoryHandler, historyState *historyBrowseState, buffer []byte) (handled bool, newBuffer []byte, ok bool) {
	next, err := reader.ReadByte()
	if err != nil {
		return false, buffer, false
	}
	if next != '[' {
		return false, buffer, false
	}

	code, err := reader.ReadByte()
	if err != nil {
		return false, buffer, false
	}

	if code != 'A' {
		return false, buffer, false
	}

	command, found := historyState.stepUp(historyHandler)
	if !found {
		ringBell(w)
		return true, buffer, true
	}

	newBuffer = []byte(command)
	redrawLine(w, command)
	return true, newBuffer, true
}
