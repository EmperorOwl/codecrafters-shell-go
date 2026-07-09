package files

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
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
			if diff := cmp.Diff(tt.want, got, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("ListInDir(%q, %q) mismatch (-want +got):\n%s", root, tt.dir, diff)
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

func TestReadLines(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    []string
	}{
		{
			name:    "multiple lines",
			content: "echo hello\necho world\n",
			want:    []string{"echo hello", "echo world"},
		},
		{
			name:    "includes empty line",
			content: "echo hello\n\n",
			want:    []string{"echo hello", ""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "lines.txt")
			if err := os.WriteFile(path, []byte(tt.content), 0o644); err != nil {
				t.Fatalf("WriteFile() error = %v", err)
			}

			got, err := ReadLines(path)
			if err != nil {
				t.Fatalf("ReadLines() error = %v", err)
			}
			if diff := cmp.Diff(tt.want, got, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("ReadLines() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestWriteLines(t *testing.T) {
	tests := []struct {
		name  string
		lines []string
		want  string
	}{
		{
			name: "writes lines with trailing newline",
			lines: []string{
				"echo hello",
				"echo world",
				"history -w /tmp/hist",
			},
			want: "echo hello\necho world\nhistory -w /tmp/hist\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "lines.txt")
			if err := WriteLines(path, tt.lines); err != nil {
				t.Fatalf("WriteLines() error = %v", err)
			}

			got, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("ReadFile() error = %v", err)
			}
			if diff := cmp.Diff(tt.want, string(got)); diff != "" {
				t.Errorf("WriteLines() content mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestAppendLines(t *testing.T) {
	path := filepath.Join(t.TempDir(), "lines.txt")
	initial := "echo initial_command_1\necho initial_command_2\n\n"
	if err := os.WriteFile(path, []byte(initial), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	if err := AppendLines(path, []string{"echo new_command", "history -a " + path}); err != nil {
		t.Fatalf("AppendLines() error = %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	want := strings.Join([]string{
		"echo initial_command_1",
		"echo initial_command_2",
		"echo new_command",
		"history -a " + path,
	}, "\n") + "\n"
	if diff := cmp.Diff(want, string(got)); diff != "" {
		t.Errorf("AppendLines() content mismatch (-want +got):\n%s", diff)
	}
}
