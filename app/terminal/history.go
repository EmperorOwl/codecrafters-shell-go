package terminal

// HistoryHandler supplies previous commands for up-arrow recall.
type HistoryHandler interface {
	// HistoryPrevious returns the command stepsBack entries before the most
	// recent one. stepsBack 0 is the newest command.
	HistoryPrevious(stepsBack int) (string, bool)
}

// historyBrowseState tracks how far back the user has navigated while browsing
// history on the current prompt line.
//
// stepsBack is 0 on a fresh line. Each up-arrow increments it; each down-arrow
// decrements it back toward 0 (an empty line).
type historyBrowseState struct {
	stepsBack int
}

// stepUp moves one command further back in history.
func (s *historyBrowseState) stepUp(handler HistoryHandler) (command string, ok bool) {
	if handler == nil {
		return "", false
	}

	s.stepsBack++
	command, ok = handler.HistoryPrevious(s.stepsBack - 1)
	if !ok {
		s.stepsBack--
	}
	return command, ok
}

// stepDown moves one command forward toward the present.
// At stepsBack 0, down-arrow is a no-op. Reaching 0 returns an empty line.
func (s *historyBrowseState) stepDown(handler HistoryHandler) (command string, ok bool) {
	if s.stepsBack <= 0 {
		return "", false
	}

	s.stepsBack--
	if s.stepsBack == 0 {
		return "", true
	}
	if handler == nil {
		return "", false
	}
	return handler.HistoryPrevious(s.stepsBack - 1)
}

// reset clears browse position when the user edits the line manually.
func (s *historyBrowseState) reset() {
	s.stepsBack = 0
}
