package shell

import (
	"testing"

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
		want            TabResult
		wantPending     []string
	}{
		{
			name:   "completes unique prefix",
			buffer: "ech",
			want:   TabResult{Buffer: "echo "},
		},
		{
			name:        "rings bell on ambiguous prefix",
			buffer:      "e",
			want:        TabResult{Buffer: "e", RingBell: true},
			wantPending: []string{"echo", "exit"},
		},
		{
			name:            "shows listings on second tab",
			buffer:          "e",
			pendingListings: []string{"echo", "exit"},
			want:            TabResult{Buffer: "e", ListingsToShow: []string{"echo", "exit"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := TabState{pendingListings: tt.pendingListings}
			got := ApplyTabAction(&state, tt.buffer, builtins, tt.executables, nil, nil)
			if diff := cmp.Diff(tt.want, got, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("ApplyTabAction() mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tt.wantPending, state.pendingListings, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("pending listings mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
