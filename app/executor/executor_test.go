package executor

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/codecrafters-io/shell-starter-go/app/parser"
	"github.com/codecrafters-io/shell-starter-go/app/session"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestExecuteBuiltin(t *testing.T) {
	var out bytes.Buffer
	e := New(strings.NewReader(""))
	state := session.NewState()

	exitShell, err := e.ExecuteBuiltin(Outputs{
		Stdout: &out,
		Stderr: io.Discard,
	}, state, []string{"echo", "hello", "world"})
	if err != nil {
		t.Fatalf("ExecuteBuiltin() error = %v", err)
	}
	if exitShell {
		t.Fatal("ExecuteBuiltin() exitShell = true, want false")
	}
	if diff := cmp.Diff("hello world\n", out.String()); diff != "" {
		t.Errorf("stdout mismatch (-want +got):\n%s", diff)
	}
}

func TestExecuteBuiltinExit(t *testing.T) {
	e := New(strings.NewReader(""))
	state := session.NewState()

	exitShell, err := e.ExecuteBuiltin(Outputs{
		Stdout: io.Discard,
		Stderr: io.Discard,
	}, state, []string{"exit"})
	if err != nil {
		t.Fatalf("ExecuteBuiltin() error = %v", err)
	}
	if !exitShell {
		t.Error("ExecuteBuiltin(exit) exitShell = false, want true")
	}
}

func TestExecuteExternalForeground(t *testing.T) {
	var out bytes.Buffer
	e := New(strings.NewReader(""))

	err := e.ExecuteExternalForeground(Outputs{
		Stdout: &out,
		Stderr: io.Discard,
	}, []string{"echo", "external"})
	if err != nil {
		t.Fatalf("ExecuteExternalForeground() error = %v", err)
	}
	if diff := cmp.Diff("external\n", out.String()); diff != "" {
		t.Errorf("stdout mismatch (-want +got):\n%s", diff)
	}
}

func TestExecutePipeline(t *testing.T) {
	var out bytes.Buffer
	e := New(strings.NewReader(""))
	state := session.NewState()

	err := e.ExecutePipeline(Outputs{
		Stdout: &out,
		Stderr: io.Discard,
	}, state, [][]string{
		{"echo", "pipe"},
		{"echo", "line"},
	})
	if err != nil {
		t.Fatalf("ExecutePipeline() error = %v", err)
	}
	if diff := cmp.Diff("line\n", out.String()); diff != "" {
		t.Errorf("stdout mismatch (-want +got):\n%s", diff)
	}
}

func TestExecutePipelineBuiltinMiddle(t *testing.T) {
	var out bytes.Buffer
	e := New(strings.NewReader(""))
	state := session.NewState()

	err := e.ExecutePipeline(Outputs{
		Stdout: &out,
		Stderr: io.Discard,
	}, state, [][]string{
		{"echo", "hello"},
		{"type", "echo"},
	})
	if err != nil {
		t.Fatalf("ExecutePipeline() error = %v", err)
	}
	want := "echo is a shell builtin\n"
	if diff := cmp.Diff(want, out.String()); diff != "" {
		t.Errorf("stdout mismatch (-want +got):\n%s", diff)
	}
}

func TestStdoutRedirect(t *testing.T) {
	dir := t.TempDir()
	outPath := filepath.Join(dir, "out.txt")

	e := New(strings.NewReader(""))
	state := session.NewState()

	_, err := e.ExecuteBuiltin(Outputs{
		Stdout: io.Discard,
		Stderr: io.Discard,
		Redirect: parser.Redirect{StdoutPath: outPath},
	}, state, []string{"echo", "redirected"})
	if err != nil {
		t.Fatalf("ExecuteBuiltin() error = %v", err)
	}

	got, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if diff := cmp.Diff("redirected\n", string(got)); diff != "" {
		t.Errorf("redirect file mismatch (-want +got):\n%s", diff)
	}
}

func TestExecutePipelineTooFewSegments(t *testing.T) {
	e := New(strings.NewReader(""))
	state := session.NewState()

	err := e.ExecutePipeline(Outputs{
		Stdout: io.Discard,
		Stderr: io.Discard,
	}, state, [][]string{{"echo"}})
	if err != nil {
		t.Fatalf("ExecutePipeline() error = %v", err)
	}
}

func TestNonExitError(t *testing.T) {
	if diff := cmp.Diff(nil, nonExitError(nil)); diff != "" {
		t.Errorf("nonExitError(nil) mismatch (-want +got):\n%s", diff)
	}
}

func outputLines(text string) []string {
	text = strings.TrimSuffix(text, "\n")
	if text == "" {
		return nil
	}
	return strings.Split(text, "\n")
}

func TestExecuteBuiltinStdoutLines(t *testing.T) {
	var out bytes.Buffer
	e := New(strings.NewReader(""))
	state := session.NewState()

	_, err := e.ExecuteBuiltin(Outputs{Stdout: &out, Stderr: io.Discard}, state, []string{"echo", "a", "b"})
	if err != nil {
		t.Fatalf("ExecuteBuiltin() error = %v", err)
	}
	if diff := cmp.Diff([]string{"a b"}, outputLines(out.String()), cmpopts.EquateEmpty()); diff != "" {
		t.Errorf("stdout lines mismatch (-want +got):\n%s", diff)
	}
}
