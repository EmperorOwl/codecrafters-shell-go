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
		want string
	}{
		{name: "hello world", args: []string{"hello", "world"}, want: "hello world\n"},
		{name: "three words", args: []string{"one", "two", "three"}, want: "one two three\n"},
		{name: "no args", args: nil, want: "\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out bytes.Buffer
			Echo(&out, tt.args)
			if diff := cmp.Diff(tt.want, out.String()); diff != "" {
				t.Errorf("Echo() output mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
