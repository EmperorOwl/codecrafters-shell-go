package shell

import (
	"bytes"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"testing"
)

func TestShellRun(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.txt")
	errorsFile := filepath.Join(tmpDir, "errors.txt")

	tests := []struct {
		name  string
		input string
		want  string
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
			name:  "echo redirects stdout to file",
			input: fmt.Sprintf("echo hello > %q\ncat %q\n", outputFile, outputFile),
			want:  "$ $ hello\n$ ",
		},
		{
			name:  "echo with stderr redirect prints to terminal",
			input: fmt.Sprintf("echo Maria file cannot be found 2> %q\n", errorsFile),
			want:  "$ Maria file cannot be found\n$ ",
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
			if got := out.String(); got != tt.want {
				t.Errorf("Run() output = %q, want %q", got, tt.want)
			}
		})
	}
}
