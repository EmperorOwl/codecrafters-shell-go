package completion

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestApplyFileTab(t *testing.T) {
	cwdFiles := []string{"hello_world.py", "notes.md", "readme.txt"}

	tests := []struct {
		name         string
		fileDirs     map[string][]string
		buffer       string
		wantBuffer   string
		wantListings []string
	}{
		{
			name:       "completes partial filename",
			fileDirs:   map[string][]string{"": cwdFiles},
			buffer:     "cat re",
			wantBuffer: "cat readme.txt ",
		},
		{
			name:       "completes hello prefix",
			fileDirs:   map[string][]string{"": cwdFiles},
			buffer:     "cat hello",
			wantBuffer: "cat hello_world.py ",
		},
		{
			name:       "completes for unknown command",
			fileDirs:   map[string][]string{"": cwdFiles},
			buffer:     "xyz read",
			wantBuffer: "xyz readme.txt ",
		},
		{
			name:       "no file match leaves buffer unchanged",
			fileDirs:   map[string][]string{"": cwdFiles},
			buffer:     "cat missing",
			wantBuffer: "cat missing",
		},
		{
			name:         "ambiguous prefix lists matches",
			fileDirs:     map[string][]string{"": {"file.txt", "foo.txt", "fizz.txt"}},
			buffer:       "cat f",
			wantBuffer:   "cat f",
			wantListings: []string{"file.txt", "fizz.txt", "foo.txt"},
		},
		{
			name:         "empty argument lists all files in current directory",
			fileDirs:     map[string][]string{"": cwdFiles},
			buffer:       "cat ",
			wantBuffer:   "cat ",
			wantListings: cwdFiles,
		},
		{
			name:       "completes to longest common prefix",
			fileDirs:   map[string][]string{"": {"xyz_foo.txt", "xyz_foo_bar.txt", "xyz_foo_bar_baz.txt"}},
			buffer:     "cat xyz_",
			wantBuffer: "cat xyz_foo",
		},
		{
			name:       "completes to next longest common prefix",
			fileDirs:   map[string][]string{"": {"xyz_foo.txt", "xyz_foo_bar.txt", "xyz_foo_bar_baz.txt"}},
			buffer:     "cat xyz_foo_",
			wantBuffer: "cat xyz_foo_bar",
		},
		{
			name:       "completes final unique file",
			fileDirs:   map[string][]string{"": {"xyz_foo.txt", "xyz_foo_bar.txt", "xyz_foo_bar_baz.txt"}},
			buffer:     "cat xyz_foo_bar_",
			wantBuffer: "cat xyz_foo_bar_baz.txt ",
		},
		{
			name:       "completes nested path",
			fileDirs:   map[string][]string{"path/to/": {"file.txt", "other.txt"}},
			buffer:     "cat path/to/f",
			wantBuffer: "cat path/to/file.txt ",
		},
		{
			name:       "no nested match leaves buffer unchanged",
			fileDirs:   map[string][]string{"path/to/": {"file.txt", "other.txt"}},
			buffer:     "cat path/to/missing",
			wantBuffer: "cat path/to/missing",
		},
		{
			name:       "completes directory with trailing slash",
			fileDirs:   map[string][]string{"": {"project/", "readme.txt"}},
			buffer:     "cd proj",
			wantBuffer: "cd project/",
		},
		{
			name:       "completes lone directory after command",
			fileDirs:   map[string][]string{"": {"pig/"}},
			buffer:     "ls ",
			wantBuffer: "ls pig/",
		},
		{
			name:       "completes nested directory",
			fileDirs:   map[string][]string{"pig/": {"dog/"}},
			buffer:     "ls pig/",
			wantBuffer: "ls pig/dog/",
		},
		{
			name:       "completes later argument filename",
			fileDirs:   map[string][]string{"": {"first.txt", "second.txt"}},
			buffer:     "echo first.txt sec",
			wantBuffer: "echo first.txt second.txt ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			listFiles := func(dir string) []string {
				if tt.fileDirs == nil {
					return nil
				}
				return tt.fileDirs[dir]
			}

			gotBuffer, gotListings := applyFileTab(listFiles, tt.buffer)
			if diff := cmp.Diff(tt.wantBuffer, gotBuffer); diff != "" {
				t.Errorf("applyFileTab(%q) buffer mismatch (-want +got):\n%s", tt.buffer, diff)
			}
			if diff := cmp.Diff(tt.wantListings, gotListings, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("applyFileTab(%q) listings mismatch (-want +got):\n%s", tt.buffer, diff)
			}
		})
	}
}
