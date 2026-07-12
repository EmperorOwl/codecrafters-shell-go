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
		fmt.Fprintln(stdout, shellBuiltinMessage(command))
		return
	}
	if path, ok := external.FindExecutableInPath(command); ok {
		fmt.Fprintln(stdout, executableMessage(command, path))
		return
	}
	fmt.Fprintln(stdout, commandNotFoundMessage(command))
}

func shellBuiltinMessage(command string) string {
	return command + " is a shell builtin"
}

func executableMessage(command, path string) string {
	return command + " is " + path
}

func commandNotFoundMessage(command string) string {
	return command + ": not found"
}
