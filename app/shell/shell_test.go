package shell

import (
	"bytes"
	"strings"
	"testing"
)

func TestShellRun(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "echo builtin",
			input: "echo hello\n",
			want:  "$ hello\n$ ",
		},
		{
			name:  "echo hello with EOF does not print prompt again",
			input: "echo hello",
			want:  "$ hello\n",
		},
		{
			name:  "unknown command",
			input: "xyz\n",
			want:  "$ xyz: command not found\n$ ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shell := New()
			var out bytes.Buffer
			err := shell.Run(strings.NewReader(tt.input), &out)
			if err != nil {
				t.Fatalf("Run() error = %v", err)
			}
			if got := out.String(); got != tt.want {
				t.Errorf("Run() output = %q, want %q", got, tt.want)
			}
		})
	}
}
