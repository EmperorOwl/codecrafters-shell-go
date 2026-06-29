package builtins

import (
	"bytes"
	"testing"
)

func TestComplete(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr string
	}{
		{
			name:    "prints missing specification for -p",
			args:    []string{"-p", "git"},
			wantErr: "complete: git: no completion specification\n",
		},
		{
			name:    "ignores -C registration for now",
			args:    []string{"-C", "/path/to/script", "git"},
			wantErr: "",
		},
		{
			name:    "ignores bare complete",
			args:    nil,
			wantErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stderr bytes.Buffer
			Complete(&stderr, tt.args)
			if got := stderr.String(); got != tt.wantErr {
				t.Errorf("Complete(%v) stderr = %q, want %q", tt.args, got, tt.wantErr)
			}
		})
	}
}
