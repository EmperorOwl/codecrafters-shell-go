package completion

import (
	"slices"
	"strings"
)

// FileLister returns sorted file and directory names in dir relative to the current working directory.
// An empty dir lists the current working directory. Directory names include a trailing slash.
type FileLister func(dir string) []string

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

func completeEntry(name string) string {
	if strings.HasSuffix(name, "/") {
		return name
	}
	return name + " "
}

func applyFileTab(listFiles FileLister, buffer string) (newBuffer string, listings []string) {
	lastSpace := strings.LastIndex(buffer, " ")
	if lastSpace < 0 {
		return buffer, nil
	}

	token := buffer[lastSpace+1:]
	dirPath := ""
	prefix := token

	if idx := strings.LastIndex(token, "/"); idx >= 0 {
		dirPath = token[:idx+1]
		prefix = token[idx+1:]
	}

	matched := findMatches(listFiles(dirPath), prefix)

	switch len(matched) {
	case 0:
		return buffer, nil
	case 1:
		return buffer[:lastSpace+1] + dirPath + completeEntry(matched[0]), nil
	default:
		lcp := longestCommonPrefix(matched)
		if len(lcp) > len(prefix) {
			return buffer[:lastSpace+1] + dirPath + lcp, nil
		}
		return buffer, matched
	}
}

// ApplyTab returns an updated command buffer after Tab.
// listings is non-empty when multiple commands share the prefix.
func ApplyTab(builtins, executables []string, listFiles FileLister, buffer string) (newBuffer string, listings []string) {
	if strings.Contains(buffer, " ") {
		return applyFileTab(listFiles, buffer)
	}

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
