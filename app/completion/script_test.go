package completion

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestParseCompleterOutput(t *testing.T) {
	tests := []struct {
		name   string
		output []byte
		want   []string
	}{
		{
			name:   "empty output",
			output: nil,
			want:   nil,
		},
		{
			name:   "single candidate",
			output: []byte("run\n"),
			want:   []string{"run"},
		},
		{
			name:   "multiple candidates",
			output: []byte("stash\nstatus\n"),
			want:   []string{"stash", "status"},
		},
		{
			name:   "trims trailing newline",
			output: []byte("run"),
			want:   []string{"run"},
		},
		{
			name:   "windows line endings",
			output: []byte("run\r\nstatus\r\n"),
			want:   []string{"run", "status"},
		},
		{
			name:   "skips blank lines",
			output: []byte("run\n\nstatus\n"),
			want:   []string{"run", "status"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseCompleterOutput(tt.output)
			if diff := cmp.Diff(tt.want, got, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("parseCompleterOutput() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestRunCompleter(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(t *testing.T) (CompleterOptions, []string)
		wantErr bool
	}{
		{
			name: "runs script and parses candidates",
			setup: func(t *testing.T) (CompleterOptions, []string) {
				t.Helper()
				scriptPath := writeCompleterScript(t, "echo run\necho status\n")
				return CompleterOptions{
					Path:         scriptPath,
					Command:      "git",
					CurrentWord:  "st",
					PreviousWord: "git",
					CompLine:     "git st",
					CompPoint:    6,
				}, []string{"run", "status"}
			},
		},
		{
			name: "passes completion env vars to script",
			setup: func(t *testing.T) (CompleterOptions, []string) {
				t.Helper()
				if runtime.GOOS == "windows" {
					t.Skip("COMP_LINE env verification requires a shell script")
				}
				scriptPath := writeCompleterScript(t, "printf '%s\\n' \"$COMP_LINE\" \"$COMP_POINT\"\n")
				return CompleterOptions{
					Path:        scriptPath,
					Command:     "git",
					CompLine:    "git checkout main",
					CompPoint:   15,
					CurrentWord: "main",
				}, []string{"git checkout main", "15"}
			},
		},
		{
			name: "missing script returns error",
			setup: func(t *testing.T) (CompleterOptions, []string) {
				t.Helper()
				return CompleterOptions{
					Path: filepath.Join(t.TempDir(), "missing"),
				}, nil
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts, want := tt.setup(t)

			got, err := RunCompleter(opts)
			if tt.wantErr {
				if err == nil {
					t.Fatal("RunCompleter() error = nil, want non-nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("RunCompleter() error = %v", err)
			}
			if diff := cmp.Diff(want, got, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("RunCompleter() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func writeCompleterScript(t *testing.T, body string) string {
	t.Helper()

	dir := t.TempDir()
	if runtime.GOOS == "windows" {
		path := filepath.Join(dir, "completer.bat")
		content := "@echo off\r\n" + body
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatalf("WriteFile() error = %v", err)
		}
		return path
	}

	path := filepath.Join(dir, "completer.sh")
	content := "#!/bin/sh\n" + body
	if err := os.WriteFile(path, []byte(content), 0o755); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	return path
}
