package terminal

// HistoryHandler supplies previous commands for up-arrow recall.
type HistoryHandler interface {
	HistoryPrevious(stepsBack int) (string, bool)
}

// historyBrowseState tracks how far back the user has navigated in history.
type historyBrowseState struct {
	stepsBack int
}

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

func (s *historyBrowseState) reset() {
	s.stepsBack = 0
}
