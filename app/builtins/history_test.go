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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := &history.HistoryList{}
			for _, command := range tt.commands {
				list.Add(command)
			}

			var out bytes.Buffer
			History(&out, list)

			if diff := cmp.Diff(tt.want, out.String()); diff != "" {
				t.Errorf("History() output mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
