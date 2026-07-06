package builtins

import (
	"bytes"
	"testing"

	"github.com/codecrafters-io/shell-starter-go/app/completion"
	"github.com/codecrafters-io/shell-starter-go/app/jobs"
	"github.com/google/go-cmp/cmp"
)

func TestRun(t *testing.T) {
	tests := []struct {
		name          string
		builtinName   string
		args          []string
		ctx           *Context
		wantExitShell bool
	}{
		{
			name:          "exit terminates shell",
			builtinName:   "exit",
			ctx:           &Context{},
			wantExitShell: true,
		},
		{
			name:        "echo runs",
			builtinName: "echo",
			args:        []string{"hello", "world"},
			ctx:         &Context{},
		},
		{
			name:        "pwd runs",
			builtinName: "pwd",
			ctx:         &Context{},
		},
		{
			name:        "cd runs",
			builtinName: "cd",
			args:        []string{"/tmp"},
			ctx:         &Context{},
		},
		{
			name:        "type runs",
			builtinName: "type",
			args:        []string{"echo"},
			ctx:         &Context{},
		},
		{
			name:        "complete runs",
			builtinName: "complete",
			args:        []string{"-C", "/path/to/script", "git"},
			ctx: &Context{
				State: &State{Completion: completion.NewCompletionRegistry()},
			},
		},
		{
			name:        "complete -p runs",
			builtinName: "complete",
			args:        []string{"-p", "git"},
			ctx: &Context{
				State: &State{Completion: completion.NewCompletionRegistry()},
			},
		},
		{
			name:        "jobs runs",
			builtinName: "jobs",
			ctx: &Context{
				State: &State{Jobs: &jobs.JobTable{}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out bytes.Buffer
			var errOut bytes.Buffer
			tt.ctx.Stdout = &out
			tt.ctx.Stderr = &errOut

			exitShell, err := Run(tt.builtinName, tt.args, tt.ctx)
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
