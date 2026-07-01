package completion

import (
	"strings"
)

// ApplyTab returns an updated command buffer after Tab.
// listings is non-empty when multiple commands share the prefix.
func ApplyTab(
	builtinsList, executables []string,
	listFiles FileLister,
	completeHandler CompleteHandler,
	buffer string,
) (newBuffer string, listings []string) {
	if strings.Contains(buffer, " ") {
		if completeHandler != nil {
			if newBuffer, listings, ok := applyProgrammableTab(buffer, completeHandler); ok {
				return newBuffer, listings
			}
		}
		return applyFileTab(listFiles, buffer)
	}
	return applyCommandTab(builtinsList, executables, buffer)
}
