package terminal

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRawMode(t *testing.T) {
	tests := []struct {
		name            string
		setup           func(t *testing.T) (*RawMode, func())
		wantActive      bool
		wantPrepareRead bool
		closeBeforeRead bool
	}{
		{
			name: "non-terminal stdin is inactive",
			setup: func(t *testing.T) (*RawMode, func()) {
				t.Helper()
				r, w, err := os.Pipe()
				if err != nil {
					t.Fatalf("Pipe() error = %v", err)
				}
				cleanup := func() {
					r.Close()
					w.Close()
				}
				return NewRawMode(r), cleanup
			},
			wantActive:      false,
			wantPrepareRead: false,
		},
		{
			name: "prepare read after close is false",
			setup: func(t *testing.T) (*RawMode, func()) {
				t.Helper()
				r, w, err := os.Pipe()
				if err != nil {
					t.Fatalf("Pipe() error = %v", err)
				}
				cleanup := func() {
					r.Close()
					w.Close()
				}
				raw := NewRawMode(r)
				if err := raw.Close(); err != nil {
					t.Fatalf("Close() error = %v", err)
				}
				return raw, cleanup
			},
			closeBeforeRead: true,
			wantPrepareRead: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			raw, cleanup := tt.setup(t)
			t.Cleanup(cleanup)

			if diff := cmp.Diff(tt.wantActive, raw.Active()); diff != "" {
				t.Errorf("Active() mismatch (-want +got):\n%s", diff)
			}

			if !tt.closeBeforeRead {
				if diff := cmp.Diff(tt.wantPrepareRead, raw.PrepareRead()); diff != "" {
					t.Errorf("PrepareRead() mismatch (-want +got):\n%s", diff)
				}
				if err := raw.Close(); err != nil {
					t.Errorf("Close() error = %v", err)
				}
				return
			}

			if diff := cmp.Diff(tt.wantPrepareRead, raw.PrepareRead()); diff != "" {
				t.Errorf("PrepareRead() after Close() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
