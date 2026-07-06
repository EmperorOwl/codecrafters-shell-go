package builtins

import (
	"fmt"
	"io"

	"github.com/codecrafters-io/shell-starter-go/app/external"
)

func TypeOutput(command string) string {
	if IsBuiltin(command) {
		return command + " is a shell builtin"
	}
	if path, ok := external.FindExecutableInPath(command); ok {
		return command + " is " + path
	}
	return command + ": not found"
}

func Type(out io.Writer, command string) {
	fmt.Fprintln(out, TypeOutput(command))
}

func init() {
	register("type", typeBuiltin)
}

func typeBuiltin(ctx *Context, args []string) (bool, error) {
	target := ""
	if len(args) > 0 {
		target = args[0]
	}
	Type(ctx.Stdout, target)
	return false, nil
}
