package terminal

import (
	"os"
	"testing"
)

func TestSession_nonTerminalStdin(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Pipe() error = %v", err)
	}
	t.Cleanup(func() {
		r.Close()
		w.Close()
	})

	session := NewSession(r)
	if session.RawMode() {
		t.Error("RawMode() = true, want false for pipe stdin")
	}
	if session.PrepareRead() {
		t.Error("PrepareRead() = true, want false for pipe stdin")
	}
	if err := session.Close(); err != nil {
		t.Errorf("Close() error = %v", err)
	}
}

func TestSession_PrepareReadAfterClose(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Pipe() error = %v", err)
	}
	t.Cleanup(func() {
		r.Close()
		w.Close()
	})

	session := NewSession(r)
	if err := session.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
	if session.PrepareRead() {
		t.Error("PrepareRead() after Close() = true, want false")
	}
}
