package builtins

import (
	"fmt"
	"io"

	"github.com/codecrafters-io/shell-starter-go/app/completion"
)

func init() {
	register("complete", completeBuiltin)
}

func completeBuiltin(ctx *Context, args []string) (bool, error) {
	if ctx.State == nil {
		return false, nil
	}
	Complete(ctx.Stdout, ctx.Stderr, args, ctx.State.Completion)
	return false, nil
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
			fmt.Fprintln(stdout, registeredSpecMessage(scriptPath, command))
			return
		}
		fmt.Fprintln(stderr, noCompletionSpecMessage(command))
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

func noCompletionSpecMessage(command string) string {
	return "complete: " + command + ": no completion specification"
}

func registeredSpecMessage(scriptPath, command string) string {
	return "complete -C '" + scriptPath + "' " + command
}
