package external

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func writeMockProgram(t *testing.T, dir string) (name string, path string) {
	t.Helper()

	src := filepath.Join(dir, "mock_prog.go")
	source := `package main

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
`
	if err := os.WriteFile(src, []byte(source), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	name = "mock_prog"
	if runtime.GOOS == "windows" {
		name = "mock_prog.exe"
	}
	path = filepath.Join(dir, name)
	build := exec.Command("go", "build", "-o", path, src)
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("go build error = %v\n%s", err, out)
	}
	return strings.TrimSuffix(name, ".exe"), path
}

func TestExternalProgram_Run(t *testing.T) {
	dir := t.TempDir()
	name, programPath := writeMockProgram(t, dir)

	tests := []struct {
		name    string
		args    []string
		wantOut string
	}{
		{
			name:    "prints greeting with argument",
			args:    []string{"world"},
			wantOut: "hello world\n",
		},
		{
			name:    "prints greeting without argument",
			args:    nil,
			wantOut: "hello\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out bytes.Buffer
			prog := &ExternalProgram{
				Name:   name,
				Path:   programPath,
				Args:   tt.args,
				Stdout: &out,
				Stderr: io.Discard,
			}

			if err := prog.Run(); err != nil {
				t.Fatalf("Run() error = %v", err)
			}
			if diff := cmp.Diff(tt.wantOut, out.String()); diff != "" {
				t.Errorf("stdout mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestExternalProgram_RunInBackground(t *testing.T) {
	dir := t.TempDir()
	name, programPath := writeMockProgram(t, dir)

	exited := make(chan struct{})
	prog := &ExternalProgram{
		Name:   name,
		Path:   programPath,
		Stdout: io.Discard,
		Stderr: io.Discard,
	}

	pid, err := prog.RunInBackground(func() { close(exited) })
	if err != nil {
		t.Fatalf("RunInBackground() error = %v", err)
	}
	if pid <= 0 {
		t.Fatalf("RunInBackground() pid = %d, want > 0", pid)
	}

	select {
	case <-exited:
	case <-time.After(5 * time.Second):
		t.Fatal("onExit was not called within 5s")
	}
}
