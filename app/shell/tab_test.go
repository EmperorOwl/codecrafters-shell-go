package shell

import (
	"testing"

	"github.com/codecrafters-io/shell-starter-go/app/terminal"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestApplyTabAction(t *testing.T) {
	tests := []struct {
		name            string
		buffer          string
		newBuffer       string
		listings        []string
		pendingListings []string
		want            terminal.TabResult
		wantPending     []string
	}{
		{
			name:      "updates buffer on unique match",
			buffer:    "ech",
			newBuffer: "echo ",
			want:      terminal.TabResult{Buffer: "echo "},
		},
		{
			name:        "rings bell on ambiguous prefix",
			buffer:      "e",
			newBuffer:   "e",
			listings:    []string{"echo", "exit"},
			want:        terminal.TabResult{Buffer: "e", RingBell: true},
			wantPending: []string{"echo", "exit"},
		},
		{
			name:            "shows listings on second tab",
			buffer:          "e",
			newBuffer:       "e",
			listings:        []string{"echo", "exit"},
			pendingListings: []string{"echo", "exit"},
			want:            terminal.TabResult{Buffer: "e", ListingsToShow: []string{"echo", "exit"}},
		},
		{
			name:      "rings bell when nothing changes",
			buffer:    "xyz",
			newBuffer: "xyz",
			want:      terminal.TabResult{Buffer: "xyz", RingBell: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := terminal.TabState{PendingListings: tt.pendingListings}
			got := applyTabAction(&state, tt.buffer, tt.newBuffer, tt.listings)
			if diff := cmp.Diff(tt.want, got, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("applyTabAction() mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tt.wantPending, state.PendingListings, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("pending listings mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
