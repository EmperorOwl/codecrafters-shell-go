//go:build windows

package shell

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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()

			var filePath string
			if tt.fileName != "" {
				filePath = filepath.Join(dir, tt.fileName)
				if err := os.WriteFile(filePath, nil, 0o644); err != nil {
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
