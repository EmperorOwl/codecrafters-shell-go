package builtins

import (
	"fmt"
	"io"
)

// CompleterFunc runs a completer script and returns completion candidates.
type CompleterFunc func(scriptPath, command, currentWord, previousWord string) ([]string, error)

// Completer holds a registered completer script and its runner.
type Completer struct {
	Path string
	Func CompleterFunc
}

func NoCompletionSpecMessage(command string) string {
	return "complete: " + command + ": no completion specification"
}

func RegisteredSpecMessage(scriptPath, command string) string {
	return "complete -C '" + scriptPath + "' " + command
}

// Complete handles the complete builtin. registeredCompleters maps command names
// to their completers and is owned by the shell session.
func Complete(stdout, stderr io.Writer, args []string, registeredCompleters map[string]Completer, runCompleter CompleterFunc) {
	if len(args) == 0 {
		return
	}

	switch args[0] {
	case "-p":
		if len(args) < 2 {
			return
		}
		command := args[1]
		if completer, ok := registeredCompleters[command]; ok {
			fmt.Fprintln(stdout, RegisteredSpecMessage(completer.Path, command))
			return
		}
		fmt.Fprintln(stderr, NoCompletionSpecMessage(command))
	case "-C":
		if len(args) < 3 {
			return
		}
		registeredCompleters[args[2]] = Completer{
			Path: args[1],
			Func: runCompleter,
		}
	}
}
