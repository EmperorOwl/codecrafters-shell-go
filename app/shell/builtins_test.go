package shell

import (
	"bytes"
	"testing"

	"github.com/codecrafters-io/shell-starter-go/app/builtins"
)

func TestTryBuiltin(t *testing.T) {
	tests := []struct {
		name        string
		fields      []string
		wantHandled bool
		wantExit    bool
	}{
		{
			name:        "exit terminates shell",
			fields:      []string{"exit"},
			wantHandled: true,
			wantExit:    true,
		},
		{
			name:        "echo is handled",
			fields:      []string{"echo", "hello", "world"},
			wantHandled: true,
		},
		{
			name:        "pwd is handled",
			fields:      []string{"pwd"},
			wantHandled: true,
		},
		{
			name:        "cd is handled",
			fields:      []string{"cd", "/tmp"},
			wantHandled: true,
		},
		{
			name:        "type is handled",
			fields:      []string{"type", "echo"},
			wantHandled: true,
		},
		{
			name:        "complete is handled",
			fields:      []string{"complete", "-C", "/path/to/script", "git"},
			wantHandled: true,
		},
		{
			name:        "complete -p is handled",
			fields:      []string{"complete", "-p", "git"},
			wantHandled: true,
		},
		{
			name:        "unknown command is not handled",
			fields:      []string{"xyz"},
			wantHandled: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out bytes.Buffer
			var errOut bytes.Buffer
			handled, shouldExit := TryBuiltin(tt.fields, &out, &errOut, map[string]builtins.Completer{}, nil)
			if handled != tt.wantHandled {
				t.Errorf("TryBuiltin(%v) handled = %v, want %v", tt.fields, handled, tt.wantHandled)
			}
			if shouldExit != tt.wantExit {
				t.Errorf("TryBuiltin(%v) shouldExit = %v, want %v", tt.fields, shouldExit, tt.wantExit)
			}
		})
	}
}
