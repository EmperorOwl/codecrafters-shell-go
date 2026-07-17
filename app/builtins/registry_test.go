package builtins

import (
	"bytes"
	"testing"

	"github.com/codecrafters-io/shell-starter-go/app/completion"
	"github.com/codecrafters-io/shell-starter-go/app/history"
	"github.com/codecrafters-io/shell-starter-go/app/jobs"
	"github.com/codecrafters-io/shell-starter-go/app/session"
	"github.com/codecrafters-io/shell-starter-go/app/variables"
	"github.com/google/go-cmp/cmp"
)

func noopHandler(ctx *Context, args []string) (bool, error) {
	return false, nil
}

func TestRegistry_Register(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*registry)
		command string
		wantOK  bool
	}{
		{
			name:    "registers command",
			setup:   func(*registry) {},
			command: "echo",
			wantOK:  true,
		},
		{
			name: "idempotent register",
			setup: func(r *registry) {
				r.register("echo", noopHandler)
			},
			command: "echo",
			wantOK:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reg := newRegistry()
			tt.setup(reg)
			reg.register(tt.command, noopHandler)

			gotOK := reg.is(tt.command)
			if diff := cmp.Diff(tt.wantOK, gotOK); diff != "" {
				t.Errorf("Is(%q) mismatch (-want +got):\n%s", tt.command, diff)
			}
		})
	}
}

func TestRegistry_Names(t *testing.T) {
	reg := newRegistry()
	reg.register("zebra", noopHandler)
	reg.register("alpha", noopHandler)

	got := reg.names()
	want := []string{"alpha", "zebra"}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Names() mismatch (-want +got):\n%s", diff)
	}
}

func TestRegistry_Run(t *testing.T) {
	reg := newRegistry()
	reg.register("exit", exitHandler)

	exitShell, err := reg.run("exit", nil, &Context{})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if !exitShell {
		t.Error("Run(exit) exitShell = false, want true")
	}

	exitShell, err = reg.run("missing", nil, &Context{})
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
				Session: &session.Session{Completion: completion.NewRegistry()},
			},
		},
		{
			name:        "complete -p runs",
			builtinName: "complete",
			args:        []string{"-p", "git"},
			ctx: &Context{
				Session: &session.Session{Completion: completion.NewRegistry()},
			},
		},
		{
			name:        "jobs runs",
			builtinName: "jobs",
			ctx: &Context{
				Session: &session.Session{Jobs: jobs.NewTable()},
			},
		},
		{
			name:        "history runs",
			builtinName: "history",
			ctx: &Context{
				Session: &session.Session{History: history.NewList()},
			},
		},
		{
			name:        "declare runs",
			builtinName: "declare",
			ctx: &Context{
				Session: &session.Session{Variables: variables.NewStore()},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout bytes.Buffer
			var stderr bytes.Buffer
			tt.ctx.Stdout = &stdout
			tt.ctx.Stderr = &stderr

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
