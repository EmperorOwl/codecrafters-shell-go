package variables

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestVariablesStore_SetGet(t *testing.T) {
	tests := []struct {
		name      string
		varName   string
		value     string
		overwrite string
		want      string
		wantOK    bool
	}{
		{
			name:    "stores and retrieves value",
			varName: "foo",
			value:   "bar",
			want:    "bar",
			wantOK:  true,
		},
		{
			name:      "overwrites existing value",
			varName:   "foo",
			value:     "bar",
			overwrite: "bar2",
			want:      "bar2",
			wantOK:    true,
		},
		{
			name:    "missing variable",
			varName: "missing",
			wantOK:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := NewVariablesStore()
			if tt.value != "" {
				store.Set(tt.varName, tt.value)
			}
			if tt.overwrite != "" {
				store.Set(tt.varName, tt.overwrite)
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
