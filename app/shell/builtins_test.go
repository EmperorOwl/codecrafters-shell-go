package shell

import (
	"bytes"
	"strings"
	"testing"
)

func TestHandleBuiltin(t *testing.T) {
	tests := []struct {
		name       string
		command    string
		shouldExit bool
	}{
		{name: "exit terminates shell", command: "exit", shouldExit: true},
		{name: "unknown command is not builtin", command: "xyz", shouldExit: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HandleBuiltin(tt.command); got != tt.shouldExit {
				t.Errorf("HandleBuiltin(%q) = %v, want %v", tt.command, got, tt.shouldExit)
			}
		})
	}
}

func TestShellRunExitBuiltin(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "exit after invalid command",
			input: "invalid_command_1\nexit\n",
			want:  "$ invalid_command_1: command not found\n$ ",
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
