package builtins

import (
	"fmt"
	"io"

	"github.com/codecrafters-io/shell-starter-go/app/external"
)

func init() {
	register("type", typeHandler)
}

func typeHandler(ctx *Context, args []string) (bool, error) {
	target := ""
	if len(args) > 0 {
		target = args[0]
	}
	typeBuiltin(ctx.Stdout, target)
	return false, nil
}

func typeBuiltin(stdout io.Writer, command string) {
	if IsBuiltin(command) {
		fmt.Fprintln(stdout, typeBuiltinMessage(command))
		return
	}
	if path, ok := external.FindExecutableInPath(command); ok {
		fmt.Fprintln(stdout, typeExecutableMessage(command, path))
		return
	}
	fmt.Fprintln(stdout, typeNotFoundMessage(command))
}

func typeBuiltinMessage(command string) string {
	return command + " is a shell builtin"
}

func typeExecutableMessage(command, path string) string {
	return command + " is " + path
}

func typeNotFoundMessage(command string) string {
	return command + ": not found"
}
