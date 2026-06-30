package terminal

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestWrapWriter(t *testing.T) {
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
			writer := WrapWriter(&buf, true)
			if _, err := writer.Write([]byte(tt.input)); err != nil {
				t.Fatalf("Write() error = %v", err)
			}
			if diff := cmp.Diff(tt.want, buf.String()); diff != "" {
				t.Errorf("Write(%q) mismatch (-want +got):\n%s", tt.input, diff)
			}
		})
	}
}
