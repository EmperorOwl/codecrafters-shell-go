package builtins

import (
	"io"
	"slices"

	"github.com/codecrafters-io/shell-starter-go/app/jobs"
)

type Context struct {
	Stdout io.Writer
	Stderr io.Writer

	Completers map[string]string
	Jobs       *jobs.JobTable
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
