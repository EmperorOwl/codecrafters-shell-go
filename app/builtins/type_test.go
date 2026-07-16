package builtins

import (
	"bytes"
	"testing"

	"github.com/codecrafters-io/shell-starter-go/app/utils"
	"github.com/google/go-cmp/cmp"
)

func TestType(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(t *testing.T) string
		args    []string
		wantOut string
	}{
		{name: "echo builtin", args: []string{"echo"}, wantOut: "echo is a shell builtin\n"},
		{name: "exit builtin", args: []string{"exit"}, wantOut: "exit is a shell builtin\n"},
		{name: "type builtin", args: []string{"type"}, wantOut: "type is a shell builtin\n"},
		{name: "pwd builtin", args: []string{"pwd"}, wantOut: "pwd is a shell builtin\n"},
		{name: "cd builtin", args: []string{"cd"}, wantOut: "cd is a shell builtin\n"},
		{name: "complete builtin", args: []string{"complete"}, wantOut: "complete is a shell builtin\n"},
		{name: "jobs builtin", args: []string{"jobs"}, wantOut: "jobs is a shell builtin\n"},
		{name: "history builtin", args: []string{"history"}, wantOut: "history is a shell builtin\n"},
		{name: "declare builtin", args: []string{"declare"}, wantOut: "declare is a shell builtin\n"},
		{name: "invalid command", args: []string{"invalid_command"}, wantOut: "invalid_command: not found\n"},
		{
			name: "reports executable",
			setup: func(t *testing.T) string {
				executable := utils.CreateTempExecutable(t, "mycommand")
				return "mycommand is " + executable + "\n"
			},
			args: []string{"mycommand"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wantOut := tt.wantOut
			if tt.setup != nil {
				wantOut = tt.setup(t)
			}

			var stdout bytes.Buffer
			typeBuiltin(&stdout, tt.args[0])

			if diff := cmp.Diff(wantOut, stdout.String()); diff != "" {
				t.Errorf("typeBuiltin(%v) stdout mismatch (-want +got):\n%s", tt.args, diff)
			}
		})
	}
}
