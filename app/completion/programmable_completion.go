package completion

import (
	"os/exec"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/builtins"
)

func RunCompleterScript(scriptPath, command, currentWord, previousWord string) ([]string, error) {
	output, err := exec.Command(scriptPath, command, currentWord, previousWord).Output()
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

func parseCompletionContext(buffer string) (command, currentWord, previousWord string) {
	commandEnd := strings.Index(buffer, " ")
	if commandEnd < 0 {
		return buffer, "", ""
	}

	command = buffer[:commandEnd]
	afterCommand := buffer[commandEnd+1:]

	lastSpace := strings.LastIndex(afterCommand, " ")
	if lastSpace < 0 {
		return command, afterCommand, ""
	}

	currentWord = afterCommand[lastSpace+1:]
	beforeCurrent := afterCommand[:lastSpace]
	prevLastSpace := strings.LastIndex(beforeCurrent, " ")
	if prevLastSpace < 0 {
		return command, currentWord, beforeCurrent
	}
	return command, currentWord, beforeCurrent[prevLastSpace+1:]
}

func applyProgrammableTab(buffer string, completer builtins.Completer) (newBuffer string, listings []string) {
	if completer.Func == nil {
		return buffer, nil
	}

	command, currentWord, previousWord := parseCompletionContext(buffer)
	candidates, err := completer.Func(completer.Path, command, currentWord, previousWord)
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
		return buffer, candidates
	}
}
