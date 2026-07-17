package builtins

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestExit(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		wantExit bool
	}{
		{name: "no args exits shell", wantExit: true},
		{name: "ignores status arg", args: []string{"42"}, wantExit: true},
		{name: "ignores multiple args", args: []string{"1", "2"}, wantExit: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotExit, err := exitHandler(&Context{}, tt.args)
			if err != nil {
				t.Fatalf("exitHandler() error = %v", err)
			}
			if diff := cmp.Diff(tt.wantExit, gotExit); diff != "" {
				t.Errorf("exitHandler(%v) exit mismatch (-want +got):\n%s", tt.args, diff)
			}
		})
	}
}
