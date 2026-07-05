package shell

import (
	"io"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestCompleteArgument(t *testing.T) {
	cwdFiles := []string{"hello_world.py", "notes.md", "readme.txt"}

	tests := []struct {
		name         string
		candidates   []string
		buffer       string
		wantBuffer   string
		wantListings []string
	}{
		{
			name:       "completes partial filename",
			candidates: cwdFiles,
			buffer:     "cat re",
			wantBuffer: "cat readme.txt ",
		},
		{
			name:       "completes hello prefix",
			candidates: cwdFiles,
			buffer:     "cat hello",
			wantBuffer: "cat hello_world.py ",
		},
		{
			name:       "no file match leaves buffer unchanged",
			candidates: cwdFiles,
			buffer:     "cat missing",
			wantBuffer: "cat missing",
		},
		{
			name:         "ambiguous prefix lists matches",
			candidates:   []string{"file.txt", "foo.txt", "fizz.txt"},
			buffer:       "cat f",
			wantBuffer:   "cat f",
			wantListings: []string{"file.txt", "fizz.txt", "foo.txt"},
		},
		{
			name:         "empty argument lists all files",
			candidates:   cwdFiles,
			buffer:       "cat ",
			wantBuffer:   "cat ",
			wantListings: cwdFiles,
		},
		{
			name:       "completes directory with trailing slash",
			candidates: []string{"project/", "readme.txt"},
			buffer:     "cd proj",
			wantBuffer: "cd project/",
		},
		{
			name:       "completes later argument filename",
			candidates: []string{"first.txt", "second.txt"},
			buffer:     "echo first.txt sec",
			wantBuffer: "echo first.txt second.txt ",
		},
		{
			name:       "programmable single candidate",
			candidates: []string{"run"},
			buffer:     "docker ",
			wantBuffer: "docker run ",
		},
		{
			name:         "programmable multiple candidates",
			candidates:   []string{"status", "stash"},
			buffer:       "git sta",
			wantBuffer:   "git sta",
			wantListings: []string{"stash", "status"},
		},
	}

	s := New(strings.NewReader(""), io.Discard, io.Discard)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBuffer, gotListings := s.completeArgument(tt.buffer, tt.candidates)
			if diff := cmp.Diff(tt.wantBuffer, gotBuffer); diff != "" {
				t.Errorf("completeArgument(%q) buffer mismatch (-want +got):\n%s", tt.buffer, diff)
			}
			if diff := cmp.Diff(tt.wantListings, gotListings, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("completeArgument(%q) listings mismatch (-want +got):\n%s", tt.buffer, diff)
			}
		})
	}
}

func TestCompleteCommand(t *testing.T) {
	s := New(strings.NewReader(""), io.Discard, io.Discard)

	tests := []struct {
		buffer     string
		wantBuffer string
	}{
		{buffer: "ech", wantBuffer: "echo "},
		{buffer: "exi", wantBuffer: "exit "},
	}

	for _, tt := range tests {
		t.Run(tt.buffer, func(t *testing.T) {
			gotBuffer, gotListings := s.completeCommand(tt.buffer)
			if diff := cmp.Diff(tt.wantBuffer, gotBuffer); diff != "" {
				t.Errorf("completeCommand(%q) buffer mismatch (-want +got):\n%s", tt.buffer, diff)
			}
			if diff := cmp.Diff([]string(nil), gotListings, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("completeCommand(%q) listings mismatch (-want +got):\n%s", tt.buffer, diff)
			}
		})
	}
}
