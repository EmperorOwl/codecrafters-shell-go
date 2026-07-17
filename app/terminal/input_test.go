package terminal

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestReadLine_cookedMode(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantEOF bool
	}{
		{
			name:  "normal line",
			input: "hello\n",
			want:  "hello",
		},
		{
			name:  "trims surrounding whitespace",
			input: "  hello  \n",
			want:  "hello",
		},
		{
			name:  "empty line",
			input: "\n",
			want:  "",
		},
		{
			name:    "eof with partial line",
			input:   "partial",
			want:    "partial",
			wantEOF: true,
		},
		{
			name:    "eof with no input",
			input:   "",
			want:    "",
			wantEOF: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bufio.NewReader(bytes.NewBufferString(tt.input))
			var out bytes.Buffer

			got, eof, err := readLine(reader, &out, false, nil, nil)
			if err != nil {
				t.Fatalf("readLine() error = %v", err)
			}
			if eof != tt.wantEOF {
				t.Fatalf("readLine() eof = %v, want %v", eof, tt.wantEOF)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("readLine() line mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestReadLineRaw(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		history []string
		want    string
		wantEOF bool
	}{
		{
			name:    "up arrow recalls latest history",
			input:   "\x1b[A\r",
			history: []string{"echo hello", "echo world"},
			want:    "echo world",
		},
		{
			name:    "up arrow twice recalls earlier history",
			input:   "\x1b[A\x1b[A\r",
			history: []string{"echo hello", "echo world"},
			want:    "echo hello",
		},
		{
			name:    "down arrow after up recalls newer history",
			input:   "\x1b[A\x1b[A\x1b[B\r",
			history: []string{"echo hello", "echo world"},
			want:    "echo world",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := stubHistoryHandler{commands: tt.history}
			input := bytes.NewBufferString(tt.input)
			var out bytes.Buffer

			got, eof, err := readLineRaw(bufio.NewReader(input), &out, nil, handler)
			if err != nil {
				t.Fatalf("readLineRaw() error = %v", err)
			}
			if eof != tt.wantEOF {
				t.Fatalf("readLineRaw() eof = %v, want %v", eof, tt.wantEOF)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("readLineRaw() line mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
