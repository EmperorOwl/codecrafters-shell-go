package completion

import (
	"slices"
	"strings"

	shellpath "github.com/codecrafters-io/shell-starter-go/app/path"
)

func builtinMatches(builtins []string, prefix string) []string {
	var result []string
	for _, name := range builtins {
		if strings.HasPrefix(name, prefix) {
			result = append(result, name)
		}
	}
	slices.Sort(result)
	return result
}

func allMatches(builtins []string, prefix string) []string {
	matched := append(builtinMatches(builtins, prefix), shellpath.FindMatchingExecutablesInPath(prefix)...)
	slices.Sort(matched)
	return slices.Compact(matched)
}

// ApplyTab returns an updated command buffer after Tab.
// listings is non-empty when multiple commands share the prefix.
func ApplyTab(builtins []string, buffer string) (newBuffer string, listings []string) {
	matched := allMatches(builtins, buffer)
	switch len(matched) {
	case 0:
		return buffer, nil
	case 1:
		return matched[0] + " ", nil
	default:
		return buffer, matched
	}
}
