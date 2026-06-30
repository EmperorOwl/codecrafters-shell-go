package completion

import (
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// CompleterFuncOptions holds the context passed to a completer script.
type CompleterFuncOptions struct {
	ScriptPath   string
	Command      string
	CurrentWord  string
	PreviousWord string
	CompLine     string
	CompPoint    int
}

// CompleterFunc runs a completer and returns completion candidates.
type CompleterFunc func(opts CompleterFuncOptions) ([]string, error)

// BuildCompleterFuncs maps registered script paths to completer functions.
func BuildCompleterFuncs(registeredCompleters map[string]string) map[string]CompleterFunc {
	funcs := make(map[string]CompleterFunc, len(registeredCompleters))
	for command, scriptPath := range registeredCompleters {
		funcs[command] = completerFuncFor(scriptPath)
	}
	return funcs
}

func completerFuncFor(scriptPath string) CompleterFunc {
	return func(opts CompleterFuncOptions) ([]string, error) {
		opts.ScriptPath = scriptPath
		return RunCompleterScript(opts)
	}
}

func RunCompleterScript(opts CompleterFuncOptions) ([]string, error) {
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

func buildCompleterFuncOptions(buffer string) CompleterFuncOptions {
	commandEnd := strings.Index(buffer, " ")
	if commandEnd < 0 {
		return CompleterFuncOptions{
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
		return CompleterFuncOptions{
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

	return CompleterFuncOptions{
		Command:      command,
		CurrentWord:  currentWord,
		PreviousWord: previousWord,
		CompLine:     buffer,
		CompPoint:    len(buffer),
	}
}

func applyProgrammableTab(buffer string, completer CompleterFunc) (newBuffer string, listings []string) {
	if completer == nil {
		return buffer, nil
	}

	opts := buildCompleterFuncOptions(buffer)
	candidates, err := completer(opts)
	if err != nil || len(candidates) == 0 {
		return buffer, nil
	}

	lastSpace := strings.LastIndex(buffer, " ")
	if lastSpace < 0 {
		return buffer, nil
	}

	prefix := buffer[:lastSpace+1]
	matches := findMatches(candidates, opts.CurrentWord)
	switch len(matches) {
	case 0:
		return buffer, nil
	case 1:
		return prefix + matches[0] + " ", nil
	default:
		lcp := longestCommonPrefix(matches)
		if len(lcp) > len(opts.CurrentWord) {
			return prefix + lcp, nil
		}
		return buffer, matches
	}
}
