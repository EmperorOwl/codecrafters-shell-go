package builtins

import (
	"bytes"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestPwd(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(t *testing.T) string
		wantOut func(t *testing.T, setupDir string) string
	}{
		{
			name: "prints working directory",
			wantOut: func(t *testing.T, _ string) string {
				t.Helper()
				cwd, err := os.Getwd()
				if err != nil {
					t.Fatalf("Getwd() error = %v", err)
				}
				return cwd + "\n"
			},
		},
		{
			name: "prints directory after chdir",
			setup: func(t *testing.T) string {
				t.Helper()
				dir := t.TempDir()
				t.Chdir(dir)
				return dir
			},
			wantOut: func(_ *testing.T, setupDir string) string {
				return setupDir + "\n"
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupDir := ""
			if tt.setup != nil {
				setupDir = tt.setup(t)
			}

			var stdout bytes.Buffer
			pwdBuiltin(&stdout)

			wantOut := tt.wantOut(t, setupDir)
			if diff := cmp.Diff(wantOut, stdout.String()); diff != "" {
				t.Errorf("pwdBuiltin() stdout mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
