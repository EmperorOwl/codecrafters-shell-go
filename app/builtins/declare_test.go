package builtins

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDeclare(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantOut string
		wantErr string
	}{
		{
			name:    "prints not found for -p",
			args:    []string{"-p", "missing_variable"},
			wantErr: "declare: missing_variable: not found\n",
		},
		{
			name: "ignores bare declare",
			args: nil,
		},
		{
			name: "ignores -p without variable name",
			args: []string{"-p"},
		},
		{
			name: "ignores assignment form",
			args: []string{"variable=value"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			Declare(&stdout, &stderr, tt.args)

			if diff := cmp.Diff(tt.wantOut, stdout.String()); diff != "" {
				t.Errorf("Declare(%v) stdout mismatch (-want +got):\n%s", tt.args, diff)
			}
			if diff := cmp.Diff(tt.wantErr, stderr.String()); diff != "" {
				t.Errorf("Declare(%v) stderr mismatch (-want +got):\n%s", tt.args, diff)
			}
		})
	}
}
