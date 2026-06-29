package completion

import "testing"

func TestApplyTab(t *testing.T) {
	t.Setenv("PATH", "")

	builtins := []string{"cd", "echo", "exit", "pwd", "type"}

	tests := []struct {
		name         string
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBuffer, gotListings := ApplyTab(builtins, tt.buffer)
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
