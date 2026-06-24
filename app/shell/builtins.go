package shell

import (
	"fmt"
	"io"
	"os"
	"strings"

	shellpath "github.com/codecrafters-io/shell-starter-go/app/path"
)

var shellBuiltins = map[string]struct{}{
	"echo": {},
	"exit": {},
	"type": {},
}

func IsShellBuiltin(command string) bool {
	_, ok := shellBuiltins[command]
	return ok
}

func EchoOutput(args []string) string {
	return strings.Join(args, " ")
}

func TypeOutput(command string) string {
	if IsShellBuiltin(command) {
		return command + " is a shell builtin"
	}
	if path, ok := shellpath.FindExecutableInPath(command, os.Getenv("PATH")); ok {
		return command + " is " + path
	}
	return command + ": not found"
}

func TryBuiltin(line string, out io.Writer) (handled bool, shouldExit bool) {
	fields := strings.Fields(line)
	if len(fields) == 0 {
		return false, false
	}

	switch fields[0] {
	case "exit":
		return true, true
	case "echo":
		fmt.Fprintln(out, EchoOutput(fields[1:]))
		return true, false
	case "type":
		target := ""
		if len(fields) > 1 {
			target = fields[1]
		}
		fmt.Fprintln(out, TypeOutput(target))
		return true, false
	default:
		return false, false
	}
}
