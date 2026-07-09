package repl

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/codecrafters-io/shell-starter-go/app/history"
	"github.com/google/go-cmp/cmp"
)

func TestNewStateLoadsHistfile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "histfile")
	content := "echo hello\necho world\n\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	t.Setenv("HISTFILE", path)

	state := NewState()
	state.History.Add("history")

	got := state.History.List()
	want := []history.Entry{
		{Number: 1, Command: "echo hello"},
		{Number: 2, Command: "echo world"},
		{Number: 3, Command: "history"},
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("NewState() history mismatch (-want +got):\n%s", diff)
	}
	if state.Histfile != path {
		t.Errorf("NewState() Histfile = %q, want %q", state.Histfile, path)
	}
}
