package variables

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestIsValidIdentifier(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want bool
	}{
		{name: "accepts letter start", arg: "foo", want: true},
		{name: "accepts underscore start", arg: "_FOO", want: true},
		{name: "accepts digits after first character", arg: "foo1", want: true},
		{name: "rejects digit start", arg: "67", want: false},
		{name: "rejects empty", arg: "", want: false},
		{name: "rejects hyphen", arg: "foo-bar", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsValidIdentifier(tt.arg)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("IsValidIdentifier(%q) mismatch (-want +got):\n%s", tt.arg, diff)
			}
		})
	}
}
