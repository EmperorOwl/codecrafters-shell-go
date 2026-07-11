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
	State  *session.State
}

// Handler runs a builtin command. The bool is true when the shell should exit.
type Handler func(ctx *Context, args []string) (exit bool, err error)

// Registry stores registered builtin command handlers.
type Registry struct {
	handlers map[string]Handler
}

var defaultRegistry = NewRegistry()

// NewRegistry returns an empty builtin registry.
func NewRegistry() *Registry {
	return &Registry{
		handlers: make(map[string]Handler),
	}
}

// Register adds a builtin command handler to the registry.
func (r *Registry) Register(name string, handler Handler) {
	r.handlers[name] = handler
}

// Run executes a registered builtin. The bool is true when the shell should exit.
func (r *Registry) Run(name string, args []string, ctx *Context) (bool, error) {
	handler, ok := r.handlers[name]
	if !ok {
		return false, nil
	}
	return handler(ctx, args)
}

// Is reports whether name is a registered builtin command.
func (r *Registry) Is(name string) bool {
	_, ok := r.handlers[name]
	return ok
}

// Names returns registered builtin names in sorted order.
func (r *Registry) Names() []string {
	names := make([]string, 0, len(r.handlers))
	for name := range r.handlers {
		names = append(names, name)
	}
	slices.Sort(names)
	return names
}

func register(name string, handler Handler) {
	defaultRegistry.Register(name, handler)
}

// Run executes a registered builtin on the default registry.
func Run(name string, args []string, ctx *Context) (bool, error) {
	return defaultRegistry.Run(name, args, ctx)
}

// IsBuiltin reports whether name is a registered builtin command.
func IsBuiltin(name string) bool {
	return defaultRegistry.Is(name)
}

// Names returns registered builtin names in sorted order.
func Names() []string {
	return defaultRegistry.Names()
}
