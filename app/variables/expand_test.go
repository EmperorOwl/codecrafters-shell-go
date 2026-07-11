package variables

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestExpandField(t *testing.T) {
	store := NewVariablesStore()
	store.Set("Variable_1", "Value_1")
	store.Set("Variable_2", "Value2")
	store.Set("Var1", "foo")
	store.Set("Var2", "bar")
	store.Set("Item", "widget")
	store.Set("Foo1", "Bar2")

	tests := []struct {
		name  string
		field string
		want  string
	}{
		{name: "no expansion", field: "echo", want: "echo"},
		{name: "expands variable", field: "$Variable_1", want: "Value_1"},
		{name: "expands second variable", field: "$Variable_2", want: "Value2"},
		{name: "expands multiple variables", field: "$Variable_1$Variable_2", want: "Value_1Value2"},
		{name: "expands braced variable with suffix", field: "${Var1}end", want: "fooend"},
		{name: "expands multiple braced variables", field: "${Var1}and${Var2}", want: "fooandbar"},
		{name: "expands braced variable within word", field: "stock_${Item}_id", want: "stock_widget_id"},
		{name: "expands braced variable", field: "${Foo1}", want: "Bar2"},
		{name: "leaves undefined variable literal", field: "$HOME", want: "$HOME"},
		{name: "leaves undefined braced variable literal", field: "${HOME}", want: "${HOME}"},
		{name: "leaves bare dollar sign", field: "a$", want: "a$"},
		{name: "leaves invalid identifier", field: "$1foo", want: "$1foo"},
		{name: "leaves invalid braced identifier", field: "${1foo}", want: "${1foo}"},
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
	store.Set("Item", "widget")
	store.Set("Foo1", "Bar2")

	tests := []struct {
		name   string
		fields []string
		want   []string
	}{
		{
			name:   "expands simple variables",
			fields: []string{"custom_exe", "$Variable_1", "$Variable_2"},
			want:   []string{"custom_exe", "Value_1", "Value2"},
		},
		{
			name:   "expands braced variables within words",
			fields: []string{"custom_exe_1234", "stock_${Item}_id", "${Foo1}"},
			want:   []string{"custom_exe_1234", "stock_widget_id", "Bar2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExpandFields(store, tt.fields)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("ExpandFields() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
