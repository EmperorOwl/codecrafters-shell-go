package completion

import "strings"

// FileLister returns sorted file and directory names in dir relative to a base directory.
// An empty dir lists the base directory. Directory names include a trailing slash.
type FileLister func(dir string) []string

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
