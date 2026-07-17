//go:build !windows

package external

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestFindExecutableInPath(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(t *testing.T) (command string, wantPath string)
		fileName  string
		filePerms os.FileMode
		pathEnv   string
		wantFound bool
	}{
		{
			name:      "finds executable in path",
			command:   "valid_command",
			fileName:  "valid_command",
			filePerms: 0o755,
			wantFound: true,
		},
		{
			name:      "skips file without execute permission",
			command:   "noexec",
			fileName:  "noexec",
			filePerms: 0o644,
			wantFound: false,
		},
		{
			name:      "skips missing directory",
			command:   "valid_command",
			fileName:  "valid_command",
			filePerms: 0o755,
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
				firstPath := filepath.Join(first, "tool")
				if err := os.WriteFile(firstPath, nil, 0o755); err != nil {
					t.Fatalf("WriteFile() error = %v", err)
				}
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
					if err := os.WriteFile(wantPath, nil, tt.filePerms); err != nil {
						t.Fatalf("WriteFile() error = %v", err)
					}
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
				for name, perms := range map[string]os.FileMode{
					"custom_executable": 0o755,
					"other_tool":        0o755,
					"custom_not_exec":   0o644,
				} {
					if err := os.WriteFile(filepath.Join(dir, name), nil, perms); err != nil {
						t.Fatalf("WriteFile(%q) error = %v", name, err)
					}
				}
				return dir
			},
			want: []string{"custom_executable", "other_tool"},
		},
		{
			name: "skips missing directory",
			setup: func(t *testing.T) string {
				t.Helper()
				dir := t.TempDir()
				if err := os.WriteFile(filepath.Join(dir, "custom_executable"), nil, 0o755); err != nil {
					t.Fatalf("WriteFile() error = %v", err)
				}
				return "missing" + string(os.PathListSeparator) + dir
			},
			want: []string{"custom_executable"},
		},
		{
			name: "empty path returns nil",
			setup: func(t *testing.T) string {
				t.Helper()
				return ""
			},
		},
		{
			name: "deduplicates executables across path entries",
			setup: func(t *testing.T) string {
				t.Helper()
				first := t.TempDir()
				second := t.TempDir()
				for _, dir := range []string{first, second} {
					if err := os.WriteFile(filepath.Join(dir, "shared_tool"), nil, 0o755); err != nil {
						t.Fatalf("WriteFile() error = %v", err)
					}
				}
				return first + string(os.PathListSeparator) + second
			},
			want: []string{"shared_tool"},
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
