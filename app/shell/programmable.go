package shell

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/builtins"
	"github.com/codecrafters-io/shell-starter-go/app/completion"
	"github.com/codecrafters-io/shell-starter-go/app/external"
	"github.com/codecrafters-io/shell-starter-go/app/files"
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

func (s *Shell) programmableCandidates(buffer string) ([]string, bool) {
	opts := buildCompleterFuncOptions(buffer)
	scriptPath, ok := s.completionRegistry.Lookup(opts.Command)
	if !ok {
		return nil, false
	}
	opts.ScriptPath = scriptPath

	candidates, err := runCompleterScript(opts)
	if err != nil {
		return []string{}, true
	}
	return candidates, true
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

func (s *Shell) completeBuffer(buffer string) (newBuffer string, listings []string) {
	if !strings.Contains(buffer, " ") {
		return s.completeCommand(buffer)
	}

	if candidates, ok := s.programmableCandidates(buffer); ok {
		return s.completeArgument(buffer, candidates)
	}

	return s.completeFile(buffer)
}

func (s *Shell) completeCommand(buffer string) (string, []string) {
	token, listings, unique := completion.Complete(buffer, commandCandidates())
	if len(listings) > 0 {
		return buffer, listings
	}
	if token != buffer {
		suffix := ""
		if unique {
			suffix = completionSuffix(token)
		}
		return token + suffix, nil
	}
	return buffer, nil
}

func (s *Shell) completeFile(buffer string) (string, []string) {
	return s.completeArgument(buffer, s.fileCandidates(buffer))
}

func (s *Shell) fileCandidates(buffer string) []string {
	lastSpace := strings.LastIndex(buffer, " ")
	if lastSpace < 0 {
		return nil
	}

	token := buffer[lastSpace+1:]
	dirPath, _ := splitDirToken(token)

	cwd, err := os.Getwd()
	if err != nil {
		return nil
	}
	return files.ListInDir(cwd, dirPath)
}

func (s *Shell) completeArgument(buffer string, candidates []string) (string, []string) {
	lastSpace := strings.LastIndex(buffer, " ")
	if lastSpace < 0 {
		return buffer, nil
	}

	bufferPrefix := buffer[:lastSpace+1]
	token := buffer[lastSpace+1:]
	dirPath, prefix := splitDirToken(token)

	token, listings, unique := completion.Complete(prefix, candidates)
	if len(listings) > 0 {
		return buffer, listings
	}
	if token != prefix {
		suffix := ""
		if unique {
			suffix = completionSuffix(token)
		}
		return bufferPrefix + dirPath + token + suffix, nil
	}
	return buffer, nil
}

func commandCandidates() []string {
	candidates := builtins.Names()
	seen := make(map[string]struct{}, len(candidates))
	for _, name := range candidates {
		seen[name] = struct{}{}
	}
	for _, name := range external.FindAllExecutablesInPath() {
		name = completionCommandName(name)
		if _, ok := seen[name]; ok {
			continue
		}
		seen[name] = struct{}{}
		candidates = append(candidates, name)
	}
	return candidates
}

func completionCommandName(name string) string {
	switch strings.ToLower(filepath.Ext(name)) {
	case ".exe", ".com", ".bat", ".cmd":
		return strings.TrimSuffix(name, filepath.Ext(name))
	}
	return name
}

func splitDirToken(token string) (dirPath, prefix string) {
	dirPath = ""
	prefix = token

	if idx := strings.LastIndex(token, "/"); idx >= 0 {
		dirPath = token[:idx+1]
		prefix = token[idx+1:]
	}
	return dirPath, prefix
}

func completionSuffix(token string) string {
	if strings.HasSuffix(token, "/") {
		return ""
	}
	return " "
}
