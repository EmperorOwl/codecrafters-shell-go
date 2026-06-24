package shell

import (
	"io"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/builtins"
)

var shellBuiltins = map[string]struct{}{
	"cd":   {},
	"echo": {},
	"exit": {},
	"pwd":  {},
	"type": {},
}

func IsShellBuiltin(command string) bool {
	_, ok := shellBuiltins[command]
	return ok
}

func TryBuiltin(line string, out io.Writer) (handled bool, shouldExit bool) {
	fields := strings.Fields(line)
	if len(fields) == 0 {
		return false, false
	}

	switch fields[0] {
	case "exit":
		return true, builtins.Exit()
	case "echo":
		builtins.Echo(out, fields[1:])
		return true, false
	case "pwd":
		builtins.Pwd(out)
		return true, false
	case "cd":
		directory := ""
		if len(fields) > 1 {
			directory = fields[1]
		}
		builtins.Cd(out, directory)
		return true, false
	case "type":
		target := ""
		if len(fields) > 1 {
			target = fields[1]
		}
		builtins.Type(out, target, IsShellBuiltin(target))
		return true, false
	default:
		return false, false
	}
}
