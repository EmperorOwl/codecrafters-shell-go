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

func TestFindAllExecutablesInPath(t *testing.T) {
	dir := t.TempDir()
	executable := filepath.Join(dir, "custom_executable")
	if err := os.WriteFile(executable, nil, 0o755); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "other_tool"), nil, 0o755); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "custom_not_exec"), nil, 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	tests := []struct {
		name    string
		pathEnv string
		want    []string
	}{
		{
			name:    "lists executables in path",
			pathEnv: dir,
			want:    []string{"custom_executable", "other_tool"},
		},
		{
			name:    "skips missing directory",
			pathEnv: "missing:" + dir,
			want:    []string{"custom_executable", "other_tool"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("PATH", tt.pathEnv)

			got := FindAllExecutablesInPath()
			if len(got) != len(tt.want) {
				t.Fatalf("FindAllExecutablesInPath() = %v, want %v", got, tt.want)
			}
			for i := range tt.want {
				if got[i] != tt.want[i] {
					t.Errorf("FindAllExecutablesInPath()[%d] = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}
