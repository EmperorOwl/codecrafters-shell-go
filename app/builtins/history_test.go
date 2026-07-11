package builtins

import (
	"bytes"
	"strings"
	"testing"

	"github.com/codecrafters-io/shell-starter-go/app/history"
	"github.com/google/go-cmp/cmp"
)

func TestHistory(t *testing.T) {
	tests := []struct {
		name     string
		commands []string
		limit    int
		want     string
	}{
		{
			name: "empty history",
		},
		{
			name:     "lists previous commands",
			commands: []string{"previous_command_1", "previous_command_2", "history"},
			want: strings.Join([]string{
				"    1  previous_command_1",
				"    2  previous_command_2",
				"    3  history",
			}, "\n") + "\n",
		},
		{
			name:     "limits to last two commands",
			commands: []string{"echo hello", "echo world", "invalid_command", "history 2"},
			limit:    2,
			want: strings.Join([]string{
				"    3  invalid_command",
				"    4  history 2",
			}, "\n") + "\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := &history.List{}
			for _, command := range tt.commands {
				list.Add(command)
			}

			var out bytes.Buffer
			History(&out, list, tt.limit)

			if diff := cmp.Diff(tt.want, out.String()); diff != "" {
				t.Errorf("History() output mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
