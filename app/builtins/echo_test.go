package builtins

import (
	"bytes"
	"testing"
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
			if got := out.String(); got != tt.want {
				t.Errorf("Echo() output = %q, want %q", got, tt.want)
			}
		})
	}
}
