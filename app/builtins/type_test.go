package builtins

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestType(t *testing.T) {
	tests := []struct {
		name    string
		command string
		wantOut string
	}{
		{name: "echo builtin", command: "echo", wantOut: "echo is a shell builtin\n"},
		{name: "exit builtin", command: "exit", wantOut: "exit is a shell builtin\n"},
		{name: "type builtin", command: "type", wantOut: "type is a shell builtin\n"},
		{name: "pwd builtin", command: "pwd", wantOut: "pwd is a shell builtin\n"},
		{name: "cd builtin", command: "cd", wantOut: "cd is a shell builtin\n"},
		{name: "complete builtin", command: "complete", wantOut: "complete is a shell builtin\n"},
		{name: "jobs builtin", command: "jobs", wantOut: "jobs is a shell builtin\n"},
		{name: "history builtin", command: "history", wantOut: "history is a shell builtin\n"},
		{name: "declare builtin", command: "declare", wantOut: "declare is a shell builtin\n"},
		{name: "invalid command", command: "invalid_command", wantOut: "invalid_command: not found\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout bytes.Buffer
			typeBuiltin(&stdout, tt.command)
			if diff := cmp.Diff(tt.wantOut, stdout.String()); diff != "" {
				t.Errorf("typeBuiltin(%v) stdout mismatch (-want +got):\n%s", tt.command, diff)
			}
		})
	}

	t.Run("reports executable", func(t *testing.T) {
		dir := t.TempDir()
		command := "mycommand"
		fileName := command
		if runtime.GOOS == "windows" {
			fileName += ".exe"
		}
		executable := filepath.Join(dir, fileName)
		perms := os.FileMode(0o755)
		if runtime.GOOS == "windows" {
			perms = 0o644
		}
		if err := os.WriteFile(executable, nil, perms); err != nil {
			t.Fatalf("WriteFile() error = %v", err)
		}
		t.Setenv("PATH", dir)

		var stdout bytes.Buffer
		typeBuiltin(&stdout, command)
		want := command + " is " + executable + "\n"
		if diff := cmp.Diff(want, stdout.String()); diff != "" {
			t.Errorf("typeBuiltin(%q) stdout mismatch (-want +got):\n%s", command, diff)
		}
	})
}
