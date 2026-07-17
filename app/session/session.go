package session

import (
	"github.com/codecrafters-io/shell-starter-go/app/completion"
	"github.com/codecrafters-io/shell-starter-go/app/history"
	"github.com/codecrafters-io/shell-starter-go/app/jobs"
	"github.com/codecrafters-io/shell-starter-go/app/variables"
)

// Session holds mutable shell state for the lifetime of a shell session.
type Session struct {
	Jobs       *jobs.Table
	History    *history.List
	Histfile   string
	Completion *completion.Registry
	Variables  *variables.Store
}

// NewSession returns a fresh shell session with empty jobs, variables, and completion registry.
func NewSession() *Session {
	return &Session{
		Jobs:       jobs.NewTable(),
		History:    history.NewList(),
		Completion: completion.NewRegistry(),
		Variables:  variables.NewStore(),
	}
}
