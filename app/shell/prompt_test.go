package shell

import (
	"bytes"
	"testing"
)

func TestWritePrompt(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{name: "prints dollar prompt", want: "$ "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			WritePrompt(&buf)
			if got := buf.String(); got != tt.want {
				t.Errorf("WritePrompt() = %q, want %q", got, tt.want)
			}
		})
	}
}
