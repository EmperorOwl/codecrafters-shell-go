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
