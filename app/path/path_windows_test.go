//go:build windows

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
			t.Setenv("PATH", pathEnv)

			gotPath, gotFound := FindExecutableInPath(tt.command)
			if gotFound != tt.wantFound {
				t.Fatalf("FindExecutableInPath(%q) found = %v, want %v", tt.command, gotFound, tt.wantFound)
			}
			if tt.wantFound && gotPath != filePath {
				t.Errorf("FindExecutableInPath(%q) path = %q, want %q", tt.command, gotPath, filePath)
			}
		})
	}
}

func TestFindMatchingExecutablesInPath(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "custom_executable.exe"), nil, 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "custom_not_exec.txt"), nil, 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	tests := []struct {
		name    string
		prefix  string
		pathEnv string
		want    []string
	}{
		{
			name:    "completes executable in path",
			prefix:  "custom",
			pathEnv: dir,
			want:    []string{"custom_executable.exe"},
		},
		{
			name:    "no match",
			prefix:  "missing",
			pathEnv: dir,
			want:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("PATH", tt.pathEnv)

			got := FindMatchingExecutablesInPath(tt.prefix)
			if len(got) != len(tt.want) {
				t.Fatalf("FindMatchingExecutablesInPath(%q) = %v, want %v", tt.prefix, got, tt.want)
			}
			for i := range tt.want {
				if got[i] != tt.want[i] {
					t.Errorf("FindMatchingExecutablesInPath(%q)[%d] = %q, want %q", tt.prefix, i, got[i], tt.want[i])
				}
			}
		})
	}
}
