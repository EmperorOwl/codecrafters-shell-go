package builtins

import (
	"bytes"
	"testing"

	"github.com/codecrafters-io/shell-starter-go/app/completion"
	"github.com/google/go-cmp/cmp"
)

func TestComplete(t *testing.T) {
	tests := []struct {
		name         string
		setup        func(*completion.Registry)
		args         []string
		wantOut      string
		wantErr      string
		wantRegistry map[string]string
	}{
		{
			name:    "prints missing specification for -p",
			args:    []string{"-p", "git"},
			wantErr: "complete: git: no completion specification\n",
		},
		{
			name: "displays registered completion for -p",
			setup: func(registry *completion.Registry) {
				registry.Register("git", "/path/to/script")
			},
			args:    []string{"-p", "git"},
			wantOut: "complete -C '/path/to/script' git\n",
		},
		{
			name: "displays docker completion for -p",
			setup: func(registry *completion.Registry) {
				registry.Register("docker", "/path/to/docker/completer")
			},
			args:    []string{"-p", "docker"},
			wantOut: "complete -C '/path/to/docker/completer' docker\n",
		},
		{
			name:         "registers completion with -C",
			args:         []string{"-C", "/path/to/script", "git"},
			wantRegistry: map[string]string{"git": "/path/to/script"},
		},
		{
			name: "ignores -C with too few arguments",
			args: []string{"-C", "/path/to/script"},
		},
		{
			name: "ignores bare complete",
			args: nil,
		},
		{
			name: "unregisters completion with -r",
			setup: func(registry *completion.Registry) {
				registry.Register("git", "/path/to/script")
			},
			args:         []string{"-r", "git"},
			wantRegistry: map[string]string{},
		},
		{
			name: "ignores -r for unregistered command",
			args: []string{"-r", "git"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			registry := completion.NewRegistry()
			if tt.setup != nil {
				tt.setup(registry)
			}

			var stdout, stderr bytes.Buffer
			completeBuiltin(&stdout, &stderr, tt.args, registry)

			if tt.wantRegistry != nil {
				if diff := cmp.Diff(tt.wantRegistry, registry.Entries()); diff != "" {
					t.Errorf("registry entries mismatch (-want +got):\n%s", diff)
				}
			}

			if diff := cmp.Diff(tt.wantOut, stdout.String()); diff != "" {
				t.Errorf("completeBuiltin(%v) stdout mismatch (-want +got):\n%s", tt.args, diff)
			}
			if diff := cmp.Diff(tt.wantErr, stderr.String()); diff != "" {
				t.Errorf("completeBuiltin(%v) stderr mismatch (-want +got):\n%s", tt.args, diff)
			}
		})
	}
}
