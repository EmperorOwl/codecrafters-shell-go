package terminal

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

type stubHistoryHandler struct {
	commands []string
}

func (h stubHistoryHandler) HistoryPrevious(stepsBack int) (string, bool) {
	if stepsBack < 0 || stepsBack >= len(h.commands) {
		return "", false
	}
	return h.commands[len(h.commands)-1-stepsBack], true
}

func TestHistoryBrowseState_stepUp(t *testing.T) {
	handler := stubHistoryHandler{commands: []string{"echo hello", "echo world"}}

	tests := []struct {
		name      string
		stepCount int
		want      string
		wantOK    bool
		wantSteps int
	}{
		{
			name:      "first step shows most recent command",
			stepCount: 1,
			want:      "echo world",
			wantOK:    true,
			wantSteps: 1,
		},
		{
			name:      "second step shows earlier command",
			stepCount: 2,
			want:      "echo hello",
			wantOK:    true,
			wantSteps: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var state historyBrowseState
			var got string
			var ok bool

			for range tt.stepCount {
				got, ok = state.stepUp(handler)
			}

			if diff := cmp.Diff(tt.wantOK, ok); diff != "" {
				t.Fatalf("stepUp() ok mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("stepUp() command mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tt.wantSteps, state.stepsBack); diff != "" {
				t.Errorf("stepsBack mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestHistoryBrowseState_stepUpEdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		handler   HistoryHandler
		stepCount int
		wantOK    bool
		wantSteps int
	}{
		{
			name:      "cannot step past start of history",
			handler:   stubHistoryHandler{commands: []string{"echo hello"}},
			stepCount: 2,
			wantOK:    false,
			wantSteps: 1,
		},
		{
			name:      "nil handler rejects step",
			handler:   nil,
			stepCount: 1,
			wantOK:    false,
			wantSteps: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var state historyBrowseState
			var ok bool

			for range tt.stepCount {
				_, ok = state.stepUp(tt.handler)
			}

			if diff := cmp.Diff(tt.wantOK, ok); diff != "" {
				t.Errorf("stepUp() ok mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tt.wantSteps, state.stepsBack); diff != "" {
				t.Errorf("stepsBack mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestHistoryBrowseState_stepDown(t *testing.T) {
	tests := []struct {
		name              string
		commands          []string
		setup             func(*historyBrowseState, HistoryHandler)
		stepDownNilHandler bool
		want              string
		wantOK            bool
		wantSteps         int
	}{
		{
			name:     "steps forward one command",
			commands: []string{"echo hello", "echo world"},
			setup: func(state *historyBrowseState, handler HistoryHandler) {
				state.stepUp(handler)
				state.stepUp(handler)
			},
			want:      "echo world",
			wantOK:    true,
			wantSteps: 1,
		},
		{
			name:     "returns empty line at present",
			commands: []string{"echo hello", "echo world"},
			setup: func(state *historyBrowseState, handler HistoryHandler) {
				state.stepUp(handler)
			},
			want:      "",
			wantOK:    true,
			wantSteps: 0,
		},
		{
			name:      "no-op at bottom of history",
			commands:  []string{"echo hello"},
			setup:     func(*historyBrowseState, HistoryHandler) {},
			wantOK:    false,
			wantSteps: 0,
		},
		{
			name:     "nil handler after partial step down returns false",
			commands: []string{"echo hello", "echo world"},
			setup: func(state *historyBrowseState, handler HistoryHandler) {
				state.stepUp(handler)
				state.stepUp(handler)
			},
			stepDownNilHandler: true,
			wantOK:             false,
			wantSteps:          1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := stubHistoryHandler{commands: tt.commands}
			var state historyBrowseState
			if tt.setup != nil {
				tt.setup(&state, handler)
			}

			var stepDownHandler HistoryHandler = handler
			if tt.stepDownNilHandler {
				stepDownHandler = nil
			}

			got, ok := state.stepDown(stepDownHandler)
			if diff := cmp.Diff(tt.wantOK, ok); diff != "" {
				t.Fatalf("stepDown() ok mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("stepDown() command mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tt.wantSteps, state.stepsBack); diff != "" {
				t.Errorf("stepsBack mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestHistoryBrowseState_reset(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(*historyBrowseState, HistoryHandler)
		wantSteps int
	}{
		{
			name: "clears browse position",
			setup: func(state *historyBrowseState, handler HistoryHandler) {
				state.stepUp(handler)
				state.stepUp(handler)
				state.reset()
			},
			wantSteps: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := stubHistoryHandler{commands: []string{"echo hello", "echo world"}}
			var state historyBrowseState
			if tt.setup != nil {
				tt.setup(&state, handler)
			}

			if diff := cmp.Diff(tt.wantSteps, state.stepsBack); diff != "" {
				t.Errorf("reset() stepsBack mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
