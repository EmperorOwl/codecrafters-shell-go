package repl

import (
	"github.com/codecrafters-io/shell-starter-go/app/completion"
	"github.com/codecrafters-io/shell-starter-go/app/history"
	"github.com/codecrafters-io/shell-starter-go/app/jobs"
)

// State holds mutable shell state for the lifetime of the REPL loop.
type State struct {
	Jobs       *jobs.JobTable
	History    *history.HistoryList
	Completion *completion.CompletionRegistry
}

// NewState returns a fresh REPL state with an empty job table and completion registry.
func NewState() *State {
	return &State{
		Jobs:       &jobs.JobTable{},
		History:    &history.HistoryList{},
		Completion: completion.NewCompletionRegistry(),
	}
}
