package completion

import (
	"testing"

	"github.com/codecrafters-io/shell-starter-go/app/builtins"
)

func TestBuildCompleterFuncOptions(t *testing.T) {
	tests := []struct {
		name   string
		buffer string
		want   builtins.CompleterFuncOptions
	}{
		{
			name:   "first argument completion",
			buffer: "git ",
			want: builtins.CompleterFuncOptions{
				Command:   "git",
				CompLine:  "git ",
				CompPoint: 4,
			},
		},
		{
			name:   "partial first argument",
			buffer: "git remot",
			want: builtins.CompleterFuncOptions{
				Command:     "git",
				CurrentWord: "remot",
				CompLine:    "git remot",
				CompPoint:   9,
			},
		},
		{
			name:   "later argument completion",
			buffer: "git remote set",
			want: builtins.CompleterFuncOptions{
				Command:      "git",
				CurrentWord:  "set",
				PreviousWord: "remote",
				CompLine:     "git remote set",
				CompPoint:    14,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildCompleterFuncOptions(tt.buffer)
			if got.Command != tt.want.Command {
				t.Errorf("Command = %q, want %q", got.Command, tt.want.Command)
			}
			if got.CurrentWord != tt.want.CurrentWord {
				t.Errorf("CurrentWord = %q, want %q", got.CurrentWord, tt.want.CurrentWord)
			}
			if got.PreviousWord != tt.want.PreviousWord {
				t.Errorf("PreviousWord = %q, want %q", got.PreviousWord, tt.want.PreviousWord)
			}
			if got.CompLine != tt.want.CompLine {
				t.Errorf("CompLine = %q, want %q", got.CompLine, tt.want.CompLine)
			}
			if got.CompPoint != tt.want.CompPoint {
				t.Errorf("CompPoint = %d, want %d", got.CompPoint, tt.want.CompPoint)
			}
		})
	}
}

func TestApplyTabProgrammableTab(t *testing.T) {
	registeredCompleters := map[string]builtins.Completer{
		"git": {
			Path: "/path/to/completer",
			Func: func(opts builtins.CompleterFuncOptions) ([]string, error) {
				if opts.ScriptPath != "/path/to/completer" || opts.Command != "git" || opts.CurrentWord != "set" || opts.PreviousWord != "remote" {
					return nil, nil
				}
				if opts.CompLine != "git remote set" || opts.CompPoint != len("git remote set") {
					return nil, nil
				}
				return []string{"set-url"}, nil
			},
		},
	}

	gotBuffer, gotListings := ApplyTab(nil, nil, nil, registeredCompleters, "git remote set")
	wantBuffer := "git remote set-url "
	if gotBuffer != wantBuffer {
		t.Errorf("ApplyTab(%q) buffer = %q, want %q", "git remote set", gotBuffer, wantBuffer)
	}
	if len(gotListings) != 0 {
		t.Errorf("ApplyTab(%q) listings = %v, want nil", "git remote set", gotListings)
	}
}
