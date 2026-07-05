package shell

import (
	"slices"

	"github.com/codecrafters-io/shell-starter-go/app/terminal"
)

func (s *Shell) HandleTab(state *terminal.TabState, buffer string) terminal.TabResult {
	newBuffer, listings := s.completeBuffer(buffer)
	return applyTabAction(state, buffer, newBuffer, listings)
}

func applyTabAction(state *terminal.TabState, buffer string, newBuffer string, listings []string) terminal.TabResult {
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
