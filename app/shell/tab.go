package shell

import (
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/builtins"
	"github.com/codecrafters-io/shell-starter-go/app/completion"
	"github.com/codecrafters-io/shell-starter-go/app/external"
	"github.com/codecrafters-io/shell-starter-go/app/files"
	"github.com/codecrafters-io/shell-starter-go/app/terminal"
)

// HandleTab is the TabHandler entry point called by terminal on each Tab press.
// It computes a completion for the current buffer and applies bash-style tab
// behavior (extend buffer, ring bell, or show listings on a second tab).
func (s *Shell) HandleTab(state *terminal.TabState, buffer string) terminal.TabResult {
	newBuffer, listings := s.completeBuffer(buffer)
	return applyTabAction(state, buffer, newBuffer, listings)
}

// applyTabAction turns a completion result into terminal instructions.
// When listings are returned, the first tab rings the bell and stores them;
// a second tab with the same listings prints them. When the buffer changes,
// the new buffer replaces the input line.
func applyTabAction(state *terminal.TabState, buffer string, newBuffer string, listings []string) terminal.TabResult {
	if len(listings) > 0 {
		if slices.Equal(state.PendingListings, listings) {
			state.PendingListings = nil
			return terminal.TabResult{Buffer: buffer, ListingsToShow: listings}
		}
		state.PendingListings = listings
		return terminal.TabResult{Buffer: buffer, RingBell: true}
	}

	state.PendingListings = nil
	if newBuffer != buffer {
		return terminal.TabResult{Buffer: newBuffer}
	}
	return terminal.TabResult{Buffer: buffer, RingBell: true}
}

// completeBuffer routes the input buffer to the right completion strategy:
// command names before the first space, programmable scripts when registered,
// otherwise filename completion for the current argument.
func (s *Shell) completeBuffer(buffer string) (newBuffer string, listings []string) {
	if !strings.Contains(buffer, " ") {
		return s.completeCommand(buffer)
	}

	if candidates, ok := s.programmableCandidates(buffer); ok {
		return s.completeArgument(buffer, candidates)
	}

	return s.completeFile(buffer)
}

// completeCommand completes the first token against builtin names and PATH
// executables. A unique match appends a trailing space; ambiguous prefixes
// return listings for a second tab press.
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

// commandCandidates returns deduplicated command names from builtins and PATH.
// PATH entries that duplicate a builtin or differ only by Windows extension
// (e.g. echo.exe) are normalized and skipped.
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

// completionCommandName strips Windows executable extensions so echo.exe
// deduplicates against the builtin echo.
func completionCommandName(name string) string {
	switch strings.ToLower(filepath.Ext(name)) {
	case ".exe", ".com", ".bat", ".cmd":
		return strings.TrimSuffix(name, filepath.Ext(name))
	}
	return name
}

// completeFile completes the last argument token against files in the cwd.
func (s *Shell) completeFile(buffer string) (string, []string) {
	return s.completeArgument(buffer, s.fileCandidates(buffer))
}

// fileCandidates lists files in the directory implied by the last argument
// token (e.g. "cat dir/" lists files inside dir/).
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

// completeArgument completes the last argument token against the given
// candidates. Directory prefixes in the token (e.g. "dir/fi") are preserved
// while only the basename is matched.
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

// splitDirToken separates a path-style token into its directory prefix and
// the basename being completed (e.g. "src/ma" → "src/", "ma").
func splitDirToken(token string) (dirPath, prefix string) {
	dirPath = ""
	prefix = token

	if idx := strings.LastIndex(token, "/"); idx >= 0 {
		dirPath = token[:idx+1]
		prefix = token[idx+1:]
	}
	return dirPath, prefix
}

// completionSuffix returns a trailing space after a unique completion, or
// nothing for directories (which already end with "/").
func completionSuffix(token string) string {
	if strings.HasSuffix(token, "/") {
		return ""
	}
	return " "
}

// programmableCandidates runs a registered completer script for the command in
// buffer when one exists. The second return value is false when no script is
// registered, signalling that filename completion should be used instead.
func (s *Shell) programmableCandidates(buffer string) ([]string, bool) {
	opts := buildCompleterOptions(buffer)
	scriptPath, ok := s.completionRegistry.Lookup(opts.Command)
	if !ok {
		return nil, false
	}
	opts.Path = scriptPath

	candidates, err := completion.RunCompleter(opts)
	if err != nil {
		return []string{}, true
	}
	return candidates, true
}

// buildCompleterOptions parses the input buffer into the arguments and
// environment variables expected by a programmable completer script.
func buildCompleterOptions(buffer string) completion.CompleterOptions {
	commandEnd := strings.Index(buffer, " ")
	if commandEnd < 0 {
		return completion.CompleterOptions{
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
		return completion.CompleterOptions{
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

	return completion.CompleterOptions{
		Command:      command,
		CurrentWord:  currentWord,
		PreviousWord: previousWord,
		CompLine:     buffer,
		CompPoint:    len(buffer),
	}
}
