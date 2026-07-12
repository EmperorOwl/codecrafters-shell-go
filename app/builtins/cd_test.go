package builtins

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type cdFixture struct {
	startDir   string
	directory  string
	wantDir    string
	wantErr string
}

func TestCd(t *testing.T) {
	tests := []struct {
		name  string
		setup func(t *testing.T, base string) cdFixture
	}{
		{
			name: "empty directory",
			setup: func(t *testing.T, base string) cdFixture {
				local := makeDir(t, base, "local")
				return cdFixture{
					startDir:  local,
					directory: "",
					wantDir:   local,
				}
			},
		},
		{
			name: "absolute path",
			setup: func(t *testing.T, base string) cdFixture {
				localBin := makeDir(t, base, "local/bin")
				return cdFixture{
					startDir:  base,
					directory: localBin,
					wantDir:   localBin,
				}
			},
		},
		{
			name: "relative path with dot prefix",
			setup: func(t *testing.T, base string) cdFixture {
				localBin := makeDir(t, base, "local/bin")
				return cdFixture{
					startDir:  base,
					directory: "./local/bin",
					wantDir:   localBin,
				}
			},
		},
		{
			name: "parent directories",
			setup: func(t *testing.T, base string) cdFixture {
				nested := makeDir(t, base, "a/b")
				return cdFixture{
					startDir:  nested,
					directory: "../../",
					wantDir:   base,
				}
			},
		},
		{
			name: "subdirectory without dot prefix",
			setup: func(t *testing.T, base string) cdFixture {
				local := makeDir(t, base, "local")
				return cdFixture{
					startDir:  base,
					directory: "local",
					wantDir:   local,
				}
			},
		},
		{
			name: "home directory",
			setup: func(t *testing.T, base string) cdFixture {
				home := makeDir(t, base, "home")
				setHome(t, home)
				return cdFixture{
					startDir:  base,
					directory: "~",
					wantDir:   home,
				}
			},
		},
		{
			name: "missing relative directory",
			setup: func(t *testing.T, base string) cdFixture {
				return cdFixture{
					startDir:   base,
					directory:  "./does_not_exist",
					wantDir:    base,
					wantErr: "cd: ./does_not_exist: No such file or directory\n",
				}
			},
		},
		{
			name: "missing absolute directory",
			setup: func(t *testing.T, base string) cdFixture {
				invalidAbs := invalidAbsolutePath()
				return cdFixture{
					startDir:   base,
					directory:  invalidAbs,
					wantDir:    base,
					wantErr: fmt.Sprintf("cd: %s: No such file or directory\n", invalidAbs),
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			base := t.TempDir()
			fixture := tt.setup(t, base)

			t.Chdir(fixture.startDir)

			var stderr bytes.Buffer
			cdBuiltin(&stderr, fixture.directory)

			if diff := cmp.Diff(fixture.wantErr, stderr.String()); diff != "" {
				t.Errorf("cdBuiltin() stderr mismatch (-want +got):\n%s", diff)
			}

			cwd, err := os.Getwd()
			if err != nil {
				t.Fatalf("Getwd() error = %v", err)
			}
			if diff := cmp.Diff(fixture.wantDir, cwd); diff != "" {
				t.Errorf("Getwd() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func setHome(t *testing.T, home string) {
	t.Helper()
	t.Setenv("HOME", home)
	if runtime.GOOS == "windows" {
		t.Setenv("USERPROFILE", home)
	}
}

func makeDir(t *testing.T, base, rel string) string {
	t.Helper()
	path := filepath.Join(base, filepath.FromSlash(rel))
	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatalf("MkdirAll(%q) error = %v", path, err)
	}
	return path
}

func invalidAbsolutePath() string {
	invalidAbs := "/does_not_exist"
	if runtime.GOOS == "windows" {
		vol := os.Getenv("SystemDrive")
		if vol == "" {
			vol = "C:"
		}
		invalidAbs = filepath.Join(vol+string(filepath.Separator), "does_not_exist")
	}
	return invalidAbs
}
