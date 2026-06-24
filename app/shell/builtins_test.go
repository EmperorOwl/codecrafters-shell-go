package shell

import (
	"bytes"
	"strings"
	"testing"
)

func TestEchoOutput(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{name: "hello world", args: []string{"hello", "world"}, want: "hello world"},
		{name: "three words", args: []string{"one", "two", "three"}, want: "one two three"},
		{name: "no args", args: nil, want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EchoOutput(tt.args); got != tt.want {
				t.Errorf("EchoOutput(%v) = %q, want %q", tt.args, got, tt.want)
			}
		})
	}
}

func TestTryBuiltin(t *testing.T) {
	tests := []struct {
		name        string
		line        string
		wantOutput  string
		wantHandled bool
		wantExit    bool
	}{
		{
			name:        "exit terminates shell",
			line:        "exit",
			wantHandled: true,
			wantExit:    true,
		},
		{
			name:        "echo prints arguments",
			line:        "echo hello world",
			wantOutput:  "hello world\n",
			wantHandled: true,
		},
		{
			name:        "unknown command is not builtin",
			line:        "xyz",
			wantHandled: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out bytes.Buffer
			handled, shouldExit := TryBuiltin(tt.line, &out)
			if handled != tt.wantHandled {
				t.Errorf("TryBuiltin(%q) handled = %v, want %v", tt.line, handled, tt.wantHandled)
			}
			if shouldExit != tt.wantExit {
				t.Errorf("TryBuiltin(%q) shouldExit = %v, want %v", tt.line, shouldExit, tt.wantExit)
			}
			if got := out.String(); got != tt.wantOutput {
				t.Errorf("TryBuiltin(%q) output = %q, want %q", tt.line, got, tt.wantOutput)
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

func TestShellRunEchoBuiltin(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "multiple echo commands",
			input: "echo hello world\necho pineapple strawberry\n",
			want:  "$ hello world\n$ pineapple strawberry\n$ ",
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
