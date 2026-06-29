package completion

import (
	"os/exec"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/builtins"
)

func RunCompleterScript(scriptPath string) ([]string, error) {
	output, err := exec.Command(scriptPath).Output()
	if err != nil {
		return nil, err
	}

	text := strings.TrimRight(string(output), "\r\n")
	if text == "" {
		return nil, nil
	}

	lines := strings.Split(text, "\n")
	candidates := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSuffix(line, "\r")
		if line != "" {
			candidates = append(candidates, line)
		}
	}
	return candidates, nil
}

func applyProgrammableTab(buffer string, completer builtins.Completer) (newBuffer string, listings []string) {
	if completer.Func == nil {
		return buffer, nil
	}

	candidates, err := completer.Func(completer.Path)
	if err != nil || len(candidates) == 0 {
		return buffer, nil
	}

	commandEnd := strings.Index(buffer, " ")
	command := buffer[:commandEnd]

	switch len(candidates) {
	case 1:
		return command + " " + candidates[0] + " ", nil
	default:
		return buffer, candidates
	}
}
