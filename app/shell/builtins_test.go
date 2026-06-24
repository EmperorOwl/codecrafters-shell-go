package shell

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestTryBuiltin(t *testing.T) {
	tests := []struct {
		name        string
		line        string
		wantOutput  string
		wantHandled bool
		wantExit    bool
	}{
		{
			name:        "exit terminates shell",
			line:        "exit",
			wantHandled: true,
			wantExit:    true,
		},
		{
			name:        "echo prints arguments",
			line:        "echo hello world",
			wantOutput:  "hello world\n",
			wantHandled: true,
		},
		{
			name:        "echo three words",
			line:        "echo one two three",
			wantOutput:  "one two three\n",
			wantHandled: true,
		},
		{
			name:        "echo no args",
			line:        "echo",
			wantOutput:  "\n",
			wantHandled: true,
		},
		{
			name:        "type reports echo builtin",
			line:        "type echo",
			wantOutput:  "echo is a shell builtin\n",
			wantHandled: true,
		},
		{
			name:        "type reports exit builtin",
			line:        "type exit",
			wantOutput:  "exit is a shell builtin\n",
			wantHandled: true,
		},
		{
			name:        "type reports type builtin",
			line:        "type type",
			wantOutput:  "type is a shell builtin\n",
			wantHandled: true,
		},
		{
			name:        "type reports pwd builtin",
			line:        "type pwd",
			wantOutput:  "pwd is a shell builtin\n",
			wantHandled: true,
		},
		{
			name:        "type reports cd builtin",
			line:        "type cd",
			wantOutput:  "cd is a shell builtin\n",
			wantHandled: true,
		},
		{
			name:        "type reports not found",
			line:        "type invalid_command",
			wantOutput:  "invalid_command: not found\n",
			wantHandled: true,
		},
		{
			name:        "unknown command is not builtin",
			line:        "xyz",
			wantHandled: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out bytes.Buffer
			handled, shouldExit := TryBuiltin(tt.line, &out)
			if handled != tt.wantHandled {
				t.Errorf("TryBuiltin(%q) handled = %v, want %v", tt.line, handled, tt.wantHandled)
			}
			if shouldExit != tt.wantExit {
				t.Errorf("TryBuiltin(%q) shouldExit = %v, want %v", tt.line, shouldExit, tt.wantExit)
			}
			if got := out.String(); got != tt.wantOutput {
				t.Errorf("TryBuiltin(%q) output = %q, want %q", tt.line, got, tt.wantOutput)
			}
		})
	}

	t.Run("pwd prints working directory", func(t *testing.T) {
		cwd, err := os.Getwd()
		if err != nil {
			t.Fatalf("Getwd() error = %v", err)
		}

		var out bytes.Buffer
		handled, shouldExit := TryBuiltin("pwd", &out)
		if !handled {
			t.Fatalf("TryBuiltin() handled = false, want true")
		}
		if shouldExit {
			t.Errorf("TryBuiltin() shouldExit = true, want false")
		}
		wantOutput := cwd + "\n"
		if got := out.String(); got != wantOutput {
			t.Errorf("TryBuiltin() output = %q, want %q", got, wantOutput)
		}
	})

	t.Run("cd changes to absolute path", func(t *testing.T) {
		target, err := filepath.Abs(t.TempDir())
		if err != nil {
			t.Fatalf("Abs() error = %v", err)
		}

		t.Chdir(t.TempDir())

		var out bytes.Buffer
		handled, shouldExit := TryBuiltin("cd "+target, &out)
		if !handled {
			t.Fatalf("TryBuiltin() handled = false, want true")
		}
		if shouldExit {
			t.Errorf("TryBuiltin() shouldExit = true, want false")
		}
		if got := out.String(); got != "" {
			t.Errorf("TryBuiltin() output = %q, want empty", got)
		}

		cwd, err := os.Getwd()
		if err != nil {
			t.Fatalf("Getwd() error = %v", err)
		}
		if cwd != target {
			t.Errorf("Getwd() = %q, want %q", cwd, target)
		}
	})

	t.Run("cd prints error for missing directory", func(t *testing.T) {
		invalid := "/does_not_exist_codecrafters_test"
		if runtime.GOOS == "windows" {
			vol := os.Getenv("SystemDrive")
			if vol == "" {
				vol = "C:"
			}
			invalid = filepath.Join(vol+string(filepath.Separator), "does_not_exist_codecrafters_test")
		}

		cwd, err := os.Getwd()
		if err != nil {
			t.Fatalf("Getwd() error = %v", err)
		}

		var out bytes.Buffer
		handled, shouldExit := TryBuiltin("cd "+invalid, &out)
		if !handled {
			t.Fatalf("TryBuiltin() handled = false, want true")
		}
		if shouldExit {
			t.Errorf("TryBuiltin() shouldExit = true, want false")
		}
		wantOutput := CdErrorMessage(invalid) + "\n"
		if got := out.String(); got != wantOutput {
			t.Errorf("TryBuiltin() output = %q, want %q", got, wantOutput)
		}

		afterCwd, err := os.Getwd()
		if err != nil {
			t.Fatalf("Getwd() error = %v", err)
		}
		if afterCwd != cwd {
			t.Errorf("Getwd() after failed cd = %q, want %q", afterCwd, cwd)
		}
	})

	t.Run("type reports executable", func(t *testing.T) {
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
		handled, shouldExit := TryBuiltin("type "+command, &out)
		if !handled {
			t.Fatalf("TryBuiltin() handled = false, want true")
		}
		if shouldExit {
			t.Errorf("TryBuiltin() shouldExit = true, want false")
		}
		wantOutput := command + " is " + executable + "\n"
		if got := out.String(); got != wantOutput {
			t.Errorf("TryBuiltin() output = %q, want %q", got, wantOutput)
		}
	})
}
