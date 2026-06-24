package shell

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestFindExecutableInPath(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX execute permissions are not enforced on Windows")
	}
	dir1 := t.TempDir()
	dir2 := t.TempDir()

	executable := filepath.Join(dir1, "valid_command")
	if err := os.WriteFile(executable, nil, 0o755); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	nonExecutable := filepath.Join(dir1, "noexec")
	if err := os.WriteFile(nonExecutable, nil, 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	laterExecutable := filepath.Join(dir2, "valid_command")
	if err := os.WriteFile(laterExecutable, nil, 0o755); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	tests := []struct {
		name      string
		command   string
		pathEnv   string
		wantPath  string
		wantFound bool
	}{
		{
			name:      "finds executable in path",
			command:   "valid_command",
			pathEnv:   dir1,
			wantPath:  executable,
			wantFound: true,
		},
		{
			name:      "skips file without execute permission",
			command:   "noexec",
			pathEnv:   dir1,
			wantFound: false,
		},
		{
			name:      "skips missing directory",
			command:   "valid_command",
			pathEnv:   filepath.Join(dir1, "missing") + string(os.PathListSeparator) + dir2,
			wantPath:  laterExecutable,
			wantFound: true,
		},
		{
			name:      "command not found",
			command:   "missing_command",
			pathEnv:   dir1,
			wantFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPath, gotFound := FindExecutableInPath(tt.command, tt.pathEnv)
			if gotFound != tt.wantFound {
				t.Fatalf("FindExecutableInPath(%q) found = %v, want %v", tt.command, gotFound, tt.wantFound)
			}
			if gotPath != tt.wantPath {
				t.Errorf("FindExecutableInPath(%q) path = %q, want %q", tt.command, gotPath, tt.wantPath)
			}
		})
	}
}

func TestTypeOutputExecutable(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX execute permissions are not enforced on Windows")
	}
	dir := t.TempDir()
	command := "mycommand"
	executable := filepath.Join(dir, command)
	if err := os.WriteFile(executable, nil, 0o755); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	t.Setenv("PATH", dir)

	got := TypeOutput(command)
	want := command + " is " + executable
	if got != want {
		t.Errorf("TypeOutput(%q) = %q, want %q", command, got, want)
	}
}
