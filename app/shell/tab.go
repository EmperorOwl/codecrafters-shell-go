package shell

import (
	"slices"

	"github.com/codecrafters-io/shell-starter-go/app/builtins"
	"github.com/codecrafters-io/shell-starter-go/app/completion"
	"github.com/codecrafters-io/shell-starter-go/app/external"
	"github.com/codecrafters-io/shell-starter-go/app/terminal"
)

func (s *Shell) HandleTab(state *terminal.TabState, buffer string) terminal.TabResult {
	return applyTabAction(
		state,
		buffer,
		builtins.Names(),
		external.FindAllExecutablesInPath(),
		s.listFiles,
		s.programmableComplete,
	)
}

func applyTabAction(
	state *terminal.TabState,
	buffer string,
	builtinsList, executables []string,
	listFiles completion.FileLister,
	completeHandler completion.CompleteHandler,
) terminal.TabResult {
	newBuffer, listings := completion.ApplyTab(
		builtinsList,
		executables,
		listFiles,
		completeHandler,
		buffer,
	)

	if len(listings) > 0 {
		if slices.Equal(state.PendingListings, listings) {
			state.PendingListings = nil
			return terminal.TabResult{Buffer: buffer, ListingsToShow: listings}
		}
		state.PendingListings = listings
		return terminal.TabResult{Buffer: buffer, RingBell: true}
	}

	state.PendingListings = nil
	if newBuffer != buffer {
		return terminal.TabResult{Buffer: newBuffer}
	}
	return terminal.TabResult{Buffer: buffer, RingBell: true}
}
