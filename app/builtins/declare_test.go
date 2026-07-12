package builtins

import (
	"bytes"
	"testing"

	"github.com/codecrafters-io/shell-starter-go/app/variables"
	"github.com/google/go-cmp/cmp"
)

func TestDeclare(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(*variables.Store)
		args       []string
		wantOut    string
		wantErr    string
		wantStore map[string]string
	}{
		{
			name:    "prints not found for -p",
			args:    []string{"-p", "missing_variable"},
			wantErr: "declare: missing_variable: not found\n",
		},
		{
			name: "prints description for -p",
			setup: func(store *variables.Store) {
				store.Set("foo", "bar")
			},
			args:       []string{"-p", "foo"},
			wantOut:    `declare -- foo="bar"` + "\n",
			wantStore: map[string]string{"foo": "bar"},
		},
		{
			name:       "stores assignment",
			args:       []string{"foo=bar"},
			wantStore: map[string]string{"foo": "bar"},
		},
		{
			name:       "stores underscore assignment",
			args:       []string{"_FOO=bar"},
			wantStore: map[string]string{"_FOO": "bar"},
		},
		{
			name:    "rejects digit at start",
			args:    []string{"67=x"},
			wantErr: "declare: `67=x': not a valid identifier\n",
		},
		{
			name: "overwrites existing variable",
			setup: func(store *variables.Store) {
				store.Set("foo", "bar")
			},
			args:       []string{"foo=bar2"},
			wantStore: map[string]string{"foo": "bar2"},
		},
		{
			name: "ignores bare declare",
			args: nil,
		},
		{
			name: "ignores -p without variable name",
			args: []string{"-p"},
		},
		{
			name: "ignores assignment without name",
			args: []string{"=value"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := variables.NewStore()
			if tt.setup != nil {
				tt.setup(store)
			}

			var stdout, stderr bytes.Buffer
			declareBuiltin(&stdout, &stderr, tt.args, store)

			if tt.wantStore != nil {
				if diff := cmp.Diff(tt.wantStore, store.Entries()); diff != "" {
					t.Errorf("store entries mismatch (-want +got):\n%s", diff)
				}
			}

			if diff := cmp.Diff(tt.wantOut, stdout.String()); diff != "" {
				t.Errorf("declareBuiltin(%v) stdout mismatch (-want +got):\n%s", tt.args, diff)
			}
			if diff := cmp.Diff(tt.wantErr, stderr.String()); diff != "" {
				t.Errorf("declareBuiltin(%v) stderr mismatch (-want +got):\n%s", tt.args, diff)
			}
		})
	}
}
