package executor

import (
	"io"

	"github.com/codecrafters-io/shell-starter-go/app/completion"
	"github.com/codecrafters-io/shell-starter-go/app/jobs"
)

// Executor runs parsed commands using injected shell state.
type Executor struct {
	jobTable           *jobs.JobTable
	completionRegistry *completion.CompletionRegistry
	stdin              io.Reader
}

// New returns an executor wired to the given shell state.
func New(jobTable *jobs.JobTable, completionRegistry *completion.CompletionRegistry) *Executor {
	return &Executor{
		jobTable:           jobTable,
		completionRegistry: completionRegistry,
	}
}

// SetStdin configures stdin used for foreground command execution.
func (e *Executor) SetStdin(stdin io.Reader) {
	e.stdin = stdin
}
