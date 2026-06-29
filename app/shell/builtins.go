package shell

import (
	"io"
	"slices"

	"github.com/codecrafters-io/shell-starter-go/app/builtins"
)

var shellBuiltins = map[string]struct{}{
	"cd":       {},
	"complete": {},
	"echo":     {},
	"exit":     {},
	"pwd":      {},
	"type":     {},
}

func IsShellBuiltin(command string) bool {
	_, ok := shellBuiltins[command]
	return ok
}

func BuiltinNames() []string {
	names := make([]string, 0, len(shellBuiltins))
	for name := range shellBuiltins {
		names = append(names, name)
	}
	slices.Sort(names)
	return names
}

func TryBuiltin(fields []string, stdout, stderr io.Writer) (handled bool, shouldExit bool) {
	if len(fields) == 0 {
		return false, false
	}

	switch fields[0] {
	case "exit":
		return true, builtins.Exit()
	case "echo":
		builtins.Echo(stdout, fields[1:])
		return true, false
	case "pwd":
		builtins.Pwd(stdout)
		return true, false
	case "cd":
		directory := ""
		if len(fields) > 1 {
			directory = fields[1]
		}
		builtins.Cd(stderr, directory)
		return true, false
	case "type":
		target := ""
		if len(fields) > 1 {
			target = fields[1]
		}
		builtins.Type(stdout, target, IsShellBuiltin(target))
		return true, false
	case "complete":
		builtins.Complete(stderr, fields[1:])
		return true, false
	default:
		return false, false
	}
}
