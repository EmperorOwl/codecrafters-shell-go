package session

import (
	"os"

	"github.com/codecrafters-io/shell-starter-go/app/completion"
	"github.com/codecrafters-io/shell-starter-go/app/history"
	"github.com/codecrafters-io/shell-starter-go/app/jobs"
	"github.com/codecrafters-io/shell-starter-go/app/variables"
)

// State holds mutable shell state for the lifetime of a shell session.
type State struct {
	Jobs       *jobs.Table
	History    *history.List
	Histfile   string
	Completion *completion.Registry
	Variables  *variables.Store
}

// NewState returns a fresh shell session with empty jobs, variables, and completion registry.
func NewState() *State {
	histfile := os.Getenv("HISTFILE")
	list := history.NewList()
	_ = list.AppendFromFile(histfile)

	return &State{
		Jobs:       jobs.NewTable(),
		History:    list,
		Histfile:   histfile,
		Completion: completion.NewRegistry(),
		Variables:  variables.NewStore(),
	}
}
