package completion

import "testing"

func TestApplyTab(t *testing.T) {
	builtins := []string{"cd", "echo", "exit", "pwd", "type"}

	tests := []struct {
		name         string
		executables  []string
		fileDirs     map[string][]string
		buffer       string
		wantBuffer   string
		wantListings []string
	}{
		{
			name:       "completes echo",
			buffer:     "ech",
			wantBuffer: "echo ",
		},
		{
			name:       "completes exit",
			buffer:     "exi",
			wantBuffer: "exit ",
		},
		{
			name:       "no match",
			buffer:     "xyz",
			wantBuffer: "xyz",
		},
		{
			name:         "ambiguous prefix lists matches",
			buffer:       "e",
			wantBuffer:   "e",
			wantListings: []string{"echo", "exit"},
		},
		{
			name:         "empty buffer lists all builtins",
			buffer:       "",
			wantBuffer:   "",
			wantListings: builtins,
		},
		{
			name:        "completes executable",
			executables: []string{"custom_executable"},
			buffer:      "custom",
			wantBuffer:  "custom_executable ",
		},
		{
			name:         "ambiguous executable prefix lists matches",
			executables:  []string{"xyz_bar", "xyz_baz", "xyz_quz"},
			buffer:       "xyz_",
			wantBuffer:   "xyz_",
			wantListings: []string{"xyz_bar", "xyz_baz", "xyz_quz"},
		},
		{
			name:        "completes to longest common prefix",
			executables: []string{"xyz_foo", "xyz_foo_bar", "xyz_foo_bar_baz"},
			buffer:      "xyz_",
			wantBuffer:  "xyz_foo",
		},
		{
			name:        "completes to next longest common prefix",
			executables: []string{"xyz_foo", "xyz_foo_bar", "xyz_foo_bar_baz"},
			buffer:      "xyz_foo_",
			wantBuffer:  "xyz_foo_bar",
		},
		{
			name:        "completes final unique executable",
			executables: []string{"xyz_foo", "xyz_foo_bar", "xyz_foo_bar_baz"},
			buffer:      "xyz_foo_bar_",
			wantBuffer:  "xyz_foo_bar_baz ",
		},
		{
			name:       "completes partial filename",
			fileDirs:   map[string][]string{"": {"hello_world.py", "notes.md", "readme.txt"}},
			buffer:     "cat re",
			wantBuffer: "cat readme.txt ",
		},
		{
			name:       "completes hello prefix",
			fileDirs:   map[string][]string{"": {"hello_world.py", "notes.md", "readme.txt"}},
			buffer:     "cat hello",
			wantBuffer: "cat hello_world.py ",
		},
		{
			name:       "completes for unknown command",
			fileDirs:   map[string][]string{"": {"hello_world.py", "notes.md", "readme.txt"}},
			buffer:     "xyz read",
			wantBuffer: "xyz readme.txt ",
		},
		{
			name:       "no file match leaves buffer unchanged",
			fileDirs:   map[string][]string{"": {"hello_world.py", "notes.md", "readme.txt"}},
			buffer:     "cat missing",
			wantBuffer: "cat missing",
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			listFiles := func(dir string) []string {
				if tt.fileDirs == nil {
					return nil
				}
				return tt.fileDirs[dir]
			}

			gotBuffer, gotListings := ApplyTab(builtins, tt.executables, listFiles, tt.buffer)
			if gotBuffer != tt.wantBuffer {
				t.Errorf("ApplyTab(%q) buffer = %q, want %q", tt.buffer, gotBuffer, tt.wantBuffer)
			}
			if len(gotListings) != len(tt.wantListings) {
				t.Fatalf("ApplyTab(%q) listings = %v, want %v", tt.buffer, gotListings, tt.wantListings)
			}
			for i := range tt.wantListings {
				if gotListings[i] != tt.wantListings[i] {
					t.Errorf("ApplyTab(%q) listings[%d] = %q, want %q", tt.buffer, i, gotListings[i], tt.wantListings[i])
				}
			}
		})
	}
}
