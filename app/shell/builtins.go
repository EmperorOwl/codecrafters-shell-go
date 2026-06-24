package shell

import (
	"fmt"
	"io"
	"os"
	"strings"

	shellpath "github.com/codecrafters-io/shell-starter-go/app/path"
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

func EchoOutput(args []string) string {
	return strings.Join(args, " ")
}

func PwdOutput() (string, error) {
	return os.Getwd()
}

func CdErrorMessage(directory string) string {
	return "cd: " + directory + ": No such file or directory"
}

func ChangeDirectory(directory string) error {
	return os.Chdir(directory)
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
	case "pwd":
		cwd, err := PwdOutput()
		if err != nil {
			return true, false
		}
		fmt.Fprintln(out, cwd)
		return true, false
	case "cd":
		if len(fields) < 2 {
			return true, false
		}
		directory := fields[1]
		if err := ChangeDirectory(directory); err != nil {
			fmt.Fprintln(out, CdErrorMessage(directory))
		}
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
