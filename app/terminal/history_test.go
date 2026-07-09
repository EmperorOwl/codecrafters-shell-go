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

			if ok != tt.wantOK {
				t.Fatalf("stepUp() ok = %v, want %v", ok, tt.wantOK)
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

func TestHistoryBrowseState_stepUpAtStart(t *testing.T) {
	handler := stubHistoryHandler{commands: []string{"echo hello"}}
	var state historyBrowseState

	if _, ok := state.stepUp(handler); !ok {
		t.Fatal("stepUp() ok = false, want true")
	}
	if _, ok := state.stepUp(handler); ok {
		t.Fatal("stepUp() ok = true, want false at start of history")
	}
	if diff := cmp.Diff(1, state.stepsBack); diff != "" {
		t.Errorf("stepsBack mismatch (-want +got):\n%s", diff)
	}
}
