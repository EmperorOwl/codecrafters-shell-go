package completer

import (
	"testing"

	_ "github.com/codecrafters-io/shell-starter-go/app/builtins"
	"github.com/codecrafters-io/shell-starter-go/app/session"
	"github.com/codecrafters-io/shell-starter-go/app/terminal"
	"github.com/codecrafters-io/shell-starter-go/app/testutils"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestCompleteBufferRouting(t *testing.T) {
	tests := []struct {
		name         string
		setup        func(t *testing.T) *Completer
		buffer       string
		wantBuffer   string
		wantListings []string
	}{
		{
			name: "routes to command completion without space",
			setup: func(t *testing.T) *Completer {
				return New(session.NewSession())
			},
			buffer:     "ech",
			wantBuffer: "echo ",
		},
		{
			name: "routes to file completion without programmable script",
			setup: func(t *testing.T) *Completer {
				root := t.TempDir()
				for _, path := range []string{"readme.txt", "notes.md"} {
					testutils.CreatePath(t, root, path)
				}
				t.Chdir(root)
				return New(session.NewSession())
			},
			buffer:     "cat re",
			wantBuffer: "cat readme.txt ",
		},
		{
			name: "routes to programmable completion when registered",
			setup: func(t *testing.T) *Completer {
				dir := t.TempDir()
				scriptPath := testutils.WriteCompleterScript(t, dir, "stash", "status")
				sess := session.NewSession()
				sess.Completion.Register("git", scriptPath)
				return New(sess)
			},
			buffer:       "git sta",
			wantBuffer:   "git sta",
			wantListings: []string{"stash", "status"},
		},
		{
			name: "programmable completion unique match extends buffer",
			setup: func(t *testing.T) *Completer {
				dir := t.TempDir()
				scriptPath := testutils.WriteCompleterScript(t, dir, "run")
				sess := session.NewSession()
				sess.Completion.Register("docker", scriptPath)
				return New(sess)
			},
			buffer:     "docker ru",
			wantBuffer: "docker run ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.setup(t)
			gotBuffer, gotListings := c.completeBuffer(tt.buffer)
			if diff := cmp.Diff(tt.wantBuffer, gotBuffer); diff != "" {
				t.Errorf("completeBuffer(%q) buffer mismatch (-want +got):\n%s", tt.buffer, diff)
			}
			if diff := cmp.Diff(tt.wantListings, gotListings, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("completeBuffer(%q) listings mismatch (-want +got):\n%s", tt.buffer, diff)
			}
		})
	}
}

func TestHandleTab(t *testing.T) {
	tests := []struct {
		name  string
		setup func(t *testing.T) (*Completer, *terminal.TabState)
		input string
		want  terminal.TabResult
	}{
		{
			name: "completes builtin command",
			setup: func(t *testing.T) (*Completer, *terminal.TabState) {
				return New(session.NewSession()), &terminal.TabState{}
			},
			input: "ech",
			want:  terminal.TabResult{Buffer: "echo "},
		},
		{
			name: "rings bell on ambiguous command prefix",
			setup: func(t *testing.T) (*Completer, *terminal.TabState) {
				return New(session.NewSession()), &terminal.TabState{}
			},
			input: "e",
			want:  terminal.TabResult{Buffer: "e", RingBell: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, state := tt.setup(t)
			got := c.HandleTab(state, tt.input)
			if diff := cmp.Diff(tt.want, got, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("HandleTab() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
