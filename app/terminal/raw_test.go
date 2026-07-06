package terminal

import (
	"os"
	"testing"
)

func TestRawMode_nonTerminalStdin(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Pipe() error = %v", err)
	}
	t.Cleanup(func() {
		r.Close()
		w.Close()
	})

	raw := NewRawMode(r)
	if raw.Active() {
		t.Error("Active() = true, want false for pipe stdin")
	}
	if raw.PrepareRead() {
		t.Error("PrepareRead() = true, want false for pipe stdin")
	}
	if err := raw.Close(); err != nil {
		t.Errorf("Close() error = %v", err)
	}
}

func TestRawMode_PrepareReadAfterClose(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Pipe() error = %v", err)
	}
	t.Cleanup(func() {
		r.Close()
		w.Close()
	})

	raw := NewRawMode(r)
	if err := raw.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
	if raw.PrepareRead() {
		t.Error("PrepareRead() after Close() = true, want false")
	}
}
