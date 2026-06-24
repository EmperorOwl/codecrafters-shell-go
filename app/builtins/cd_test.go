package builtins

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestCd(t *testing.T) {
	t.Run("changes to absolute path", func(t *testing.T) {
		target, err := filepath.Abs(t.TempDir())
		if err != nil {
			t.Fatalf("Abs() error = %v", err)
		}

		t.Chdir(t.TempDir())

		var out bytes.Buffer
		Cd(&out, target)
		if got := out.String(); got != "" {
			t.Errorf("Cd() output = %q, want empty", got)
		}

		cwd, err := os.Getwd()
		if err != nil {
			t.Fatalf("Getwd() error = %v", err)
		}
		if cwd != target {
			t.Errorf("Getwd() = %q, want %q", cwd, target)
		}
	})

	t.Run("prints error for missing directory", func(t *testing.T) {
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
		Cd(&out, invalid)
		want := CdErrorMessage(invalid) + "\n"
		if got := out.String(); got != want {
			t.Errorf("Cd() output = %q, want %q", got, want)
		}

		afterCwd, err := os.Getwd()
		if err != nil {
			t.Fatalf("Getwd() error = %v", err)
		}
		if afterCwd != cwd {
			t.Errorf("Getwd() after failed cd = %q, want %q", afterCwd, cwd)
		}
	})
}
