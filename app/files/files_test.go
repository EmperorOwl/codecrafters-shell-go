package files

import (
	"os"
	"path/filepath"
	"testing"
)

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
