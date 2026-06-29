package shell

import (
	"bufio"
	"bytes"
	"strings"
	"testing"
)

func TestReadLineRaw(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantLine string
		wantEOF  bool
		wantOut  string
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
			name:     "tab lists ambiguous matches",
			input:    "e\t\n",
			wantLine: "e",
			wantOut:  "\r$ e\r\necho\r\nexit\r\n\r\033[K$ e\r\n",
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

			gotLine, gotEOF, err := readLine(reader, &out, true)
			if err != nil {
				t.Fatalf("readLine() error = %v", err)
			}
			if gotLine != tt.wantLine {
				t.Errorf("readLine() line = %q, want %q", gotLine, tt.wantLine)
			}
			if gotEOF != tt.wantEOF {
				t.Errorf("readLine() eof = %v, want %v", gotEOF, tt.wantEOF)
			}
			if got := out.String(); got != tt.wantOut {
				t.Errorf("readLine() output = %q, want %q", got, tt.wantOut)
			}
		})
	}
}

func TestReadLineRaw_SkipsLFAfterCR(t *testing.T) {
	skipNextLF = false

	var out bytes.Buffer

	line, eof, err := readLine(bufio.NewReader(strings.NewReader("hi\r")), &out, true)
	if err != nil {
		t.Fatalf("first readLine() error = %v", err)
	}
	if line != "hi" {
		t.Errorf("first readLine() line = %q, want %q", line, "hi")
	}
	if eof {
		t.Error("first readLine() eof = true, want false")
	}
	if !skipNextLF {
		t.Error("skipNextLF = false after CR, want true")
	}

	line, eof, err = readLine(bufio.NewReader(strings.NewReader("\n")), &out, true)
	if err != nil {
		t.Fatalf("second readLine() error = %v", err)
	}
	if line != "" {
		t.Errorf("second readLine() line = %q, want empty (stray LF skipped)", line)
	}
	if !eof {
		t.Error("second readLine() eof = false, want true")
	}
	if skipNextLF {
		t.Error("skipNextLF = true after skipped LF, want false")
	}

	wantOut := "\r$ hi\r\n\r$ "
	if got := out.String(); got != wantOut {
		t.Errorf("combined output = %q, want %q", got, wantOut)
	}
}
