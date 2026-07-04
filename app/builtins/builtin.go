package builtins

import (
	"io"
	"slices"

	"github.com/codecrafters-io/shell-starter-go/app/jobs"
)

var shellBuiltins = map[string]struct{}{
	"cd":       {},
	"complete": {},
	"echo":     {},
	"exit":     {},
	"jobs":     {},
	"pwd":      {},
	"type":     {},
}

// Builtin describes a shell builtin command.
type Builtin struct {
	Name   string
	Args   []string
	Stdout io.Writer
	Stderr io.Writer

	completers map[string]string
	jobs       *jobs.JobTable
}

// IsBuiltin reports whether name is a shell builtin.
func IsBuiltin(name string) bool {
	_, ok := shellBuiltins[name]
	return ok
}

// Names returns sorted builtin command names.
func Names() []string {
	names := make([]string, 0, len(shellBuiltins))
	for name := range shellBuiltins {
		names = append(names, name)
	}
	slices.Sort(names)
	return names
}

// New returns a Builtin for name and args.
func New(name string, args []string, stdout, stderr io.Writer, completers map[string]string, jobs *jobs.JobTable) *Builtin {
	return &Builtin{
		Name:       name,
		Args:       args,
		Stdout:     stdout,
		Stderr:     stderr,
		completers: completers,
		jobs:       jobs,
	}
}

// Run executes the builtin. The first return value is true when the shell should exit.
func (b *Builtin) Run() (exitShell bool, err error) {
	switch b.Name {
	case "exit":
		return true, nil
	case "echo":
		Echo(b.Stdout, b.Args)
	case "pwd":
		Pwd(b.Stdout)
	case "cd":
		directory := ""
		if len(b.Args) > 0 {
			directory = b.Args[0]
		}
		Cd(b.Stderr, directory)
	case "type":
		target := ""
		if len(b.Args) > 0 {
			target = b.Args[0]
		}
		Type(b.Stdout, target, IsBuiltin(target))
	case "complete":
		Complete(b.Stdout, b.Stderr, b.Args, b.completers)
	case "jobs":
		Jobs(b.Stdout, b.jobs)
	}
	return false, nil
}
