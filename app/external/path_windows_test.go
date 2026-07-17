//go:build windows

package external

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/codecrafters-io/shell-starter-go/app/testutils"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestFindExecutableInPath(t *testing.T) {
	tests := []struct {
		name      string
		command   string
		setup     func(t *testing.T) (command string, wantPath string)
		fileName  string
		pathEnv   string
		wantFound bool
	}{
		{
			name:      "finds executable without extension",
			command:   "mycommand",
			fileName:  "mycommand.exe",
			wantFound: true,
		},
		{
			name:      "finds executable with extension",
			command:   "mycommand.exe",
			fileName:  "mycommand.exe",
			wantFound: true,
		},
		{
			name:      "skips non-executable extension",
			command:   "noexec",
			fileName:  "noexec.txt",
			wantFound: false,
		},
		{
			name:      "skips missing directory",
			command:   "valid_command",
			fileName:  "valid_command.exe",
			pathEnv:   "missing",
			wantFound: false,
		},
		{
			name:      "command not found",
			command:   "missing_command",
			wantFound: false,
		},
		{
			name: "finds first match across path entries",
			setup: func(t *testing.T) (string, string) {
				t.Helper()
				first := t.TempDir()
				second := t.TempDir()
				firstPath := filepath.Join(first, "tool.exe")
				testutils.CreatePath(t, first, "tool.exe")
				t.Setenv("PATH", first+string(os.PathListSeparator)+second)
				return "tool", firstPath
			},
			wantFound: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			command := tt.command
			wantPath := ""

			if tt.setup != nil {
				command, wantPath = tt.setup(t)
			} else {
				dir := t.TempDir()
				if tt.fileName != "" {
					wantPath = filepath.Join(dir, tt.fileName)
					testutils.CreatePath(t, dir, tt.fileName)
				}

				pathEnv := dir
				if tt.pathEnv != "" {
					pathEnv = tt.pathEnv
				}
				t.Setenv("PATH", pathEnv)
			}

			gotPath, gotFound := FindExecutableInPath(command)
			if diff := cmp.Diff(tt.wantFound, gotFound); diff != "" {
				t.Fatalf("FindExecutableInPath(%q) found mismatch (-want +got):\n%s", command, diff)
			}
			if tt.wantFound {
				if diff := cmp.Diff(wantPath, gotPath); diff != "" {
					t.Errorf("FindExecutableInPath(%q) path mismatch (-want +got):\n%s", command, diff)
				}
			}
		})
	}
}

func TestFindAllExecutablesInPath(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(t *testing.T) string
		want    []string
	}{
		{
			name: "lists executables in path",
			setup: func(t *testing.T) string {
				t.Helper()
				dir := t.TempDir()
				testutils.CreatePath(t, dir, "custom_executable.exe")
				testutils.CreatePath(t, dir, "custom_not_exec.txt")
				return dir
			},
			want: []string{"custom_executable.exe"},
		},
		{
			name: "skips missing directory",
			setup: func(t *testing.T) string {
				t.Helper()
				dir := t.TempDir()
				testutils.CreatePath(t, dir, "custom_executable.exe")
				return "missing" + string(os.PathListSeparator) + dir
			},
			want: []string{"custom_executable.exe"},
		},
		{
			name: "empty path returns nil",
			setup: func(t *testing.T) string {
				t.Helper()
				return ""
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("PATH", tt.setup(t))

			got := FindAllExecutablesInPath()
			if diff := cmp.Diff(tt.want, got, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("FindAllExecutablesInPath() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
