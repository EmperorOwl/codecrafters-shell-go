package completion

import (
	"testing"

	"github.com/codecrafters-io/shell-starter-go/app/builtins"
)

func TestParseCompletionContext(t *testing.T) {
	tests := []struct {
		name         string
		buffer       string
		wantCommand  string
		wantCurrent  string
		wantPrevious string
	}{
		{
			name:         "first argument completion",
			buffer:       "git ",
			wantCommand:  "git",
			wantCurrent:  "",
			wantPrevious: "",
		},
		{
			name:         "partial first argument",
			buffer:       "git remot",
			wantCommand:  "git",
			wantCurrent:  "remot",
			wantPrevious: "",
		},
		{
			name:         "later argument completion",
			buffer:       "git remote set",
			wantCommand:  "git",
			wantCurrent:  "set",
			wantPrevious: "remote",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			command, currentWord, previousWord := parseCompletionContext(tt.buffer)
			if command != tt.wantCommand {
				t.Errorf("command = %q, want %q", command, tt.wantCommand)
			}
			if currentWord != tt.wantCurrent {
				t.Errorf("currentWord = %q, want %q", currentWord, tt.wantCurrent)
			}
			if previousWord != tt.wantPrevious {
				t.Errorf("previousWord = %q, want %q", previousWord, tt.wantPrevious)
			}
		})
	}
}

func TestApplyTabProgrammableTab(t *testing.T) {
	registeredCompleters := map[string]builtins.Completer{
		"git": {
			Path: "/path/to/completer",
			Func: func(scriptPath, command, currentWord, previousWord string) ([]string, error) {
				if scriptPath != "/path/to/completer" || command != "git" || currentWord != "set" || previousWord != "remote" {
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
