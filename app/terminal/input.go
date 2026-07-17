package terminal

import (
	"bufio"
	"io"
	"strings"
)

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
	// skipNextLF is set after CR so the LF from a Windows CRLF Enter is discarded
	// instead of submitting an empty line.
	var skipNextLF bool

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
			// ESC starts an ANSI escape sequence. Arrow keys send ESC [ A (up) or
			// ESC [ B (down); readLineRaw already consumed ESC, so delegate the rest.
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

// handleEscapeSequence reads the bytes after ESC and handles arrow-key history
// navigation in raw mode.
//
// Terminals send arrow keys as 3-byte CSI sequences:
//   - up:   ESC [ A  (\x1b[A)
//   - down: ESC [ B  (\x1b[B)
//
// Because readLineRaw already consumed ESC, this function reads the remaining
// two bytes ([ and A/B). Unrecognized sequences return ok=false so the caller
// can ignore them.
//
// Return values tell readLineRaw how to update the input buffer:
//   - ok=false: not a handled arrow sequence; ignore ESC
//   - ok=true, handled=true, newBuffer: replace buffer and redraw the prompt
//   - ok=true, handled=true, unchanged buffer: arrow pressed at a history
//     boundary; bell was rung and the current line is kept
func handleEscapeSequence(reader *bufio.Reader, w io.Writer, historyHandler HistoryHandler, historyState *historyBrowseState, buffer []byte) (handled bool, newBuffer []byte, ok bool) {
	// Expect CSI introducer '[' after ESC.
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

	if code != 'A' && code != 'B' {
		return false, buffer, false
	}

	var command string
	var found bool
	switch code {
	case 'A':
		command, found = historyState.stepUp(historyHandler)
	case 'B':
		command, found = historyState.stepDown(historyHandler)
	}

	if !found {
		ringBell(w)
		return true, buffer, true
	}

	newBuffer = []byte(command)
	redrawLine(w, command)
	return true, newBuffer, true
}
