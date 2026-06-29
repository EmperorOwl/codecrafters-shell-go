package builtins

import (
	"fmt"
	"io"
)

func NoCompletionSpecMessage(command string) string {
	return "complete: " + command + ": no completion specification"
}

func Complete(stderr io.Writer, args []string) {
	if len(args) >= 2 && args[0] == "-p" {
		fmt.Fprintln(stderr, NoCompletionSpecMessage(args[1]))
	}
}
