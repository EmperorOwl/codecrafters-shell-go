package variables

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestStore_Set(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(*Store)
		wantStore map[string]string
	}{
		{
			name: "stores value",
			setup: func(s *Store) {
				s.Set("foo", "bar")
			},
			wantStore: map[string]string{"foo": "bar"},
		},
		{
			name: "overwrites existing value",
			setup: func(s *Store) {
				s.Set("foo", "bar")
				s.Set("foo", "bar2")
			},
			wantStore: map[string]string{"foo": "bar2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := NewStore()
			if tt.setup != nil {
				tt.setup(store)
			}

			if diff := cmp.Diff(tt.wantStore, store.Entries()); diff != "" {
				t.Errorf("store entries mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestStore_Get(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*Store)
		varName string
		want    string
		wantOK  bool
	}{
		{
			name: "returns stored value",
			setup: func(s *Store) {
				s.Set("foo", "bar")
			},
			varName: "foo",
			want:    "bar",
			wantOK:  true,
		},
		{
			name: "returns overwritten value",
			setup: func(s *Store) {
				s.Set("foo", "bar")
				s.Set("foo", "bar2")
			},
			varName: "foo",
			want:    "bar2",
			wantOK:  true,
		},
		{
			name:    "missing variable",
			varName: "missing",
			wantOK:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := NewStore()
			if tt.setup != nil {
				tt.setup(store)
			}

			got, ok := store.Get(tt.varName)
			if diff := cmp.Diff(tt.wantOK, ok); diff != "" {
				t.Errorf("Get(%q) ok mismatch (-want +got):\n%s", tt.varName, diff)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("Get(%q) value mismatch (-want +got):\n%s", tt.varName, diff)
			}
		})
	}
}
