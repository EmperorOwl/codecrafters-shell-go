package completion

import (
	"slices"
	"strings"
)

func matches(builtins []string, prefix string) []string {
	var result []string
	for _, name := range builtins {
		if strings.HasPrefix(name, prefix) {
			result = append(result, name)
		}
	}
	slices.Sort(result)
	return result
}

// ApplyTab returns an updated command buffer after Tab.
// listings is non-empty when multiple builtins share the prefix.
func ApplyTab(builtins []string, buffer string) (newBuffer string, listings []string) {
	matched := matches(builtins, buffer)
	switch len(matched) {
	case 0:
		return buffer, nil
	case 1:
		return matched[0] + " ", nil
	default:
		return buffer, matched
	}
}
