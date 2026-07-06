package completer

import (
	"path/filepath"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/builtins"
	"github.com/codecrafters-io/shell-starter-go/app/completion"
	"github.com/codecrafters-io/shell-starter-go/app/external"
)

func (c *Completer) completeCommand(buffer string) (string, []string) {
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

func commandCandidates() []string {
	candidates := builtins.Names()
	seen := make(map[string]struct{}, len(candidates))
	for _, name := range candidates {
		seen[name] = struct{}{}
	}
	for _, name := range external.FindAllExecutablesInPath() {
		name = commandName(name)
		if _, ok := seen[name]; ok {
			continue
		}
		seen[name] = struct{}{}
		candidates = append(candidates, name)
	}
	return candidates
}

func commandName(name string) string {
	switch strings.ToLower(filepath.Ext(name)) {
	case ".exe", ".com", ".bat", ".cmd":
		return strings.TrimSuffix(name, filepath.Ext(name))
	}
	return name
}

// CompleteCommand completes the first token for tests and direct use.
func (c *Completer) CompleteCommand(buffer string) (string, []string) {
	return c.completeCommand(buffer)
}
