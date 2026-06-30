package completion

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestApplyCommandTab(t *testing.T) {
	builtins := []string{"cd", "echo", "exit", "pwd", "type"}

	tests := []struct {
		name         string
		executables  []string
		buffer       string
		wantBuffer   string
		wantListings []string
	}{
		{
			name:       "completes echo",
			buffer:     "ech",
			wantBuffer: "echo ",
		},
		{
			name:       "completes exit",
			buffer:     "exi",
			wantBuffer: "exit ",
		},
		{
			name:       "no match",
			buffer:     "xyz",
			wantBuffer: "xyz",
		},
		{
			name:         "ambiguous prefix lists matches",
			buffer:       "e",
			wantBuffer:   "e",
			wantListings: []string{"echo", "exit"},
		},
		{
			name:         "empty buffer lists all builtins",
			buffer:       "",
			wantBuffer:   "",
			wantListings: builtins,
		},
		{
			name:        "completes executable",
			executables: []string{"custom_executable"},
			buffer:      "custom",
			wantBuffer:  "custom_executable ",
		},
		{
			name:         "ambiguous executable prefix lists matches",
			executables:  []string{"xyz_bar", "xyz_baz", "xyz_quz"},
			buffer:       "xyz_",
			wantBuffer:   "xyz_",
			wantListings: []string{"xyz_bar", "xyz_baz", "xyz_quz"},
		},
		{
			name:        "completes to longest common prefix",
			executables: []string{"xyz_foo", "xyz_foo_bar", "xyz_foo_bar_baz"},
			buffer:      "xyz_",
			wantBuffer:  "xyz_foo",
		},
		{
			name:        "completes to next longest common prefix",
			executables: []string{"xyz_foo", "xyz_foo_bar", "xyz_foo_bar_baz"},
			buffer:      "xyz_foo_",
			wantBuffer:  "xyz_foo_bar",
		},
		{
			name:        "completes final unique executable",
			executables: []string{"xyz_foo", "xyz_foo_bar", "xyz_foo_bar_baz"},
			buffer:      "xyz_foo_bar_",
			wantBuffer:  "xyz_foo_bar_baz ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBuffer, gotListings := applyCommandTab(builtins, tt.executables, tt.buffer)
			if diff := cmp.Diff(tt.wantBuffer, gotBuffer); diff != "" {
				t.Errorf("applyCommandTab(%q) buffer mismatch (-want +got):\n%s", tt.buffer, diff)
			}
			if diff := cmp.Diff(tt.wantListings, gotListings, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("applyCommandTab(%q) listings mismatch (-want +got):\n%s", tt.buffer, diff)
			}
		})
	}
}
