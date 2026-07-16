package utils

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// WantLines joins lines with newlines and a trailing newline.
func WantLines(lines []string) string {
	if len(lines) == 0 {
		return ""
	}
	return strings.Join(lines, "\n") + "\n"
}

// WantStdout joins expected output lines into the trailing-newline format
// builtins write to stdout.
func WantStdout(lines []string) string {
	return WantLines(lines)
}

// CreatePath creates root/rel for tests. A trailing slash on rel creates a directory;
// otherwise it creates parent directories and an empty file.
func CreatePath(root, rel string) error {
	full := filepath.Join(root, filepath.FromSlash(rel))
	if strings.HasSuffix(rel, "/") {
		return os.MkdirAll(strings.TrimSuffix(full, string(os.PathSeparator)), 0o755)
	}

	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		return err
	}
	return os.WriteFile(full, nil, 0o644)
}

// WriteTempFile writes content to name in a new temp directory and returns the full path.
func WriteTempFile(t *testing.T, name, content string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile(%q) error = %v", path, err)
	}
	return path
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
