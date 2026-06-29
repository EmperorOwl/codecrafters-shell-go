package completion

import (
	"slices"
	"strings"
)

func findMatches(commands []string, prefix string) []string {
	var result []string
	for _, name := range commands {
		if strings.HasPrefix(name, prefix) {
			result = append(result, name)
		}
	}
	slices.Sort(result)
	return result
}

// ApplyTab returns an updated command buffer after Tab.
// listings is non-empty when multiple commands share the prefix.
func ApplyTab(builtins, executables []string, buffer string) (newBuffer string, listings []string) {
	matched := append(findMatches(builtins, buffer), findMatches(executables, buffer)...)
	slices.Sort(matched)
	matched = slices.Compact(matched)

	switch len(matched) {
	case 0:
		return buffer, nil
	case 1:
		return matched[0] + " ", nil
	default:
		return buffer, matched
	}
}
