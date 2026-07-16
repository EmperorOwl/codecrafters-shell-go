package builtins

import (
	"fmt"
	"io"

	"github.com/codecrafters-io/shell-starter-go/app/completion"
)

func init() {
	register("complete", completeHandler)
}

func completeHandler(ctx *Context, args []string) (bool, error) {
	if ctx.State == nil {
		return false, nil
	}
	completeBuiltin(ctx.Stdout, ctx.Stderr, args, ctx.State.Completion)
	return false, nil
}

func completeBuiltin(stdout, stderr io.Writer, args []string, registry *completion.Registry) {
	if len(args) == 0 || registry == nil {
		return
	}

	switch args[0] {
	case "-p": // print completion spec
		if len(args) < 2 {
			return
		}
		command := args[1]
		if scriptPath, ok := registry.Lookup(command); ok {
			fmt.Fprintln(stdout, completeRegisteredSpecMessage(scriptPath, command))
			return
		}
		fmt.Fprintln(stderr, completeNoSpecMessage(command))
	case "-C": // register completion script
		if len(args) < 3 {
			return
		}
		registry.Register(args[2], args[1])
	case "-r": // remove completion spec
		if len(args) < 2 {
			return
		}
		registry.Unregister(args[1])
	}
}

func completeNoSpecMessage(command string) string {
	return "complete: " + command + ": no completion specification"
}

func completeRegisteredSpecMessage(scriptPath, command string) string {
	return "complete -C '" + scriptPath + "' " + command
}
