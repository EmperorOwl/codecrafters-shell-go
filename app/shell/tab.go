package shell

import (
	"slices"

	"github.com/codecrafters-io/shell-starter-go/app/builtins"
	"github.com/codecrafters-io/shell-starter-go/app/completion"
	"github.com/codecrafters-io/shell-starter-go/app/external"
)

// TabState tracks double-tab listing behavior while reading one input line.
type TabState struct {
	pendingListings []string
}

// TabResult tells the terminal how to update the input line after Tab.
type TabResult struct {
	Buffer         string
	ListingsToShow []string
	RingBell       bool
}

// TabHandler handles tab completion for an in-progress input line.
type TabHandler interface {
	HandleTab(state *TabState, buffer string) TabResult
}

func (s *Shell) HandleTab(state *TabState, buffer string) TabResult {
	return ApplyTabAction(
		state,
		buffer,
		builtins.Names(),
		external.FindAllExecutablesInPath(),
		s.listFiles,
		s.complete,
	)
}

func ApplyTabAction(
	state *TabState,
	buffer string,
	builtinsList, executables []string,
	listFiles completion.FileLister,
	completeHandler completion.CompleteHandler,
) TabResult {
	newBuffer, listings := completion.ApplyTab(
		builtinsList,
		executables,
		listFiles,
		completeHandler,
		buffer,
	)

	if len(listings) > 0 {
		if slices.Equal(state.pendingListings, listings) {
			state.pendingListings = nil
			return TabResult{Buffer: buffer, ListingsToShow: listings}
		}
		state.pendingListings = listings
		return TabResult{Buffer: buffer, RingBell: true}
	}

	state.pendingListings = nil
	if newBuffer != buffer {
		return TabResult{Buffer: newBuffer}
	}
	return TabResult{Buffer: buffer, RingBell: true}
}
