package completion

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestParseCompleterOutput(t *testing.T) {
	tests := []struct {
		name   string
		output []byte
		want   []string
	}{
		{
			name:   "empty output",
			output: nil,
			want:   nil,
		},
		{
			name:   "single candidate",
			output: []byte("run\n"),
			want:   []string{"run"},
		},
		{
			name:   "multiple candidates",
			output: []byte("stash\nstatus\n"),
			want:   []string{"stash", "status"},
		},
		{
			name:   "trims trailing newline",
			output: []byte("run"),
			want:   []string{"run"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseCompleterOutput(tt.output)
			if diff := cmp.Diff(tt.want, got, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("parseCompleterOutput() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
