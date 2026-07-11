package builtins

import (
	"bytes"
	"testing"

	"github.com/codecrafters-io/shell-starter-go/app/variables"
	"github.com/google/go-cmp/cmp"
)

func TestDeclare(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*variables.VariablesStore)
		args    []string
		wantOut string
		wantErr string
	}{
		{
			name:    "prints not found for -p",
			args:    []string{"-p", "missing_variable"},
			wantErr: "declare: missing_variable: not found\n",
		},
		{
			name: "prints description for -p",
			setup: func(store *variables.VariablesStore) {
				store.Set("foo", "bar")
			},
			args:    []string{"-p", "foo"},
			wantOut: `declare -- foo="bar"` + "\n",
		},
		{
			name: "stores assignment",
			args: []string{"foo=bar"},
		},
		{
			name: "overwrites existing variable",
			setup: func(store *variables.VariablesStore) {
				store.Set("foo", "bar")
			},
			args: []string{"foo=bar2"},
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
			store := variables.NewVariablesStore()
			if tt.setup != nil {
				tt.setup(store)
			}

			var stdout, stderr bytes.Buffer
			Declare(&stdout, &stderr, tt.args, store)

			if tt.args != nil && len(tt.args) == 1 {
				if name, value, ok := parseAssignment(tt.args[0]); ok {
					got, exists := store.Get(name)
					if !exists {
						t.Fatalf("Get(%q) exists = false, want true", name)
					}
					if diff := cmp.Diff(value, got); diff != "" {
						t.Errorf("Get(%q) value mismatch (-want +got):\n%s", name, diff)
					}
				}
			}

			if diff := cmp.Diff(tt.wantOut, stdout.String()); diff != "" {
				t.Errorf("Declare(%v) stdout mismatch (-want +got):\n%s", tt.args, diff)
			}
			if diff := cmp.Diff(tt.wantErr, stderr.String()); diff != "" {
				t.Errorf("Declare(%v) stderr mismatch (-want +got):\n%s", tt.args, diff)
			}
		})
	}
}
