package shell

import (
	"bytes"
	"testing"
)

func TestTerminalWriter(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "single newline",
			input: "hello\n",
			want:  "hello\r\n",
		},
		{
			name:  "multiple newlines",
			input: "one\ntwo\n",
			want:  "one\r\ntwo\r\n",
		},
		{
			name:  "preserves existing crlf",
			input: "one\r\ntwo",
			want:  "one\r\ntwo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			writer := terminalWriter{w: &buf}
			if _, err := writer.Write([]byte(tt.input)); err != nil {
				t.Fatalf("Write() error = %v", err)
			}
			if got := buf.String(); got != tt.want {
				t.Errorf("Write(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
