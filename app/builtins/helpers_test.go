package builtins

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func wantStdout(lines []string) string {
	if len(lines) == 0 {
		return ""
	}
	return strings.Join(lines, "\n") + "\n"
}

func createTempExecutable(t *testing.T, name string) string {
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
