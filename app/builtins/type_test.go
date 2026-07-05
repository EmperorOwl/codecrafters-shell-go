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
		want    string
	}{
		{name: "echo builtin", command: "echo", want: "echo is a shell builtin\n"},
		{name: "exit builtin", command: "exit", want: "exit is a shell builtin\n"},
		{name: "type builtin", command: "type", want: "type is a shell builtin\n"},
		{name: "pwd builtin", command: "pwd", want: "pwd is a shell builtin\n"},
		{name: "cd builtin", command: "cd", want: "cd is a shell builtin\n"},
		{name: "complete builtin", command: "complete", want: "complete is a shell builtin\n"},
		{name: "jobs builtin", command: "jobs", want: "jobs is a shell builtin\n"},
		{name: "invalid command", command: "invalid_command", want: "invalid_command: not found\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out bytes.Buffer
			Type(&out, tt.command)
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
		Type(&out, command)
		want := command + " is " + executable + "\n"
		if diff := cmp.Diff(want, out.String()); diff != "" {
			t.Errorf("Type() output mismatch (-want +got):\n%s", diff)
		}
	})
}
