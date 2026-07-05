package terminal

// TabState tracks double-tab listing behavior while reading one input line.
type TabState struct {
	PendingListings []string
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
