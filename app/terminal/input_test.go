package terminal

import (
	"bufio"
	"bytes"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestReadLineRaw(t *testing.T) {
	builtins := []string{"cd", "echo", "exit", "pwd", "type"}

	tests := []struct {
		name        string
		executables []string
		input       string
		wantLine    string
		wantEOF     bool
		wantOut     string
	}{
		{
			name:     "unix enter submits line",
			input:    "hello\n",
			wantLine: "hello",
			wantOut:  "\r$ hello\r\n",
		},
		{
			name:     "windows enter submits line",
			input:    "hello\r",
			wantLine: "hello",
			wantOut:  "\r$ hello\r\n",
		},
		{
			name:     "tab completes echo",
			input:    "ech\t\n",
			wantLine: "echo ",
			wantOut:  "\r$ ech\r\033[K$ echo \r\n",
		},
		{
			name:     "tab completes exit",
			input:    "exi\t\n",
			wantLine: "exit ",
			wantOut:  "\r$ exi\r\033[K$ exit \r\n",
		},
		{
			name:     "tab rings bell on ambiguous prefix",
			input:    "e\t\n",
			wantLine: "e",
			wantOut:  "\r$ e\a\r\n",
		},
		{
			name:     "double tab lists ambiguous matches",
			input:    "e\t\t\n",
			wantLine: "e",
			wantOut:  "\r$ e\a\r\necho  exit\r\n\r\033[K$ e\r\n",
		},
		{
			name:     "tab rings bell on no match",
			input:    "xyz\t\n",
			wantLine: "xyz",
			wantOut:  "\r$ xyz\a\r\n",
		},
		{
			name:        "double tab lists executable matches",
			executables: []string{"xyz_bar", "xyz_baz", "xyz_quz"},
			input:       "xyz_\t\t\n",
			wantLine:    "xyz_",
			wantOut:     "\r$ xyz_\a\r\nxyz_bar  xyz_baz  xyz_quz\r\n\r\033[K$ xyz_\r\n",
		},
		{
			name:        "progressive tab completes to longest common prefix",
			executables: []string{"xyz_foo", "xyz_foo_bar", "xyz_foo_bar_baz"},
			input:       "xyz_\t_\t_\t\n",
			wantLine:    "xyz_foo_bar_baz ",
			wantOut:     "\r$ xyz_\r\033[K$ xyz_foo_\r\033[K$ xyz_foo_bar_\r\033[K$ xyz_foo_bar_baz \r\n",
		},
		{
			name:     "backspace removes character",
			input:    "ab\x08\n",
			wantLine: "a",
			wantOut:  "\r$ ab\b \b\r\n",
		},
		{
			name:     "eof on empty input",
			input:    "",
			wantLine: "",
			wantEOF:  true,
			wantOut:  "\r$ ",
		},
		{
			name:     "eof with partial line",
			input:    "hello",
			wantLine: "hello",
			wantEOF:  true,
			wantOut:  "\r$ hello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			skipNextLF = false

			reader := bufio.NewReader(strings.NewReader(tt.input))
			var out bytes.Buffer

			gotLine, gotEOF, err := ReadLine(reader, &out, true, builtins, tt.executables, nil, nil)
			if err != nil {
				t.Fatalf("ReadLine() error = %v", err)
			}
			if diff := cmp.Diff(tt.wantLine, gotLine); diff != "" {
				t.Errorf("ReadLine() line mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tt.wantEOF, gotEOF); diff != "" {
				t.Errorf("ReadLine() eof mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tt.wantOut, out.String()); diff != "" {
				t.Errorf("ReadLine() output mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestReadLineRaw_SkipsLFAfterCR(t *testing.T) {
	builtins := []string{"cd", "echo", "exit", "pwd", "type"}
	skipNextLF = false

	var out bytes.Buffer

	line, eof, err := ReadLine(bufio.NewReader(strings.NewReader("hi\r")), &out, true, builtins, nil, nil, nil)
	if err != nil {
		t.Fatalf("first ReadLine() error = %v", err)
	}
	if diff := cmp.Diff("hi", line); diff != "" {
		t.Errorf("first ReadLine() line mismatch (-want +got):\n%s", diff)
	}
	if eof {
		t.Error("first ReadLine() eof = true, want false")
	}
	if !skipNextLF {
		t.Error("skipNextLF = false after CR, want true")
	}

	line, eof, err = ReadLine(bufio.NewReader(strings.NewReader("\n")), &out, true, builtins, nil, nil, nil)
	if err != nil {
		t.Fatalf("second ReadLine() error = %v", err)
	}
	if diff := cmp.Diff("", line); diff != "" {
		t.Errorf("second ReadLine() line mismatch (-want +got):\n%s", diff)
	}
	if !eof {
		t.Error("second ReadLine() eof = false, want true")
	}
	if skipNextLF {
		t.Error("skipNextLF = true after skipped LF, want false")
	}

	wantOut := "\r$ hi\r\n\r$ "
	if diff := cmp.Diff(wantOut, out.String()); diff != "" {
		t.Errorf("combined output mismatch (-want +got):\n%s", diff)
	}
}

func TestReadLineRaw_tabCompletesFileOnSecondPrompt(t *testing.T) {
	builtins := []string{"cd", "echo", "exit", "pwd", "type"}
	listFiles := func(dir string) []string {
		if dir == "" {
			return []string{"app/"}
		}
		return nil
	}

	var out bytes.Buffer

	skipNextLF = false
	line, eof, err := ReadLine(bufio.NewReader(strings.NewReader("ls a\t\r")), &out, true, builtins, nil, listFiles, nil)
	if err != nil {
		t.Fatalf("first ReadLine() error = %v", err)
	}
	if diff := cmp.Diff("ls app/", line); diff != "" {
		t.Errorf("first ReadLine() line mismatch (-want +got):\n%s", diff)
	}
	if eof {
		t.Error("first ReadLine() eof = true, want false")
	}

	// Simulates a second prompt after an external command. The shell re-enables
	// raw mode via Session.PrepareRead; tab completion must still work here.
	skipNextLF = false
	line, eof, err = ReadLine(bufio.NewReader(strings.NewReader("ls a\t\r")), &out, true, builtins, nil, listFiles, nil)
	if err != nil {
		t.Fatalf("second ReadLine() error = %v", err)
	}
	if diff := cmp.Diff("ls app/", line); diff != "" {
		t.Errorf("second ReadLine() line mismatch (-want +got):\n%s", diff)
	}
	if eof {
		t.Error("second ReadLine() eof = true, want false")
	}
}
