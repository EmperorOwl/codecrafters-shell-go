package shell

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/codecrafters-io/shell-starter-go/app/builtins"
	"github.com/google/go-cmp/cmp"
)

func TestShellRun(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.txt")
	appendFile := filepath.Join(tmpDir, "append.txt")
	errorsFile := filepath.Join(tmpDir, "errors.txt")
	appendErrorsFile := filepath.Join(tmpDir, "append-errors.txt")
	missingDir := "./does_not_exist"
	cdErr := builtins.CdErrorMessage(missingDir) + "\n"

	tests := []struct {
		name             string
		input            string
		want             string
		wantFilePath     string
		wantFileContents string
	}{
		{
			name:  "echo builtin",
			input: "echo hello\n",
			want:  "$ hello\n$ ",
		},
		{
			name:  "echo hello with EOF does not print prompt again",
			input: "echo hello",
			want:  "$ hello\n",
		},
		{
			name:  "unknown command",
			input: "xyz\n",
			want:  "$ xyz: command not found\n$ ",
		},
		{
			name:  "echo with single-quoted spaces",
			input: "echo 'world     test'\n",
			want:  "$ world     test\n$ ",
		},
		{
			name:  "echo with double-quoted arguments",
			input: `echo "bar"  "shell's"  "foo"` + "\n",
			want:  "$ bar shell's foo\n$ ",
		},
		{
			name:  "echo with escaped spaces",
			input: "echo multiple\\ \\ \\ \\ spaces\n",
			want:  "$ multiple    spaces\n$ ",
		},
		{
			name:  "echo with backslashes inside double quotes",
			input: `echo "inside\"literal_quote."outside\"` + "\n",
			want:  "$ inside\"literal_quote.outside\"\n$ ",
		},
		{
			name:             "echo redirects stdout to file",
			input:            fmt.Sprintf("echo hello > %q\n", outputFile),
			want:             "$ $ ",
			wantFilePath:     outputFile,
			wantFileContents: "hello\n",
		},
		{
			name:             "stderr redirect writes errors to file",
			input:            fmt.Sprintf("cd %s 2> %q\n", missingDir, errorsFile),
			want:             "$ $ ",
			wantFilePath:     errorsFile,
			wantFileContents: cdErr,
		},
		{
			name: "echo appends stdout to file",
			input: fmt.Sprintf(
				"echo first >> %q\necho second >> %q\n",
				appendFile, appendFile,
			),
			want:             "$ $ $ ",
			wantFilePath:     appendFile,
			wantFileContents: "first\nsecond\n",
		},
		{
			name: "stderr append redirect writes errors to file",
			input: fmt.Sprintf(
				"cd %s 2>> %q\ncd %s 2>> %q\n",
				missingDir, appendErrorsFile, missingDir, appendErrorsFile,
			),
			want:             "$ $ $ ",
			wantFilePath:     appendErrorsFile,
			wantFileContents: cdErr + cdErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shell := New()
			var out bytes.Buffer
			err := shell.Run(strings.NewReader(tt.input), &out, io.Discard)
			if err != nil {
				t.Fatalf("Run() error = %v", err)
			}
			if diff := cmp.Diff(tt.want, out.String()); diff != "" {
				t.Errorf("Run() output mismatch (-want +got):\n%s", diff)
			}
			if tt.wantFilePath != "" {
				content, err := os.ReadFile(tt.wantFilePath)
				if err != nil {
					t.Fatalf("ReadFile(%q) error = %v", tt.wantFilePath, err)
				}
				if diff := cmp.Diff(tt.wantFileContents, string(content)); diff != "" {
					t.Errorf("file %q content mismatch (-want +got):\n%s", tt.wantFilePath, diff)
				}
			}
		})
	}
}
