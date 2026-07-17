package completer

import (
	"os"
	"runtime"
	"testing"

	"github.com/codecrafters-io/shell-starter-go/app/session"
	"github.com/codecrafters-io/shell-starter-go/app/testutils"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestCompleteFile(t *testing.T) {
	tests := []struct {
		name         string
		setup        func(t *testing.T, root string)
		buffer       string
		wantBuffer   string
		wantListings []string
	}{
		{
			name: "completes partial filename from cwd",
			setup: func(t *testing.T, root string) {
				for _, path := range []string{"readme.txt", "notes.md", "hello_world.py"} {
					testutils.CreatePath(t, root, path)
				}
			},
			buffer:     "cat re",
			wantBuffer: "cat readme.txt ",
		},
		{
			name: "ambiguous prefix lists matches",
			setup: func(t *testing.T, root string) {
				for _, path := range []string{"file.txt", "foo.txt", "fizz.txt"} {
					testutils.CreatePath(t, root, path)
				}
			},
			buffer:       "cat f",
			wantBuffer:   "cat f",
			wantListings: []string{"file.txt", "fizz.txt", "foo.txt"},
		},
		{
			name: "completes directory with trailing slash",
			setup: func(t *testing.T, root string) {
				testutils.CreatePath(t, root, "project/")
				testutils.CreatePath(t, root, "readme.txt")
			},
			buffer:     "cd proj",
			wantBuffer: "cd project/",
		},
		{
			name: "empty argument lists cwd files",
			setup: func(t *testing.T, root string) {
				for _, path := range []string{"alpha.txt", "beta.txt"} {
					testutils.CreatePath(t, root, path)
				}
			},
			buffer:       "cat ",
			wantBuffer:   "cat ",
			wantListings: []string{"alpha.txt", "beta.txt"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()
			tt.setup(t, root)
			t.Chdir(root)

			c := New(session.NewSession())
			gotBuffer, gotListings := c.completeFile(tt.buffer)
			if diff := cmp.Diff(tt.wantBuffer, gotBuffer); diff != "" {
				t.Errorf("completeFile(%q) buffer mismatch (-want +got):\n%s", tt.buffer, diff)
			}
			if diff := cmp.Diff(tt.wantListings, gotListings, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("completeFile(%q) listings mismatch (-want +got):\n%s", tt.buffer, diff)
			}
		})
	}
}

func TestFileCandidates(t *testing.T) {
	root := t.TempDir()
	for _, path := range []string{"one.txt", "two.txt", "nested/"} {
		testutils.CreatePath(t, root, path)
	}
	t.Chdir(root)

	c := New(session.NewSession())
	got := c.fileCandidates("cat ")
	want := []string{"nested/", "one.txt", "two.txt"}
	if diff := cmp.Diff(want, got, cmpopts.EquateEmpty()); diff != "" {
		t.Errorf("fileCandidates() mismatch (-want +got):\n%s", diff)
	}

	if got := c.fileCandidates("echo"); got != nil {
		t.Errorf("fileCandidates(no space) = %v, want nil", got)
	}
}

func TestFileCandidatesGetwdError(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("cannot remove cwd directory on Windows")
	}

	root := t.TempDir()
	t.Chdir(root)
	if err := os.Remove(root); err != nil {
		t.Fatalf("Remove() error = %v", err)
	}

	c := New(session.NewSession())
	if got := c.fileCandidates("cat "); got != nil {
		t.Errorf("fileCandidates() = %v, want nil when Getwd fails", got)
	}
}
