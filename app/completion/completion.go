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

func longestCommonPrefix(values []string) string {
	if len(values) == 0 {
		return ""
	}

	prefix := values[0]
	for _, value := range values[1:] {
		for len(prefix) > 0 && !strings.HasPrefix(value, prefix) {
			prefix = prefix[:len(prefix)-1]
		}
		if prefix == "" {
			return ""
		}
	}

	return prefix
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
		lcp := longestCommonPrefix(matched)
		if len(lcp) > len(buffer) {
			return lcp, nil
		}
		return buffer, matched
	}
}
