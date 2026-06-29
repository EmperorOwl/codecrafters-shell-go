package completion

import (
	"testing"

	"github.com/codecrafters-io/shell-starter-go/app/builtins"
)

func TestApplyTab(t *testing.T) {
	builtinsList := []string{"cd", "echo", "exit", "pwd", "type"}
	listFiles := func(dir string) []string {
		if dir == "" {
			return []string{"first.txt", "readme.txt", "second.txt"}
		}
		return nil
	}

	tests := []struct {
		name                 string
		registeredCompleters map[string]builtins.Completer
		buffer               string
		wantBuffer           string
	}{
		{
			name:       "command completion",
			buffer:     "ech",
			wantBuffer: "echo ",
		},
		{
			name:       "filename completion",
			buffer:     "cat re",
			wantBuffer: "cat readme.txt ",
		},
		{
			name:       "later argument filename completion",
			buffer:     "echo first.txt sec",
			wantBuffer: "echo first.txt second.txt ",
		},
		{
			name: "programmable completion",
			registeredCompleters: map[string]builtins.Completer{
				"docker": {
					Path: "/path/to/completer",
					Func: func(builtins.CompleterFuncOptions) ([]string, error) {
						return []string{"run"}, nil
					},
				},
			},
			buffer:     "docker ",
			wantBuffer: "docker run ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBuffer, gotListings := ApplyTab(builtinsList, nil, listFiles, tt.registeredCompleters, tt.buffer)
			if gotBuffer != tt.wantBuffer {
				t.Errorf("ApplyTab(%q) buffer = %q, want %q", tt.buffer, gotBuffer, tt.wantBuffer)
			}
			if len(gotListings) != 0 {
				t.Errorf("ApplyTab(%q) listings = %v, want nil", tt.buffer, gotListings)
			}
		})
	}
}
