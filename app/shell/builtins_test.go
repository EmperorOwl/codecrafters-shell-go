package shell

import (
	"bytes"
	"testing"

	"github.com/codecrafters-io/shell-starter-go/app/jobs"
	"github.com/google/go-cmp/cmp"
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
			name:        "jobs is handled",
			fields:      []string{"jobs"},
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
			var jobList []jobs.Job
			handled, shouldExit := TryBuiltin(tt.fields, &out, &errOut, map[string]string{}, &jobList)
			if diff := cmp.Diff(tt.wantHandled, handled); diff != "" {
				t.Errorf("TryBuiltin(%v) handled mismatch (-want +got):\n%s", tt.fields, diff)
			}
			if diff := cmp.Diff(tt.wantExit, shouldExit); diff != "" {
				t.Errorf("TryBuiltin(%v) shouldExit mismatch (-want +got):\n%s", tt.fields, diff)
			}
		})
	}
}
