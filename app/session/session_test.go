package session

import "testing"

func TestNewSession(t *testing.T) {
	sess := NewSession()

	if sess.Jobs == nil {
		t.Fatal("NewSession() Jobs = nil, want non-nil")
	}
	if sess.History == nil {
		t.Fatal("NewSession() History = nil, want non-nil")
	}
	if sess.Completion == nil {
		t.Fatal("NewSession() Completion = nil, want non-nil")
	}
	if sess.Variables == nil {
		t.Fatal("NewSession() Variables = nil, want non-nil")
	}
	if sess.Histfile != "" {
		t.Errorf("NewSession() Histfile = %q, want empty", sess.Histfile)
	}
}
