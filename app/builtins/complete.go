package builtins

import (
	"fmt"
	"io"

	"github.com/codecrafters-io/shell-starter-go/app/completion"
)

func NoCompletionSpecMessage(command string) string {
	return "complete: " + command + ": no completion specification"
}

func RegisteredSpecMessage(scriptPath, command string) string {
	return "complete -C '" + scriptPath + "' " + command
}

// Complete handles the complete builtin.
func Complete(stdout, stderr io.Writer, args []string, registry *completion.CompletionRegistry) {
	if len(args) == 0 || registry == nil {
		return
	}

	switch args[0] {
	case "-p":
		if len(args) < 2 {
			return
		}
		command := args[1]
		if scriptPath, ok := registry.Lookup(command); ok {
			fmt.Fprintln(stdout, RegisteredSpecMessage(scriptPath, command))
			return
		}
		fmt.Fprintln(stderr, NoCompletionSpecMessage(command))
	case "-C":
		if len(args) < 3 {
			return
		}
		registry.Register(args[2], args[1])
	case "-r":
		if len(args) < 2 {
			return
		}
		registry.Unregister(args[1])
	}
}

func completeBuiltin(ctx *Context, args []string) (bool, error) {
	Complete(ctx.Stdout, ctx.Stderr, args, ctx.Completion)
	return false, nil
}
