package builtins

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestEcho(t *testing.T) {
	tests := []struct {
		name string
		args []string
		wantOut string
	}{
		{name: "hello world", args: []string{"hello", "world"}, wantOut: "hello world\n"},
		{name: "three words", args: []string{"one", "two", "three"}, wantOut: "one two three\n"},
		{name: "no args", args: nil, wantOut: "\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout bytes.Buffer
			echoBuiltin(&stdout, tt.args)
			if diff := cmp.Diff(tt.wantOut, stdout.String()); diff != "" {
				t.Errorf("echoBuiltin(%v) stdout mismatch (-want +got):\n%s", tt.args, diff)
			}
		})
	}
}
