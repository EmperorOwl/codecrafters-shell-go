package shell

import (
	"testing"

	"github.com/codecrafters-io/shell-starter-go/app/parser"
	"github.com/codecrafters-io/shell-starter-go/app/variables"
	"github.com/google/go-cmp/cmp"
)

func TestExpandParsedLine(t *testing.T) {
	tests := []struct {
		name  string
		setup func() *variables.Store
		in    parser.Line
		want  parser.Line
	}{
		{
			name: "expands command arguments",
			setup: func() *variables.Store {
				store := variables.NewStore()
				store.Set("HOME", "/home/user")
				return store
			},
			in: parser.Line{
				Commands: [][]string{{"echo", "$HOME"}},
			},
			want: parser.Line{
				Commands: [][]string{{"echo", "/home/user"}},
			},
		},
		{
			name: "expands redirect paths",
			setup: func() *variables.Store {
				store := variables.NewStore()
				store.Set("HOME", "/home/user")
				store.Set("NAME", "alice")
				return store
			},
			in: parser.Line{
				Commands: [][]string{{"echo", "hi"}},
				Redirect: parser.Redirect{
					StdoutPath: "$HOME/out.txt",
					StderrPath: "${NAME}.err",
				},
			},
			want: parser.Line{
				Commands: [][]string{{"echo", "hi"}},
				Redirect: parser.Redirect{
					StdoutPath: "/home/user/out.txt",
					StderrPath: "alice.err",
				},
			},
		},
		{
			name: "expands each pipeline segment",
			setup: func() *variables.Store {
				store := variables.NewStore()
				store.Set("HOME", "/home/user")
				store.Set("NAME", "alice")
				return store
			},
			in: parser.Line{
				Pipeline: true,
				Commands: [][]string{
					{"echo", "$NAME"},
					{"echo", "$HOME"},
				},
			},
			want: parser.Line{
				Pipeline: true,
				Commands: [][]string{
					{"echo", "alice"},
					{"echo", "/home/user"},
				},
			},
		},
		{
			name:  "nil store leaves line unchanged",
			setup: func() *variables.Store { return nil },
			in: parser.Line{
				Commands: [][]string{{"echo", "$HOME"}},
				Redirect: parser.Redirect{StdoutPath: "$HOME/out"},
			},
			want: parser.Line{
				Commands: [][]string{{"echo", "$HOME"}},
				Redirect: parser.Redirect{StdoutPath: "$HOME/out"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandParsedLine(tt.in, tt.setup())
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("expandParsedLine() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
