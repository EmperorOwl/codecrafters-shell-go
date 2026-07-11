package variables

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestExpandField(t *testing.T) {
	store := NewVariablesStore()
	store.Set("Variable_1", "Value_1")
	store.Set("Variable_2", "Value2")

	tests := []struct {
		name  string
		field string
		want  string
	}{
		{name: "no expansion", field: "echo", want: "echo"},
		{name: "expands variable", field: "$Variable_1", want: "Value_1"},
		{name: "expands second variable", field: "$Variable_2", want: "Value2"},
		{name: "expands multiple variables", field: "$Variable_1$Variable_2", want: "Value_1Value2"},
		{name: "leaves undefined variable literal", field: "$HOME", want: "$HOME"},
		{name: "leaves bare dollar sign", field: "a$", want: "a$"},
		{name: "leaves invalid identifier", field: "$1foo", want: "$1foo"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExpandField(store, tt.field)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("ExpandField(%q) mismatch (-want +got):\n%s", tt.field, diff)
			}
		})
	}
}

func TestExpandFields(t *testing.T) {
	store := NewVariablesStore()
	store.Set("Variable_1", "Value_1")
	store.Set("Variable_2", "Value2")

	got := ExpandFields(store, []string{"custom_exe", "$Variable_1", "$Variable_2"})
	want := []string{"custom_exe", "Value_1", "Value2"}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("ExpandFields() mismatch (-want +got):\n%s", diff)
	}
}
