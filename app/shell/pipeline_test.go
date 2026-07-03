package shell

import (
	"bytes"
	"io"
	"runtime"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestExecutePipeline_builtinExternal(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("cat is not available on Windows")
	}

	s := New()
	var out bytes.Buffer
	executed, notFound, err := s.ExecutePipeline(
		[][]string{{"echo", "apple-orange"}, {"cat"}},
		&out,
		io.Discard,
	)
	if err != nil {
		t.Fatalf("ExecutePipeline() error = %v", err)
	}
	if !executed {
		t.Fatalf("ExecutePipeline() executed = false, notFound = %q", notFound)
	}
	if diff := cmp.Diff("apple-orange\n", out.String()); diff != "" {
		t.Errorf("output mismatch (-want +got):\n%s", diff)
	}
}

func TestExecutePipeline_externalBuiltin(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("echo is not available on Windows")
	}

	s := New()
	var out bytes.Buffer
	executed, notFound, err := s.ExecutePipeline(
		[][]string{{"echo", "ignored"}, {"type", "exit"}},
		&out,
		io.Discard,
	)
	if err != nil {
		t.Fatalf("ExecutePipeline() error = %v", err)
	}
	if !executed {
		t.Fatalf("ExecutePipeline() executed = false, notFound = %q", notFound)
	}
	if !strings.Contains(out.String(), "exit is a shell builtin") {
		t.Errorf("output = %q, want type builtin message", out.String())
	}
	if strings.Contains(out.String(), "ignored") {
		t.Errorf("output = %q, should not include piped input", out.String())
	}
}

func TestExecutePipeline_threeStages(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("cat is not available on Windows")
	}

	s := New()
	var out bytes.Buffer
	executed, notFound, err := s.ExecutePipeline(
		[][]string{{"echo", "one"}, {"cat"}, {"cat"}},
		&out,
		io.Discard,
	)
	if err != nil {
		t.Fatalf("ExecutePipeline() error = %v", err)
	}
	if !executed {
		t.Fatalf("ExecutePipeline() executed = false, notFound = %q", notFound)
	}
	if diff := cmp.Diff("one\n", out.String()); diff != "" {
		t.Errorf("output mismatch (-want +got):\n%s", diff)
	}
}

func TestExecutePipeline_unknownCommand(t *testing.T) {
	s := New()
	executed, notFound, err := s.ExecutePipeline(
		[][]string{{"echo", "hello"}, {"missing_command_xyz"}},
		io.Discard,
		io.Discard,
	)
	if err != nil {
		t.Fatalf("ExecutePipeline() error = %v", err)
	}
	if executed {
		t.Fatal("ExecutePipeline() executed = true, want false")
	}
	if notFound != "missing_command_xyz" {
		t.Errorf("notFound = %q, want missing_command_xyz", notFound)
	}
}
