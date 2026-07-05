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

// CompleteHandler returns programmable completion candidates for the given context.
// A nil return means no completer is registered for the command.
type CompleteHandler func(opts CompleterFuncOptions) []string

// CompleteCommand runs the registered completer script for the given command,
// returning the completion candidates.
func CompleteCommand(registry *CompletionRegistry, opts CompleterFuncOptions) []string {
	if registry == nil {
		return nil
	}
	scriptPath, ok := registry.Lookup(opts.Command)
	if !ok {
		return nil
	}
	opts.ScriptPath = scriptPath

	candidates, err := runCompleterScript(opts)
	if err != nil {
		return []string{}
	}
	return candidates
}

func runCompleterScript(opts CompleterFuncOptions) ([]string, error) {
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

func applyProgrammableTab(buffer string, completeHandler CompleteHandler) (newBuffer string, listings []string, handled bool) {
	opts := buildCompleterFuncOptions(buffer)
	candidates := completeHandler(opts)
	if candidates == nil {
		return buffer, nil, false
	}
	if len(candidates) == 0 {
		return buffer, nil, true
	}

	lastSpace := strings.LastIndex(buffer, " ")
	if lastSpace < 0 {
		return buffer, nil, true
	}

	prefix := buffer[:lastSpace+1]
	matches := findMatches(candidates, opts.CurrentWord)
	switch len(matches) {
	case 0:
		return buffer, nil, true
	case 1:
		return prefix + matches[0] + " ", nil, true
	default:
		lcp := longestCommonPrefix(matches)
		if len(lcp) > len(opts.CurrentWord) {
			return prefix + lcp, nil, true
		}
		return buffer, matches, true
	}
}
