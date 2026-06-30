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
		name      string
		command   string
		isBuiltin bool
		want      string
	}{
		{name: "echo builtin", command: "echo", isBuiltin: true, want: "echo is a shell builtin\n"},
		{name: "exit builtin", command: "exit", isBuiltin: true, want: "exit is a shell builtin\n"},
		{name: "type builtin", command: "type", isBuiltin: true, want: "type is a shell builtin\n"},
		{name: "pwd builtin", command: "pwd", isBuiltin: true, want: "pwd is a shell builtin\n"},
		{name: "cd builtin", command: "cd", isBuiltin: true, want: "cd is a shell builtin\n"},
		{name: "complete builtin", command: "complete", isBuiltin: true, want: "complete is a shell builtin\n"},
		{name: "invalid command", command: "invalid_command", isBuiltin: false, want: "invalid_command: not found\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out bytes.Buffer
			Type(&out, tt.command, tt.isBuiltin)
			if diff := cmp.Diff(tt.want, out.String()); diff != "" {
				t.Errorf("Type() output mismatch (-want +got):\n%s", diff)
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

		var out bytes.Buffer
		Type(&out, command, false)
		want := command + " is " + executable + "\n"
		if diff := cmp.Diff(want, out.String()); diff != "" {
			t.Errorf("Type() output mismatch (-want +got):\n%s", diff)
		}
	})
}
