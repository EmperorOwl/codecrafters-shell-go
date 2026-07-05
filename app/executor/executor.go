package executor

import (
	"github.com/codecrafters-io/shell-starter-go/app/completion"
	"github.com/codecrafters-io/shell-starter-go/app/jobs"
)

// Executor runs parsed commands using injected shell state.
type Executor struct {
	jobManager         *jobs.JobManager
	completionRegistry *completion.CompletionRegistry
}

// New returns an executor wired to the given shell state.
func New(jobManager *jobs.JobManager, completionRegistry *completion.CompletionRegistry) *Executor {
	return &Executor{
		jobManager:         jobManager,
		completionRegistry: completionRegistry,
	}
}
