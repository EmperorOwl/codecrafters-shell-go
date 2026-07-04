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
	_, err := s.executePipeline(
		[][]string{{"echo", "apple-orange"}, {"cat"}},
		lineContext{stdout: &out, stderr: io.Discard},
	)
	if err != nil {
		t.Fatalf("executePipeline() error = %v", err)
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
	_, err := s.executePipeline(
		[][]string{{"echo", "ignored"}, {"type", "exit"}},
		lineContext{stdout: &out, stderr: io.Discard},
	)
	if err != nil {
		t.Fatalf("executePipeline() error = %v", err)
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
	_, err := s.executePipeline(
		[][]string{{"echo", "one"}, {"cat"}, {"cat"}},
		lineContext{stdout: &out, stderr: io.Discard},
	)
	if err != nil {
		t.Fatalf("executePipeline() error = %v", err)
	}
	if diff := cmp.Diff("one\n", out.String()); diff != "" {
		t.Errorf("output mismatch (-want +got):\n%s", diff)
	}
}

func TestExecutePipeline_unknownCommand(t *testing.T) {
	s := New()
	var out bytes.Buffer
	_, err := s.executePipeline(
		[][]string{{"echo", "hello"}, {"missing_command_xyz"}},
		lineContext{stdout: &out, stderr: io.Discard},
	)
	if err != nil {
		t.Fatalf("executePipeline() error = %v", err)
	}
	if !strings.Contains(out.String(), "missing_command_xyz: command not found") {
		t.Errorf("output = %q, want command not found message", out.String())
	}
}
