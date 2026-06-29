package builtins

import (
	"fmt"
	"io"

	shellpath "github.com/codecrafters-io/shell-starter-go/app/path"
)

func TypeOutput(command string, isBuiltin bool) string {
	if isBuiltin {
		return command + " is a shell builtin"
	}
	if path, ok := shellpath.FindExecutableInPath(command); ok {
		return command + " is " + path
	}
	return command + ": not found"
}

func Type(out io.Writer, command string, isBuiltin bool) {
	fmt.Fprintln(out, TypeOutput(command, isBuiltin))
}
