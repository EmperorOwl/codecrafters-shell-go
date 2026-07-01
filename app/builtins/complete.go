package builtins

import (
	"fmt"
	"io"
)

func NoCompletionSpecMessage(command string) string {
	return "complete: " + command + ": no completion specification"
}

func RegisteredSpecMessage(scriptPath, command string) string {
	return "complete -C '" + scriptPath + "' " + command
}

// Complete handles the complete builtin. completers maps command names
// to completer script paths and is owned by the caller.
func Complete(stdout, stderr io.Writer, args []string, completers map[string]string) {
	if len(args) == 0 {
		return
	}

	switch args[0] {
	// Print the registered completion spec for a command, or an error if none exists.
	case "-p":
		if len(args) < 2 {
			return
		}
		command := args[1]
		if scriptPath, ok := completers[command]; ok {
			fmt.Fprintln(stdout, RegisteredSpecMessage(scriptPath, command))
			return
		}
		fmt.Fprintln(stderr, NoCompletionSpecMessage(command))
	// Register a completer script for a command.
	case "-C":
		if len(args) < 3 {
			return
		}
		completers[args[2]] = args[1]
	// Remove the completion rule for a command.
	case "-r":
		if len(args) < 2 {
			return
		}
		delete(completers, args[1])
	}
}
