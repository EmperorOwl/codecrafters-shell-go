package builtins

import (
	"bytes"
	"testing"
)

func resetCompletionSpecs() {
	completionSpecs = map[string]string{}
}

func TestComplete(t *testing.T) {
	tests := []struct {
		name    string
		setup   func()
		args    []string
		wantOut string
		wantErr string
	}{
		{
			name:    "prints missing specification for -p",
			args:    []string{"-p", "git"},
			wantErr: "complete: git: no completion specification\n",
		},
		{
			name: "registers completion with -C",
			setup: func() {
				Complete(nil, nil, []string{"-C", "/path/to/script", "git"})
			},
			args:    []string{"-p", "git"},
			wantOut: "complete -C '/path/to/script' git\n",
		},
		{
			name: "registers and displays docker completion",
			setup: func() {
				Complete(nil, nil, []string{"-C", "/path/to/docker/completer", "docker"})
			},
			args:    []string{"-p", "docker"},
			wantOut: "complete -C '/path/to/docker/completer' docker\n",
		},
		{
			name:    "ignores -C with too few arguments",
			args:    []string{"-C", "/path/to/script"},
			wantOut: "",
			wantErr: "",
		},
		{
			name:    "ignores bare complete",
			args:    nil,
			wantOut: "",
			wantErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetCompletionSpecs()
			if tt.setup != nil {
				tt.setup()
			}

			var stdout, stderr bytes.Buffer
			Complete(&stdout, &stderr, tt.args)
			if got := stdout.String(); got != tt.wantOut {
				t.Errorf("Complete(%v) stdout = %q, want %q", tt.args, got, tt.wantOut)
			}
			if got := stderr.String(); got != tt.wantErr {
				t.Errorf("Complete(%v) stderr = %q, want %q", tt.args, got, tt.wantErr)
			}
		})
	}
}
