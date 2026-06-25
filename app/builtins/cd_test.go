package builtins

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestCd(t *testing.T) {
	// Fixture layout:
	//   base/
	//     local/
	//       bin/
	//     a/
	//       b/
	//     home/
	base := t.TempDir()
	localDir := filepath.Join(base, "local")
	localBinDir := filepath.Join(localDir, "bin")
	nestedDir := filepath.Join(base, "a", "b")
	homeDir := filepath.Join(base, "home")
	if err := os.MkdirAll(localBinDir, 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.MkdirAll(nestedDir, 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.MkdirAll(homeDir, 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}

	// Unix absolute paths start with /. On Windows, /path is not a root
	// absolute path, so use a drive-letter path (e.g. C:\...) instead.
	invalidAbs := "/does_not_exist"
	if runtime.GOOS == "windows" {
		vol := os.Getenv("SystemDrive")
		if vol == "" {
			vol = "C:"
		}
		invalidAbs = filepath.Join(vol+string(filepath.Separator), "does_not_exist")
	}

	tests := []struct {
		name       string
		startDir   string
		directory  string
		home       string
		wantDir    string
		wantOutput string
	}{
		{
			name:      "absolute path",
			startDir:  nestedDir,
			directory: localBinDir,
			wantDir:   localBinDir,
		},
		{
			name:      "relative path with dot prefix",
			startDir:  base,
			directory: "./local/bin",
			wantDir:   localBinDir,
		},
		{
			name:      "parent directories",
			startDir:  nestedDir,
			directory: "../../",
			wantDir:   base,
		},
		{
			name:      "subdirectory without dot prefix",
			startDir:  base,
			directory: "local",
			wantDir:   localDir,
		},
		{
			name:      "home directory",
			startDir:  nestedDir,
			directory: "~",
			home:      homeDir,
			wantDir:   homeDir,
		},
		{
			name:       "missing relative directory",
			startDir:   base,
			directory:  "./does_not_exist",
			wantDir:    base,
			wantOutput: "cd: ./does_not_exist: No such file or directory\n",
		},
		{
			name:       "missing absolute directory",
			startDir:   base,
			directory:  invalidAbs,
			wantDir:    base,
			wantOutput: "cd: " + invalidAbs + ": No such file or directory\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.home != "" {
				t.Setenv("HOME", tt.home)
			}

			t.Chdir(tt.startDir)

			var errOut bytes.Buffer
			Cd(&errOut, tt.directory)

			if got := errOut.String(); got != tt.wantOutput {
				t.Errorf("Cd() stderr = %q, want %q", got, tt.wantOutput)
			}

			cwd, err := os.Getwd()
			if err != nil {
				t.Fatalf("Getwd() error = %v", err)
			}
			if cwd != tt.wantDir {
				t.Errorf("Getwd() = %q, want %q", cwd, tt.wantDir)
			}
		})
	}
}
