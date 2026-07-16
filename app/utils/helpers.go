package utils

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// WantStdout joins expected output lines into the trailing-newline format
// builtins write to stdout.
func WantStdout(lines []string) string {
	if len(lines) == 0 {
		return ""
	}
	return strings.Join(lines, "\n") + "\n"
}

// CreateTempExecutable writes a dummy executable named name into a temp dir,
// prepends that dir to PATH, and returns the executable's full path.
func CreateTempExecutable(t *testing.T, name string) string {
	t.Helper()

	dir := t.TempDir()
	fileName := name
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
	return executable
}
