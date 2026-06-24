package builtins

import (
	"fmt"
	"io"
	"os"

	shellpath "github.com/codecrafters-io/shell-starter-go/app/path"
)

func TypeOutput(command string, isBuiltin bool) string {
	if isBuiltin {
		return command + " is a shell builtin"
	}
	if path, ok := shellpath.FindExecutableInPath(command, os.Getenv("PATH")); ok {
		return command + " is " + path
	}
	return command + ": not found"
}

func Type(out io.Writer, command string, isBuiltin bool) {
	fmt.Fprintln(out, TypeOutput(command, isBuiltin))
}
