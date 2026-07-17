package executor

import (
	"bytes"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/codecrafters-io/shell-starter-go/app/parser"
	"github.com/codecrafters-io/shell-starter-go/app/session"
	"github.com/codecrafters-io/shell-starter-go/app/testutils"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestExecuteBuiltin(t *testing.T) {
	tests := []struct {
		name       string
		fields     []string
		wantStdout string
		wantExit   bool
	}{
		{
			name:       "echo prints arguments",
			fields:     []string{"echo", "hello", "world"},
			wantStdout: "hello world\n",
		},
		{
			name:     "exit requests shell exit",
			fields:   []string{"exit"},
			wantExit: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out bytes.Buffer
			e := New(strings.NewReader(""), &out, io.Discard)
			state := session.NewSession()

			exitShell, err := e.ExecuteBuiltin(parser.Redirect{}, state, tt.fields)
			if err != nil {
				t.Fatalf("ExecuteBuiltin() error = %v", err)
			}
			if exitShell != tt.wantExit {
				t.Errorf("ExecuteBuiltin() exitShell = %v, want %v", exitShell, tt.wantExit)
			}
			if diff := cmp.Diff(tt.wantStdout, out.String()); diff != "" {
				t.Errorf("stdout mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestExecuteExternalForeground(t *testing.T) {
	tests := []struct {
		name       string
		fields     []string
		wantStdout string
	}{
		{
			name:       "runs external command",
			fields:     []string{"echo", "external"},
			wantStdout: "external\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out bytes.Buffer
			e := New(strings.NewReader(""), &out, io.Discard)

			err := e.ExecuteExternalForeground(parser.Redirect{}, tt.fields)
			if err != nil {
				t.Fatalf("ExecuteExternalForeground() error = %v", err)
			}
			if diff := cmp.Diff(tt.wantStdout, out.String()); diff != "" {
				t.Errorf("stdout mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestExecutePipeline(t *testing.T) {
	tests := []struct {
		name       string
		segments   [][]string
		wantStdout string
		wantErr    string
	}{
		{
			name: "external pipe passes output",
			segments: [][]string{
				{"echo", "pipe"},
				{"echo", "line"},
			},
			wantStdout: "line\n",
		},
		{
			name: "builtin in middle of pipeline",
			segments: [][]string{
				{"echo", "hello"},
				{"type", "echo"},
			},
			wantStdout: "echo is a shell builtin\n",
		},
		{
			name:     "single segment is no-op",
			segments: [][]string{{"echo"}},
		},
		{
			name: "empty segment returns error",
			segments: [][]string{
				{"echo", "a"},
				{},
				{"echo", "b"},
			},
			wantErr: "empty pipeline segment",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out bytes.Buffer
			e := New(strings.NewReader(""), &out, io.Discard)
			state := session.NewSession()

			err := e.ExecutePipeline(parser.Redirect{}, state, tt.segments)
			if tt.wantErr != "" {
				if err == nil {
					t.Fatal("ExecutePipeline() error = nil, want error")
				}
				if diff := cmp.Diff(tt.wantErr, err.Error()); diff != "" {
					t.Errorf("error mismatch (-want +got):\n%s", diff)
				}
				return
			}
			if err != nil {
				t.Fatalf("ExecutePipeline() error = %v", err)
			}
			if diff := cmp.Diff(tt.wantStdout, out.String()); diff != "" {
				t.Errorf("stdout mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestRedirects(t *testing.T) {
	tests := []struct {
		name         string
		setup        func(t *testing.T, dir string) parser.Redirect
		fields       []string
		useExternal  bool
		wantFile     string
		wantContents string
	}{
		{
			name: "stdout redirect",
			setup: func(t *testing.T, dir string) parser.Redirect {
				return parser.Redirect{StdoutPath: filepath.Join(dir, "out.txt")}
			},
			fields:       []string{"echo", "redirected"},
			wantFile:     "out.txt",
			wantContents: "redirected\n",
		},
		{
			name: "stdout append redirect",
			setup: func(t *testing.T, dir string) parser.Redirect {
				path := testutils.WriteFileIn(t, dir, "append.txt", "first\n")
				return parser.Redirect{StdoutPath: path, StdoutAppend: true}
			},
			fields:       []string{"echo", "second"},
			wantFile:     "append.txt",
			wantContents: "first\nsecond\n",
		},
		{
			name: "stderr redirect",
			setup: func(t *testing.T, dir string) parser.Redirect {
				testutils.PrependPATH(t, dir)
				testutils.WriteMockStderrProgram(t, dir, "stderr output")
				return parser.Redirect{StderrPath: filepath.Join(dir, "err.txt")}
			},
			fields:       []string{"mock_stderr"},
			useExternal:  true,
			wantFile:     "err.txt",
			wantContents: "stderr output",
		},
		{
			name: "stderr append redirect",
			setup: func(t *testing.T, dir string) parser.Redirect {
				path := testutils.WriteFileIn(t, dir, "err-append.txt", "first err\n")
				testutils.PrependPATH(t, dir)
				testutils.WriteMockStderrProgram(t, dir, "second err")
				return parser.Redirect{StderrPath: path, StderrAppend: true}
			},
			fields:       []string{"mock_stderr"},
			useExternal:  true,
			wantFile:     "err-append.txt",
			wantContents: "first err\nsecond err",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			redirect := tt.setup(t, dir)
			e := New(strings.NewReader(""), io.Discard, io.Discard)
			state := session.NewSession()

			var err error
			if tt.useExternal {
				err = e.ExecuteExternalForeground(redirect, tt.fields)
			} else {
				_, err = e.ExecuteBuiltin(redirect, state, tt.fields)
			}
			if err != nil {
				t.Fatalf("execute error = %v", err)
			}

			got, err := os.ReadFile(filepath.Join(dir, tt.wantFile))
			if err != nil {
				t.Fatalf("ReadFile() error = %v", err)
			}
			if diff := cmp.Diff(tt.wantContents, string(got)); diff != "" {
				t.Errorf("redirect file mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestNonExitError(t *testing.T) {
	dir := t.TempDir()
	testutils.PrependPATH(t, dir)
	testutils.WriteMockExitProgram(t, dir, 1)
	runErr := exec.Command("mock_exit").Run()
	var exitErr *exec.ExitError
	if !errors.As(runErr, &exitErr) {
		t.Fatalf("mock_exit Run() error = %v, want *exec.ExitError", runErr)
	}

	tests := []struct {
		name string
		err  error
		want error
	}{
		{name: "nil error", err: nil, want: nil},
		{name: "exit error swallowed", err: exitErr, want: nil},
		{name: "other error propagated", err: errors.New("boom"), want: errors.New("boom")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := nonExitError(tt.err)
			if diff := cmp.Diff(tt.want, got, cmp.Comparer(func(a, b error) bool {
				if a == nil && b == nil {
					return true
				}
				if a == nil || b == nil {
					return false
				}
				return a.Error() == b.Error()
			})); diff != "" {
				t.Errorf("nonExitError() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestNonExitErrorWithMockProgram(t *testing.T) {
	dir := t.TempDir()
	testutils.PrependPATH(t, dir)
	testutils.WriteMockExitProgram(t, dir, 42)

	e := New(strings.NewReader(""), io.Discard, io.Discard)
	err := e.ExecuteExternalForeground(parser.Redirect{}, []string{"mock_exit"})
	if err != nil {
		t.Fatalf("ExecuteExternalForeground() error = %v, want nil for exit status", err)
	}
}

func TestExecuteExternalBackground(t *testing.T) {
	dir := t.TempDir()
	testutils.PrependPATH(t, dir)
	name, _ := testutils.WriteMockProgram(t, dir)

	e := New(strings.NewReader(""), io.Discard, io.Discard)

	var startedPID int
	started := make(chan struct{})
	var startedOnce sync.Once
	exited := make(chan struct{})

	pid, err := e.ExecuteExternalBackground(
		parser.Redirect{},
		[]string{name},
		func(childPID int) {
			startedPID = childPID
			startedOnce.Do(func() { close(started) })
		},
		func() { close(exited) },
	)
	if err != nil {
		t.Fatalf("ExecuteExternalBackground() error = %v", err)
	}
	if pid <= 0 {
		t.Fatalf("ExecuteExternalBackground() pid = %d, want > 0", pid)
	}

	select {
	case <-started:
	case <-time.After(5 * time.Second):
		t.Fatal("onStarted was not called within 5s")
	}
	if startedPID != pid {
		t.Errorf("onStarted pid = %d, want %d", startedPID, pid)
	}

	select {
	case <-exited:
	case <-time.After(5 * time.Second):
		t.Fatal("onExit was not called within 5s")
	}
}

func TestExecuteBuiltinStdoutLines(t *testing.T) {
	var out bytes.Buffer
	e := New(strings.NewReader(""), &out, io.Discard)
	state := session.NewSession()

	_, err := e.ExecuteBuiltin(parser.Redirect{}, state, []string{"echo", "a", "b"})
	if err != nil {
		t.Fatalf("ExecuteBuiltin() error = %v", err)
	}
	if diff := cmp.Diff([]string{"a b"}, testutils.OutputLines(out.String()), cmpopts.EquateEmpty()); diff != "" {
		t.Errorf("stdout lines mismatch (-want +got):\n%s", diff)
	}
}
