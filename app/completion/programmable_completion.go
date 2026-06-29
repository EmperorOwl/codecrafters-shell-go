package completion

import (
	"os"
	"os/exec"
	"slices"
	"strconv"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/builtins"
)

func RunCompleterScript(opts builtins.CompleterFuncOptions) ([]string, error) {
	cmd := exec.Command(opts.ScriptPath, opts.Command, opts.CurrentWord, opts.PreviousWord)
	cmd.Env = append(os.Environ(),
		"COMP_LINE="+opts.CompLine,
		"COMP_POINT="+strconv.Itoa(opts.CompPoint),
	)

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return parseCompleterOutput(output), nil
}

func parseCompleterOutput(output []byte) []string {
	text := strings.TrimRight(string(output), "\r\n")
	if text == "" {
		return nil
	}

	lines := strings.Split(text, "\n")
	candidates := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSuffix(line, "\r")
		if line != "" {
			candidates = append(candidates, line)
		}
	}
	return candidates
}

func buildCompleterFuncOptions(buffer string) builtins.CompleterFuncOptions {
	commandEnd := strings.Index(buffer, " ")
	if commandEnd < 0 {
		return builtins.CompleterFuncOptions{
			Command:   buffer,
			CompLine:  buffer,
			CompPoint: len(buffer),
		}
	}

	command := buffer[:commandEnd]
	afterCommand := buffer[commandEnd+1:]

	lastSpace := strings.LastIndex(afterCommand, " ")
	if lastSpace < 0 {
		currentWord := afterCommand
		previousWord := ""
		if currentWord != "" {
			previousWord = command
		}
		return builtins.CompleterFuncOptions{
			Command:      command,
			CurrentWord:  currentWord,
			PreviousWord: previousWord,
			CompLine:     buffer,
			CompPoint:    len(buffer),
		}
	}

	currentWord := afterCommand[lastSpace+1:]
	beforeCurrent := afterCommand[:lastSpace]
	prevLastSpace := strings.LastIndex(beforeCurrent, " ")
	previousWord := beforeCurrent
	if prevLastSpace >= 0 {
		previousWord = beforeCurrent[prevLastSpace+1:]
	}

	return builtins.CompleterFuncOptions{
		Command:      command,
		CurrentWord:  currentWord,
		PreviousWord: previousWord,
		CompLine:     buffer,
		CompPoint:    len(buffer),
	}
}

func applyProgrammableTab(buffer string, completer builtins.Completer) (newBuffer string, listings []string) {
	if completer.Func == nil {
		return buffer, nil
	}

	opts := buildCompleterFuncOptions(buffer)
	opts.ScriptPath = completer.Path
	candidates, err := completer.Func(opts)
	if err != nil || len(candidates) == 0 {
		return buffer, nil
	}

	switch len(candidates) {
	case 1:
		lastSpace := strings.LastIndex(buffer, " ")
		if lastSpace < 0 {
			return buffer, nil
		}
		return buffer[:lastSpace+1] + candidates[0] + " ", nil
	default:
		slices.Sort(candidates)
		return buffer, candidates
	}
}
