package completion

import (
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/builtins"
)

// ApplyTab returns an updated command buffer after Tab.
// listings is non-empty when multiple commands share the prefix.
func ApplyTab(
	builtinsList, executables []string,
	listFiles FileLister,
	registeredCompleters map[string]builtins.Completer,
	buffer string,
) (newBuffer string, listings []string) {
	if strings.Contains(buffer, " ") {
		command := buffer[:strings.Index(buffer, " ")]
		if completer, ok := registeredCompleters[command]; ok {
			return applyProgrammableTab(buffer, completer)
		}
		return applyFileTab(listFiles, buffer)
	}
	return applyCommandTab(builtinsList, executables, buffer)
}
