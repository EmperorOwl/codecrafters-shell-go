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
				Command:      "git",
				CurrentWord:  "remot",
				PreviousWord: "git",
				CompLine:     "git remot",
				CompPoint:    9,
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
	tests := []struct {
		name         string
		buffer       string
		candidates   []string
		wantBuffer   string
		wantListings []string
	}{
		{
			name:       "completes single candidate",
			buffer:     "git remote set",
			candidates: []string{"set-url"},
			wantBuffer: "git remote set-url ",
		},
		{
			name:         "lists multiple candidates",
			buffer:       "git sta",
			candidates:   []string{"status", "stash"},
			wantBuffer:   "git sta",
			wantListings: []string{"stash", "status"},
		},
		{
			name:       "completes to longest common prefix",
			buffer:     "git c",
			candidates: []string{"checkout", "cherry-pick"},
			wantBuffer: "git che",
		},
		{
			name:       "completes single remaining match",
			buffer:     "git chec",
			candidates: []string{"checkout", "cherry-pick"},
			wantBuffer: "git checkout ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			candidates := tt.candidates
			registeredCompleters := map[string]builtins.Completer{
				"git": {
					Func: func(builtins.CompleterFuncOptions) ([]string, error) {
						return candidates, nil
					},
				},
			}
			gotBuffer, gotListings := ApplyTab(nil, nil, nil, registeredCompleters, tt.buffer)
			if gotBuffer != tt.wantBuffer {
				t.Errorf("ApplyTab(%q) buffer = %q, want %q", tt.buffer, gotBuffer, tt.wantBuffer)
			}
			if len(gotListings) != len(tt.wantListings) {
				t.Fatalf("ApplyTab(%q) listings = %v, want %v", tt.buffer, gotListings, tt.wantListings)
			}
			for i, listing := range gotListings {
				if listing != tt.wantListings[i] {
					t.Errorf("listings[%d] = %q, want %q", i, listing, tt.wantListings[i])
				}
			}
		})
	}
}
