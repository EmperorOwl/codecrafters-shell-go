package shellio

import (
	"bufio"
	"bytes"
	"strings"
	"testing"
)

func TestReadLineRaw(t *testing.T) {
	builtins := []string{"cd", "echo", "exit", "pwd", "type"}

	tests := []struct {
		name          string
		executables   []string
		input         string
		wantLine      string
		wantEOF       bool
		wantOut       string
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

			gotLine, gotEOF, err := ReadLine(reader, &out, true, builtins, tt.executables)
			if err != nil {
				t.Fatalf("ReadLine() error = %v", err)
			}
			if gotLine != tt.wantLine {
				t.Errorf("ReadLine() line = %q, want %q", gotLine, tt.wantLine)
			}
			if gotEOF != tt.wantEOF {
				t.Errorf("ReadLine() eof = %v, want %v", gotEOF, tt.wantEOF)
			}
			if got := out.String(); got != tt.wantOut {
				t.Errorf("ReadLine() output = %q, want %q", got, tt.wantOut)
			}
		})
	}
}

func TestReadLineRaw_SkipsLFAfterCR(t *testing.T) {
	builtins := []string{"cd", "echo", "exit", "pwd", "type"}
	skipNextLF = false

	var out bytes.Buffer

	line, eof, err := ReadLine(bufio.NewReader(strings.NewReader("hi\r")), &out, true, builtins, nil)
	if err != nil {
		t.Fatalf("first ReadLine() error = %v", err)
	}
	if line != "hi" {
		t.Errorf("first ReadLine() line = %q, want %q", line, "hi")
	}
	if eof {
		t.Error("first ReadLine() eof = true, want false")
	}
	if !skipNextLF {
		t.Error("skipNextLF = false after CR, want true")
	}

	line, eof, err = ReadLine(bufio.NewReader(strings.NewReader("\n")), &out, true, builtins, nil)
	if err != nil {
		t.Fatalf("second ReadLine() error = %v", err)
	}
	if line != "" {
		t.Errorf("second ReadLine() line = %q, want empty (stray LF skipped)", line)
	}
	if !eof {
		t.Error("second ReadLine() eof = false, want true")
	}
	if skipNextLF {
		t.Error("skipNextLF = true after skipped LF, want false")
	}

	wantOut := "\r$ hi\r\n\r$ "
	if got := out.String(); got != wantOut {
		t.Errorf("combined output = %q, want %q", got, wantOut)
	}
}
