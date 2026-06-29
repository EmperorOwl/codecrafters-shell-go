package completion

import "strings"

// ApplyTab returns an updated command buffer after Tab.
// listings is non-empty when multiple commands share the prefix.
func ApplyTab(builtins, executables []string, listFiles FileLister, buffer string) (newBuffer string, listings []string) {
	if strings.Contains(buffer, " ") {
		return applyFileTab(listFiles, buffer)
	}
	return applyCommandTab(builtins, executables, buffer)
}
