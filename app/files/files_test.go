package files

import (
	"os"
	"path/filepath"
	"testing"
)

func TestListInDir(t *testing.T) {
	dir := t.TempDir()
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd() error = %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("Chdir() error = %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(origDir)
	})

	nested := filepath.Join(dir, "path", "to")
	if err := os.MkdirAll(nested, 0755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	for _, spec := range []struct {
		path string
		name string
	}{
		{dir, "readme.txt"},
		{dir, "notes.md"},
		{nested, "file.txt"},
	} {
		if err := os.WriteFile(filepath.Join(spec.path, spec.name), nil, 0644); err != nil {
			t.Fatalf("WriteFile(%q) error = %v", spec.name, err)
		}
	}
	if err := os.Mkdir(filepath.Join(dir, "subdir"), 0755); err != nil {
		t.Fatalf("Mkdir() error = %v", err)
	}

	tests := []struct {
		name string
		dir  string
		want []string
	}{
		{
			name: "current directory",
			dir:  "",
			want: []string{"notes.md", "readme.txt"},
		},
		{
			name: "nested directory",
			dir:  "path/to/",
			want: []string{"file.txt"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ListInDir(tt.dir)
			if len(got) != len(tt.want) {
				t.Fatalf("ListInDir(%q) = %v, want %v", tt.dir, got, tt.want)
			}
			for i := range tt.want {
				if got[i] != tt.want[i] {
					t.Errorf("ListInDir(%q)[%d] = %q, want %q", tt.dir, i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestListInCurrentDir(t *testing.T) {
	dir := t.TempDir()
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd() error = %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("Chdir() error = %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(origDir)
	})

	for _, name := range []string{"readme.txt", "notes.md"} {
		if err := os.WriteFile(filepath.Join(dir, name), nil, 0644); err != nil {
			t.Fatalf("WriteFile(%q) error = %v", name, err)
		}
	}
	if err := os.Mkdir(filepath.Join(dir, "subdir"), 0755); err != nil {
		t.Fatalf("Mkdir() error = %v", err)
	}

	got := ListInCurrentDir()
	want := []string{"notes.md", "readme.txt"}
	if len(got) != len(want) {
		t.Fatalf("ListInCurrentDir() = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("ListInCurrentDir()[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}
