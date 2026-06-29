package builtins

import (
	"fmt"
	"io"
)

var completionSpecs = map[string]string{}

func NoCompletionSpecMessage(command string) string {
	return "complete: " + command + ": no completion specification"
}

func RegisteredSpecMessage(scriptPath, command string) string {
	return "complete -C '" + scriptPath + "' " + command
}

func Complete(stdout, stderr io.Writer, args []string) {
	if len(args) == 0 {
		return
	}

	switch args[0] {
	case "-p":
		if len(args) < 2 {
			return
		}
		command := args[1]
		if scriptPath, ok := completionSpecs[command]; ok {
			fmt.Fprintln(stdout, RegisteredSpecMessage(scriptPath, command))
			return
		}
		fmt.Fprintln(stderr, NoCompletionSpecMessage(command))
	case "-C":
		if len(args) < 3 {
			return
		}
		completionSpecs[args[2]] = args[1]
	}
}
