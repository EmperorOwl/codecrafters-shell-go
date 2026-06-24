//go:build !windows

package path

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindExecutableInPath(t *testing.T) {
	tests := []struct {
		name      string
		command   string
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()

			var filePath string
			if tt.fileName != "" {
				filePath = filepath.Join(dir, tt.fileName)
				if err := os.WriteFile(filePath, nil, tt.filePerms); err != nil {
					t.Fatalf("WriteFile() error = %v", err)
				}
			}

			pathEnv := dir
			if tt.pathEnv != "" {
				pathEnv = tt.pathEnv
			}

			gotPath, gotFound := FindExecutableInPath(tt.command, pathEnv)
			if gotFound != tt.wantFound {
				t.Fatalf("FindExecutableInPath(%q) found = %v, want %v", tt.command, gotFound, tt.wantFound)
			}
			if tt.wantFound && gotPath != filePath {
				t.Errorf("FindExecutableInPath(%q) path = %q, want %q", tt.command, gotPath, filePath)
			}
		})
	}
}
