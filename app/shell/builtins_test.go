package shell

import (
	"bytes"
	"testing"
)

func TestTryBuiltin(t *testing.T) {
	tests := []struct {
		name        string
		line        string
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
			name:        "echo is handled",
			line:        "echo hello world",
			wantHandled: true,
		},
		{
			name:        "pwd is handled",
			line:        "pwd",
			wantHandled: true,
		},
		{
			name:        "cd is handled",
			line:        "cd /tmp",
			wantHandled: true,
		},
		{
			name:        "type is handled",
			line:        "type echo",
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
		})
	}
}
