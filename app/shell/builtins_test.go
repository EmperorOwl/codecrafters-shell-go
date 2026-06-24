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

func TestTypeOutput(t *testing.T) {
	tests := []struct {
		name    string
		command string
		want    string
	}{
		{name: "echo builtin", command: "echo", want: "echo is a shell builtin"},
		{name: "exit builtin", command: "exit", want: "exit is a shell builtin"},
		{name: "type builtin", command: "type", want: "type is a shell builtin"},
		{name: "invalid command", command: "invalid_command", want: "invalid_command: not found"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TypeOutput(tt.command); got != tt.want {
				t.Errorf("TypeOutput(%q) = %q, want %q", tt.command, got, tt.want)
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
			name:        "type reports builtin",
			line:        "type echo",
			wantOutput:  "echo is a shell builtin\n",
			wantHandled: true,
		},
		{
			name:        "type reports not found",
			line:        "type invalid_command",
			wantOutput:  "invalid_command: not found\n",
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

func TestShellRunTypeBuiltin(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "type builtins and invalid command",
			input: "type echo\ntype exit\ntype type\ntype invalid_command\n",
			want:  "$ echo is a shell builtin\n$ exit is a shell builtin\n$ type is a shell builtin\n$ invalid_command: not found\n$ ",
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
