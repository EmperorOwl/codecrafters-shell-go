package shell

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestBuildCompleterFuncOptions(t *testing.T) {
	tests := []struct {
		name   string
		buffer string
		want   CompleterFuncOptions
	}{
		{
			name:   "first argument completion",
			buffer: "git ",
			want: CompleterFuncOptions{
				Command:   "git",
				CompLine:  "git ",
				CompPoint: 4,
			},
		},
		{
			name:   "partial first argument",
			buffer: "git remot",
			want: CompleterFuncOptions{
				Command:      "git",
				CurrentWord:  "remot",
				PreviousWord: "git",
				CompLine:     "git remot",
				CompPoint:    9,
			},
		},
		{
			name:   "later argument completion",
			buffer: "git remote set",
			want: CompleterFuncOptions{
				Command:      "git",
				CurrentWord:  "set",
				PreviousWord: "remote",
				CompLine:     "git remote set",
				CompPoint:    14,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildCompleterFuncOptions(tt.buffer)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("buildCompleterFuncOptions(%q) mismatch (-want +got):\n%s", tt.buffer, diff)
			}
		})
	}
}
