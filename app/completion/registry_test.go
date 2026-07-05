package completion

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCompletionRegistry_Register(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(*CompletionRegistry)
		command    string
		scriptPath string
		wantPath   string
	}{
		{
			name:       "registers script",
			setup:      func(*CompletionRegistry) {},
			command:    "git",
			scriptPath: "/path/to/git-completer",
			wantPath:   "/path/to/git-completer",
		},
		{
			name: "overwrites existing script",
			setup: func(r *CompletionRegistry) {
				r.Register("git", "/old/script")
			},
			command:    "git",
			scriptPath: "/new/script",
			wantPath:   "/new/script",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			registry := NewCompletionRegistry()
			tt.setup(registry)

			registry.Register(tt.command, tt.scriptPath)

			gotPath, ok := registry.Lookup(tt.command)
			if !ok {
				t.Fatalf("Lookup(%q) missing, want %q", tt.command, tt.wantPath)
			}
			if diff := cmp.Diff(tt.wantPath, gotPath); diff != "" {
				t.Errorf("Lookup(%q) mismatch (-want +got):\n%s", tt.command, diff)
			}
		})
	}
}

func TestCompletionRegistry_Unregister(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(*CompletionRegistry)
		command   string
		wantFound bool
		wantPath  string
	}{
		{
			name: "removes registered script",
			setup: func(r *CompletionRegistry) {
				r.Register("git", "/path/to/script")
			},
			command:   "git",
			wantFound: false,
		},
		{
			name:      "missing command is no-op",
			setup:     func(*CompletionRegistry) {},
			command:   "git",
			wantFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			registry := NewCompletionRegistry()
			tt.setup(registry)

			registry.Unregister(tt.command)

			gotPath, ok := registry.Lookup(tt.command)
			if ok != tt.wantFound {
				if tt.wantFound {
					t.Fatalf("Lookup(%q) missing, want %q", tt.command, tt.wantPath)
				}
				t.Errorf("Lookup(%q) = %q, want missing", tt.command, gotPath)
			}
			if tt.wantFound {
				if diff := cmp.Diff(tt.wantPath, gotPath); diff != "" {
					t.Errorf("Lookup(%q) mismatch (-want +got):\n%s", tt.command, diff)
				}
			}
		})
	}
}

func TestCompletionRegistry_Lookup(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(*CompletionRegistry)
		command   string
		wantFound bool
		wantPath  string
	}{
		{
			name:      "empty registry",
			setup:     func(*CompletionRegistry) {},
			command:   "git",
			wantFound: false,
		},
		{
			name: "returns registered script",
			setup: func(r *CompletionRegistry) {
				r.Register("git", "/path/to/script")
			},
			command:   "git",
			wantFound: true,
			wantPath:  "/path/to/script",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			registry := NewCompletionRegistry()
			tt.setup(registry)

			gotPath, ok := registry.Lookup(tt.command)
			if ok != tt.wantFound {
				if tt.wantFound {
					t.Fatalf("Lookup(%q) missing, want %q", tt.command, tt.wantPath)
				}
				t.Errorf("Lookup(%q) = %q, want missing", tt.command, gotPath)
			}
			if tt.wantFound {
				if diff := cmp.Diff(tt.wantPath, gotPath); diff != "" {
					t.Errorf("Lookup(%q) mismatch (-want +got):\n%s", tt.command, diff)
				}
			}
		})
	}
}
