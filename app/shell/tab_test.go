package shell

import (
	"testing"

	"github.com/codecrafters-io/shell-starter-go/app/terminal"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestHandleTab(t *testing.T) {
	builtins := []string{"echo", "exit"}

	tests := []struct {
		name            string
		buffer          string
		pendingListings []string
		executables     []string
		want            terminal.TabResult
		wantPending     []string
	}{
		{
			name:   "completes unique prefix",
			buffer: "ech",
			want:   terminal.TabResult{Buffer: "echo "},
		},
		{
			name:        "rings bell on ambiguous prefix",
			buffer:      "e",
			want:        terminal.TabResult{Buffer: "e", RingBell: true},
			wantPending: []string{"echo", "exit"},
		},
		{
			name:            "shows listings on second tab",
			buffer:          "e",
			pendingListings: []string{"echo", "exit"},
			want:            terminal.TabResult{Buffer: "e", ListingsToShow: []string{"echo", "exit"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := terminal.TabState{PendingListings: tt.pendingListings}
			got := applyTabAction(&state, tt.buffer, builtins, tt.executables, nil, nil)
			if diff := cmp.Diff(tt.want, got, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("applyTabAction() mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tt.wantPending, state.PendingListings, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("pending listings mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
