package shell

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/codecrafters-io/shell-starter-go/app/testutils"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestExecuteLine_pipeline(t *testing.T) {
	tests := []struct {
		name       string
		line       string
		wantStdout []string
	}{
		{
			name: "pipeline runs last segment output",
			line: "echo hello | echo world",
			wantStdout: []string{
				"world",
			},
		},
		{
			name: "pipeline with builtin middle stage",
			line: "echo hello | type echo",
			wantStdout: []string{
				"echo is a shell builtin",
			},
		},
		{
			name: "pipeline missing command prints error",
			line: "echo hello | missing",
			wantStdout: []string{
				"missing: command not found",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out bytes.Buffer
			s := New(strings.NewReader(""), &out, io.Discard)

			if _, err := s.ExecuteLine(tt.line); err != nil {
				t.Fatalf("ExecuteLine(%q) error = %v", tt.line, err)
			}

			gotStdout := testutils.OutputLines(out.String())
			if diff := cmp.Diff(tt.wantStdout, gotStdout, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("ExecuteLine(%q) stdout mismatch (-want +got):\n%s", tt.line, diff)
			}
		})
	}
}

func TestExecuteLine_pipelineWithRedirectExpansion(t *testing.T) {
	dir := t.TempDir()
	outPath := filepath.Join(dir, "out.txt")

	var out bytes.Buffer
	s := New(strings.NewReader(""), &out, io.Discard)
	s.session.Variables.Set("OUT", outPath)

	if _, err := s.ExecuteLine("echo piped > $OUT"); err != nil {
		t.Fatalf("ExecuteLine() error = %v", err)
	}

	got, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if diff := cmp.Diff("piped\n", string(got)); diff != "" {
		t.Errorf("redirect file mismatch (-want +got):\n%s", diff)
	}
}
