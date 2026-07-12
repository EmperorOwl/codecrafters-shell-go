package builtins

import (
	"bytes"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestPwd(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd() error = %v", err)
	}

	tests := []struct {
		name    string
		wantOut string
	}{
		{name: "prints working directory", wantOut: cwd + "\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout bytes.Buffer
			pwdBuiltin(&stdout)

			if diff := cmp.Diff(tt.wantOut, stdout.String()); diff != "" {
				t.Errorf("pwdBuiltin() stdout mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
