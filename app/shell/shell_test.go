package shell

import (
	"bytes"
	"strings"
	"testing"
)

func TestCommandNotFoundMessage(t *testing.T) {
	tests := []struct {
		name    string
		command string
		want    string
	}{
		{name: "simple command", command: "xyz", want: "xyz: command not found"},
		{name: "command with path", command: "bin/cmd", want: "bin/cmd: command not found"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CommandNotFoundMessage(tt.command); got != tt.want {
				t.Errorf("CommandNotFoundMessage(%q) = %q, want %q", tt.command, got, tt.want)
			}
		})
	}
}

func TestShellRunInvalidCommand(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
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
