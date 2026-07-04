package builtins

import (
	"bytes"
	"testing"

	"github.com/codecrafters-io/shell-starter-go/app/jobs"
	"github.com/google/go-cmp/cmp"
)

func TestBuiltin_Run(t *testing.T) {
	tests := []struct {
		name          string
		builtin       *Builtin
		wantExitShell bool
	}{
		{
			name:          "exit terminates shell",
			builtin:       New("exit", nil, nil, nil, nil, nil),
			wantExitShell: true,
		},
		{
			name:    "echo runs",
			builtin: New("echo", []string{"hello", "world"}, nil, nil, nil, nil),
		},
		{
			name:    "pwd runs",
			builtin: New("pwd", nil, nil, nil, nil, nil),
		},
		{
			name:    "cd runs",
			builtin: New("cd", []string{"/tmp"}, nil, nil, nil, nil),
		},
		{
			name:    "type runs",
			builtin: New("type", []string{"echo"}, nil, nil, nil, nil),
		},
		{
			name:    "complete runs",
			builtin: New("complete", []string{"-C", "/path/to/script", "git"}, nil, nil, map[string]string{}, nil),
		},
		{
			name:    "complete -p runs",
			builtin: New("complete", []string{"-p", "git"}, nil, nil, map[string]string{}, nil),
		},
		{
			name:    "jobs runs",
			builtin: New("jobs", nil, nil, nil, nil, &jobs.JobTable{}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out bytes.Buffer
			var errOut bytes.Buffer
			tt.builtin.Stdout = &out
			tt.builtin.Stderr = &errOut

			exitShell, err := tt.builtin.Run()
			if err != nil {
				t.Fatalf("Run() error = %v", err)
			}
			if diff := cmp.Diff(tt.wantExitShell, exitShell); diff != "" {
				t.Errorf("Run() exitShell mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestIsBuiltin(t *testing.T) {
	if diff := cmp.Diff(true, IsBuiltin("echo")); diff != "" {
		t.Errorf("IsBuiltin(echo) mismatch (-want +got):\n%s", diff)
	}
	if diff := cmp.Diff(false, IsBuiltin("xyz")); diff != "" {
		t.Errorf("IsBuiltin(xyz) mismatch (-want +got):\n%s", diff)
	}
}
