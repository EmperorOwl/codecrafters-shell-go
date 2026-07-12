package builtins

import (
	"bytes"
	"testing"

	"github.com/codecrafters-io/shell-starter-go/app/history"
	"github.com/google/go-cmp/cmp"
)

func TestHistory(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(*history.List)
		args      []string
		wantLines []string
	}{
		{
			name: "empty history",
		},
		{
			name: "lists previous commands",
			setup: func(list *history.List) {
				list.Add("previous_command_1")
				list.Add("previous_command_2")
				list.Add("history")
			},
			wantLines: []string{
				"    1  previous_command_1",
				"    2  previous_command_2",
				"    3  history",
			},
		},
		{
			name: "limits to last two commands",
			setup: func(list *history.List) {
				list.Add("echo hello")
				list.Add("echo world")
				list.Add("invalid_command")
				list.Add("history 2")
			},
			args: []string{"2"},
			wantLines: []string{
				"    3  invalid_command",
				"    4  history 2",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := history.NewList()
			if tt.setup != nil {
				tt.setup(list)
			}

			var stdout, stderr bytes.Buffer
			historyBuiltin(&stdout, &stderr, tt.args, list)

			if diff := cmp.Diff(wantStdout(tt.wantLines), stdout.String()); diff != "" {
				t.Errorf("historyBuiltin(%v) stdout mismatch (-want +got):\n%s", tt.args, diff)
			}
		})
	}
}
