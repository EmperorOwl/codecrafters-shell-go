package completion

import "testing"

func TestApplyTab(t *testing.T) {
	builtins := []string{"cd", "echo", "exit", "pwd", "type"}

	tests := []struct {
		name         string
		executables  []string
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
			name:         "completes executable",
			executables:  []string{"custom_executable"},
			buffer:       "custom",
			wantBuffer:   "custom_executable ",
		},
		{
			name:         "ambiguous executable prefix lists matches",
			executables:  []string{"xyz_bar", "xyz_baz", "xyz_quz"},
			buffer:       "xyz_",
			wantBuffer:   "xyz_",
			wantListings: []string{"xyz_bar", "xyz_baz", "xyz_quz"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBuffer, gotListings := ApplyTab(builtins, tt.executables, tt.buffer)
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
