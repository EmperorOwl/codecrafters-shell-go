package completion

import (
	"strings"
)

// ApplyTab returns an updated command buffer after Tab.
// listings is non-empty when multiple commands share the prefix.
func ApplyTab(
	builtinsList, executables []string,
	listFiles FileLister,
	completerFuncs map[string]CompleterFunc,
	buffer string,
) (newBuffer string, listings []string) {
	if strings.Contains(buffer, " ") {
		command := buffer[:strings.Index(buffer, " ")]
		if completer, ok := completerFuncs[command]; ok {
			return applyProgrammableTab(buffer, completer)
		}
		return applyFileTab(listFiles, buffer)
	}
	return applyCommandTab(builtinsList, executables, buffer)
}
