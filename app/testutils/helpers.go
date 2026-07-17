package testutils

import (
	"fmt"
	"os"
	"os/exec"
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
func CreatePath(t *testing.T, root, rel string) {
	t.Helper()

	full := filepath.Join(root, filepath.FromSlash(rel))
	if strings.HasSuffix(rel, "/") {
		if err := os.MkdirAll(strings.TrimSuffix(full, string(os.PathSeparator)), 0o755); err != nil {
			t.Fatalf("CreatePath(%q) error = %v", rel, err)
		}
		return
	}

	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		t.Fatalf("CreatePath(%q) error = %v", rel, err)
	}
	if err := os.WriteFile(full, nil, 0o644); err != nil {
		t.Fatalf("CreatePath(%q) error = %v", rel, err)
	}
}

// WriteFileIn writes content to dir/name and returns the full path.
func WriteFileIn(t *testing.T, dir, name, content string) string {
	t.Helper()

	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile(%q) error = %v", path, err)
	}
	return path
}

// WriteTempFile writes content to name in a new temp directory and returns the full path.
func WriteTempFile(t *testing.T, name, content string) string {
	t.Helper()
	return WriteFileIn(t, t.TempDir(), name, content)
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
	PrependPATH(t, dir)
	return executable
}

// OutputLines splits trimmed text into lines without a trailing newline.
func OutputLines(text string) []string {
	text = strings.TrimSuffix(text, "\n")
	if text == "" {
		return nil
	}
	return strings.Split(text, "\n")
}

// PrependPATH prepends dir to PATH for the duration of the test.
func PrependPATH(t *testing.T, dir string) {
	t.Helper()
	path := os.Getenv("PATH")
	t.Setenv("PATH", dir+string(os.PathListSeparator)+path)
}

// WriteMockProgram builds a small Go mock program in dir and returns its name and path.
func WriteMockProgram(t *testing.T, dir string) (name string, path string) {
	t.Helper()
	return writeMockProgramSource(t, dir, "mock_prog", `package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		fmt.Printf("hello %s\n", os.Args[1])
		return
	}
	fmt.Println("hello")
}
`)
}

// WriteMockProgramOnPath builds WriteMockProgram and prepends its directory to PATH.
func WriteMockProgramOnPath(t *testing.T) (name string, path string) {
	t.Helper()

	dir := t.TempDir()
	name, path = WriteMockProgram(t, dir)
	PrependPATH(t, dir)
	return name, path
}

// WriteMockExitProgram builds a mock program that exits with the given code.
func WriteMockExitProgram(t *testing.T, dir string, exitCode int) (name string, path string) {
	t.Helper()
	return writeMockProgramSource(t, dir, "mock_exit", fmt.Sprintf(`package main

import "os"

func main() {
	os.Exit(%d)
}
`, exitCode))
}

// WriteMockStderrProgram builds a mock program that writes message to stderr.
func WriteMockStderrProgram(t *testing.T, dir string, message string) (name string, path string) {
	t.Helper()
	return writeMockProgramSource(t, dir, "mock_stderr", fmt.Sprintf(`package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Fprint(os.Stderr, %q)
}
`, message))
}

// WriteCompleterScript builds a programmable completer that prints the given lines.
func WriteCompleterScript(t *testing.T, dir string, lines ...string) string {
	t.Helper()

	var b strings.Builder
	b.WriteString("package main\n\nimport \"fmt\"\n\nfunc main() {\n")
	for _, line := range lines {
		fmt.Fprintf(&b, "\tfmt.Println(%q)\n", line)
	}
	b.WriteString("}\n")

	_, path := writeMockProgramSource(t, dir, "mock_completer", b.String())
	return path
}

func writeMockProgramSource(t *testing.T, dir, baseName, source string) (name string, path string) {
	t.Helper()

	src := filepath.Join(dir, baseName+".go")
	if err := os.WriteFile(src, []byte(source), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	name = baseName
	if runtime.GOOS == "windows" {
		name = baseName + ".exe"
	}
	path = filepath.Join(dir, name)
	build := exec.Command("go", "build", "-o", path, src)
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("go build error = %v\n%s", err, out)
	}
	return strings.TrimSuffix(name, ".exe"), path
}
