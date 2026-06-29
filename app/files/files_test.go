package files

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestListInDir(t *testing.T) {
	tests := []struct {
		name        string
		createPaths []string
		dir         string
		want        []string
	}{
		{
			name:        "current directory",
			createPaths: []string{"notes.md", "readme.txt", "subdir/", "path/to/file.txt"},
			dir:         "",
			want:        []string{"notes.md", "path/", "readme.txt", "subdir/"},
		},
		{
			name:        "nested directory",
			createPaths: []string{"path/to/file.txt"},
			dir:         "path/to/",
			want:        []string{"file.txt"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()

			for _, path := range tt.createPaths {
				if err := createPath(root, path); err != nil {
					t.Fatalf("createPath(%q) error = %v", path, err)
				}
			}

			got := ListInDir(root, tt.dir)
			if len(got) != len(tt.want) {
				t.Fatalf("ListInDir(%q, %q) = %v, want %v", root, tt.dir, got, tt.want)
			}
			for i := range tt.want {
				if got[i] != tt.want[i] {
					t.Errorf("ListInDir(%q, %q)[%d] = %q, want %q", root, tt.dir, i, got[i], tt.want[i])
				}
			}
		})
	}
}

func createPath(root, rel string) error {
	full := filepath.Join(root, filepath.FromSlash(rel))
	if strings.HasSuffix(rel, "/") {
		return os.MkdirAll(strings.TrimSuffix(full, string(os.PathSeparator)), 0755)
	}

	if err := os.MkdirAll(filepath.Dir(full), 0755); err != nil {
		return err
	}
	return os.WriteFile(full, nil, 0644)
}
