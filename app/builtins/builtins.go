package builtins

import (
	"io"
	"slices"

	"github.com/codecrafters-io/shell-starter-go/app/completion"
	"github.com/codecrafters-io/shell-starter-go/app/jobs"
)

// State holds stable shell session state shared across command invocations.
type State struct {
	Jobs       *jobs.JobTable
	Completion *completion.CompletionRegistry
}

// Context holds per-invocation I/O and shell state for a builtin command.
type Context struct {
	Stdout io.Writer
	Stderr io.Writer
	State  *State
}

type Handler func(ctx *Context, args []string) (exit bool, err error)

var handlers map[string]Handler

func init() {
	handlers = map[string]Handler{
		"cd":       cdBuiltin,
		"complete": completeBuiltin,
		"echo":     echoBuiltin,
		"exit":     exitBuiltin,
		"jobs":     jobsBuiltin,
		"pwd":      pwdBuiltin,
		"type":     typeBuiltin,
	}
}

func IsBuiltin(name string) bool {
	_, ok := handlers[name]
	return ok
}

func Names() []string {
	names := make([]string, 0, len(handlers))
	for name := range handlers {
		names = append(names, name)
	}
	slices.Sort(names)
	return names
}

func Run(name string, args []string, ctx *Context) (bool, error) {
	handler, ok := handlers[name]
	if !ok {
		return false, nil
	}

	return handler(ctx, args)
}

func exitBuiltin(ctx *Context, args []string) (bool, error) {
	return true, nil
}
