package builtins

import (
	"bytes"
	"testing"

	"github.com/codecrafters-io/shell-starter-go/app/completion"
	"github.com/codecrafters-io/shell-starter-go/app/history"
	"github.com/codecrafters-io/shell-starter-go/app/jobs"
	"github.com/codecrafters-io/shell-starter-go/app/repl"
	"github.com/google/go-cmp/cmp"
)

func noopBuiltin(ctx *Context, args []string) (bool, error) {
	return false, nil
}

func TestBuiltinRegistry_Register(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*BuiltinRegistry)
		command string
		wantOK  bool
	}{
		{
			name:    "registers command",
			setup:   func(*BuiltinRegistry) {},
			command: "echo",
			wantOK:  true,
		},
		{
			name: "idempotent register",
			setup: func(r *BuiltinRegistry) {
				r.Register("echo", noopBuiltin)
			},
			command: "echo",
			wantOK:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reg := NewBuiltinRegistry()
			tt.setup(reg)
			reg.Register(tt.command, noopBuiltin)

			gotOK := reg.Is(tt.command)
			if diff := cmp.Diff(tt.wantOK, gotOK); diff != "" {
				t.Errorf("Is(%q) mismatch (-want +got):\n%s", tt.command, diff)
			}
		})
	}
}

func TestBuiltinRegistry_Names(t *testing.T) {
	reg := NewBuiltinRegistry()
	reg.Register("zebra", noopBuiltin)
	reg.Register("alpha", noopBuiltin)

	got := reg.Names()
	want := []string{"alpha", "zebra"}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Names() mismatch (-want +got):\n%s", diff)
	}
}

func TestBuiltinRegistry_Run(t *testing.T) {
	reg := NewBuiltinRegistry()
	reg.Register("exit", exitBuiltin)

	exitShell, err := reg.Run("exit", nil, &Context{})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if !exitShell {
		t.Error("Run(exit) exitShell = false, want true")
	}

	exitShell, err = reg.Run("missing", nil, &Context{})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if exitShell {
		t.Error("Run(missing) exitShell = true, want false")
	}
}

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
				State: &repl.State{Completion: completion.NewCompletionRegistry()},
			},
		},
		{
			name:        "complete -p runs",
			builtinName: "complete",
			args:        []string{"-p", "git"},
			ctx: &Context{
				State: &repl.State{Completion: completion.NewCompletionRegistry()},
			},
		},
		{
			name:        "jobs runs",
			builtinName: "jobs",
			ctx: &Context{
				State: &repl.State{Jobs: &jobs.JobTable{}},
			},
		},
		{
			name:        "history runs",
			builtinName: "history",
			ctx: &Context{
				State: &repl.State{History: &history.HistoryList{}},
			},
		},
		{
			name:        "declare runs",
			builtinName: "declare",
			ctx:         &Context{},
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

func TestNames(t *testing.T) {
	want := []string{"cd", "complete", "declare", "echo", "exit", "history", "jobs", "pwd", "type"}
	got := Names()
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Names() mismatch (-want +got):\n%s", diff)
	}
}
