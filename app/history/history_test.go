package history

import (
	"bytes"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestHistoryList_AddAndList(t *testing.T) {
	tests := []struct {
		name     string
		commands []string
		want     []Entry
	}{
		{
			name: "empty history",
		},
		{
			name:     "single command",
			commands: []string{"echo hello"},
			want: []Entry{{
				Number:  1,
				Command: "echo hello",
			}},
		},
		{
			name:     "multiple commands",
			commands: []string{"previous_command_1", "previous_command_2", "history"},
			want: []Entry{
				{Number: 1, Command: "previous_command_1"},
				{Number: 2, Command: "previous_command_2"},
				{Number: 3, Command: "history"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var list HistoryList
			for _, command := range tt.commands {
				list.Add(command)
			}

			got := list.List()
			if diff := cmp.Diff(tt.want, got, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("List() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestHistoryList_ListLast(t *testing.T) {
	tests := []struct {
		name     string
		commands []string
		limit    int
		want     []Entry
	}{
		{
			name:  "empty history",
			limit: 2,
		},
		{
			name:     "limit greater than history length",
			commands: []string{"echo first"},
			limit:    2,
			want: []Entry{{
				Number:  1,
				Command: "echo first",
			}},
		},
		{
			name:     "shows last two commands with original numbers",
			commands: []string{"echo first", "echo second", "history 2"},
			limit:    2,
			want: []Entry{
				{Number: 2, Command: "echo second"},
				{Number: 3, Command: "history 2"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var list HistoryList
			for _, command := range tt.commands {
				list.Add(command)
			}

			got := list.ListLast(tt.limit)
			if diff := cmp.Diff(tt.want, got, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("ListLast(%d) mismatch (-want +got):\n%s", tt.limit, diff)
			}
		})
	}
}

func TestWriteAll(t *testing.T) {
	tests := []struct {
		name     string
		commands []string
		want     string
	}{
		{
			name: "empty history",
		},
		{
			name:     "numbered commands",
			commands: []string{"previous_command_1", "previous_command_2", "history"},
			want: strings.Join([]string{
				"    1  previous_command_1",
				"    2  previous_command_2",
				"    3  history",
			}, "\n") + "\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var list HistoryList
			for _, command := range tt.commands {
				list.Add(command)
			}

			var out bytes.Buffer
			WriteAll(&out, list.List())

			if diff := cmp.Diff(tt.want, out.String()); diff != "" {
				t.Errorf("WriteAll() output mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
