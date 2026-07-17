package builtins

import (
	"io"
	"slices"

	"github.com/codecrafters-io/shell-starter-go/app/session"
)

// Context holds per-invocation I/O and REPL state for a builtin command.
type Context struct {
	Stdout io.Writer
	Stderr io.Writer
	Session *session.Session
}

// Handler runs a builtin command. The bool is true when the shell should exit.
type Handler func(ctx *Context, args []string) (exit bool, err error)

type registry struct {
	handlers map[string]Handler
}

var defaultRegistry = newRegistry()

func newRegistry() *registry {
	return &registry{
		handlers: make(map[string]Handler),
	}
}

func (r *registry) register(name string, handler Handler) {
	r.handlers[name] = handler
}

func (r *registry) run(name string, args []string, ctx *Context) (bool, error) {
	handler, ok := r.handlers[name]
	if !ok {
		return false, nil
	}
	return handler(ctx, args)
}

func (r *registry) is(name string) bool {
	_, ok := r.handlers[name]
	return ok
}

func (r *registry) names() []string {
	names := make([]string, 0, len(r.handlers))
	for name := range r.handlers {
		names = append(names, name)
	}
	slices.Sort(names)
	return names
}

func register(name string, handler Handler) {
	defaultRegistry.register(name, handler)
}

// Run executes a registered builtin on the default registry.
func Run(name string, args []string, ctx *Context) (bool, error) {
	return defaultRegistry.run(name, args, ctx)
}

// IsBuiltin reports whether name is a registered builtin command.
func IsBuiltin(name string) bool {
	return defaultRegistry.is(name)
}

// Names returns registered builtin names in sorted order.
func Names() []string {
	return defaultRegistry.names()
}
