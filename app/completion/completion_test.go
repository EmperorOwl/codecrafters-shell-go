package completion

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
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
		name           string
		executables    []string
		completeHandler CompleteHandler
		buffer         string
		wantBuffer     string
	}{
		{
			name:       "builtin command completion",
			buffer:     "ech",
			wantBuffer: "echo ",
		},
		{
			name:        "external command completion",
			executables: []string{"custom_executable"},
			buffer:      "custom",
			wantBuffer:  "custom_executable ",
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
			completeHandler: func(CompleterFuncOptions) []string {
				return []string{"run"}
			},
			buffer:     "docker ",
			wantBuffer: "docker run ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBuffer, gotListings := ApplyTab(builtinsList, tt.executables, listFiles, tt.completeHandler, tt.buffer)
			if diff := cmp.Diff(tt.wantBuffer, gotBuffer); diff != "" {
				t.Errorf("ApplyTab(%q) buffer mismatch (-want +got):\n%s", tt.buffer, diff)
			}
			if diff := cmp.Diff([]string(nil), gotListings, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("ApplyTab(%q) listings mismatch (-want +got):\n%s", tt.buffer, diff)
			}
		})
	}
}
