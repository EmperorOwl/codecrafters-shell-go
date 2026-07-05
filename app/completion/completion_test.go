package completion

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestComplete(t *testing.T) {
	builtins := []string{"cd", "echo", "exit", "pwd", "type"}

	tests := []struct {
		name         string
		candidates   []string
		prefix       string
		wantToken    string
		wantListings []string
		wantUnique   bool
	}{
		{
			name:       "completes echo",
			candidates: builtins,
			prefix:     "ech",
			wantToken:  "echo",
			wantUnique: true,
		},
		{
			name:       "completes exit",
			candidates: builtins,
			prefix:     "exi",
			wantToken:  "exit",
			wantUnique: true,
		},
		{
			name:       "no match",
			candidates: builtins,
			prefix:     "xyz",
			wantToken:  "xyz",
		},
		{
			name:         "ambiguous prefix lists matches",
			candidates:   builtins,
			prefix:       "e",
			wantToken:    "e",
			wantListings: []string{"echo", "exit"},
		},
		{
			name:         "empty prefix lists all candidates",
			candidates:   builtins,
			prefix:       "",
			wantToken:    "",
			wantListings: builtins,
		},
		{
			name:       "completes executable",
			candidates: []string{"custom_executable"},
			prefix:     "custom",
			wantToken:  "custom_executable",
			wantUnique: true,
		},
		{
			name:         "ambiguous executable prefix lists matches",
			candidates:   []string{"xyz_bar", "xyz_baz", "xyz_quz"},
			prefix:       "xyz_",
			wantToken:    "xyz_",
			wantListings: []string{"xyz_bar", "xyz_baz", "xyz_quz"},
		},
		{
			name:       "completes to longest common prefix",
			candidates: []string{"xyz_foo", "xyz_foo_bar", "xyz_foo_bar_baz"},
			prefix:     "xyz_",
			wantToken:  "xyz_foo",
		},
		{
			name:       "completes to next longest common prefix",
			candidates: []string{"xyz_foo", "xyz_foo_bar", "xyz_foo_bar_baz"},
			prefix:     "xyz_foo_",
			wantToken:  "xyz_foo_bar",
		},
		{
			name:       "completes final unique match",
			candidates: []string{"xyz_foo", "xyz_foo_bar", "xyz_foo_bar_baz"},
			prefix:     "xyz_foo_bar_",
			wantToken:  "xyz_foo_bar_baz",
			wantUnique: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotToken, gotListings, gotUnique := Complete(tt.prefix, tt.candidates)
			if diff := cmp.Diff(tt.wantToken, gotToken); diff != "" {
				t.Errorf("Complete(%q) token mismatch (-want +got):\n%s", tt.prefix, diff)
			}
			if diff := cmp.Diff(tt.wantListings, gotListings, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("Complete(%q) listings mismatch (-want +got):\n%s", tt.prefix, diff)
			}
			if diff := cmp.Diff(tt.wantUnique, gotUnique); diff != "" {
				t.Errorf("Complete(%q) unique mismatch (-want +got):\n%s", tt.prefix, diff)
			}
		})
	}
}
