package variables

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestExpandField(t *testing.T) {
	tests := []struct {
		name  string
		setup func(*Store)
		field string
		want  string
	}{
		{name: "no expansion", field: "echo", want: "echo"},
		{
			name: "expands variable",
			setup: func(store *Store) {
				store.Set("Variable_1", "Value_1")
			},
			field: "$Variable_1",
			want:  "Value_1",
		},
		{
			name: "expands second variable",
			setup: func(store *Store) {
				store.Set("Variable_2", "Value2")
			},
			field: "$Variable_2",
			want:  "Value2",
		},
		{
			name: "expands multiple variables",
			setup: func(store *Store) {
				store.Set("Variable_1", "Value_1")
				store.Set("Variable_2", "Value2")
			},
			field: "$Variable_1$Variable_2",
			want:  "Value_1Value2",
		},
		{
			name: "expands braced variable with suffix",
			setup: func(store *Store) {
				store.Set("Var1", "foo")
			},
			field: "${Var1}end",
			want:  "fooend",
		},
		{
			name: "expands multiple braced variables",
			setup: func(store *Store) {
				store.Set("Var1", "foo")
				store.Set("Var2", "bar")
			},
			field: "${Var1}and${Var2}",
			want:  "fooandbar",
		},
		{
			name: "expands braced variable within word",
			setup: func(store *Store) {
				store.Set("Item", "widget")
			},
			field: "stock_${Item}_id",
			want:  "stock_widget_id",
		},
		{
			name: "expands braced variable",
			setup: func(store *Store) {
				store.Set("Foo1", "Bar2")
			},
			field: "${Foo1}",
			want:  "Bar2",
		},
		{name: "expands unset variable to empty", field: "$HOME", want: ""},
		{name: "expands unset braced variable to empty", field: "${HOME}", want: ""},
		{name: "keeps suffix when braced variable is unset", field: "${missing}world", want: "world"},
		{
			name: "expands set variable with suffix",
			setup: func(store *Store) {
				store.Set("greeting", "hello")
			},
			field: "${greeting}world",
			want:  "helloworld",
		},
		{name: "leaves bare dollar sign", field: "a$", want: "a$"},
		{name: "leaves invalid identifier", field: "$1foo", want: "$1foo"},
		{name: "leaves invalid braced identifier", field: "${1foo}", want: "${1foo}"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := NewStore()
			if tt.setup != nil {
				tt.setup(store)
			}

			got := ExpandField(store, tt.field)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("ExpandField(%q) mismatch (-want +got):\n%s", tt.field, diff)
			}
		})
	}
}

func TestExpandFields(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(*Store)
		fields []string
		want   []string
	}{
		{
			name: "expands simple variables",
			setup: func(store *Store) {
				store.Set("Variable_1", "Value_1")
				store.Set("Variable_2", "Value2")
			},
			fields: []string{"custom_exe", "$Variable_1", "$Variable_2"},
			want:   []string{"custom_exe", "Value_1", "Value2"},
		},
		{
			name: "expands braced variables within words",
			setup: func(store *Store) {
				store.Set("Item", "widget")
				store.Set("Foo1", "Bar2")
			},
			fields: []string{"custom_exe_1234", "stock_${Item}_id", "${Foo1}"},
			want:   []string{"custom_exe_1234", "stock_widget_id", "Bar2"},
		},
		{
			name: "drops empty args after unset expansion",
			setup: func(store *Store) {
				store.Set("existing", "existingsvalue")
			},
			fields: []string{"custom_exe_1234", "${missing1}end", "${existing}", "${missing2}"},
			want:   []string{"custom_exe_1234", "end", "existingsvalue"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := NewStore()
			if tt.setup != nil {
				tt.setup(store)
			}

			got := ExpandFields(store, tt.fields)
			if diff := cmp.Diff(tt.want, got, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("ExpandFields() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
