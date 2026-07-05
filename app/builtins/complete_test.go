package builtins

import (
	"bytes"
	"testing"

	"github.com/codecrafters-io/shell-starter-go/app/completion"
	"github.com/google/go-cmp/cmp"
)

func TestComplete(t *testing.T) {
	tests := []struct {
		name           string
		initial        map[string]string
		args           []string
		wantOut        string
		wantErr        string
		wantRegistered map[string]string
	}{
		{
			name:    "prints missing specification for -p",
			args:    []string{"-p", "git"},
			wantErr: "complete: git: no completion specification\n",
		},
		{
			name:    "displays registered completion for -p",
			initial: map[string]string{"git": "/path/to/script"},
			args:    []string{"-p", "git"},
			wantOut: "complete -C '/path/to/script' git\n",
		},
		{
			name:    "displays docker completion for -p",
			initial: map[string]string{"docker": "/path/to/docker/completer"},
			args:    []string{"-p", "docker"},
			wantOut: "complete -C '/path/to/docker/completer' docker\n",
		},
		{
			name:           "registers completion with -C",
			args:           []string{"-C", "/path/to/script", "git"},
			wantRegistered: map[string]string{"git": "/path/to/script"},
		},
		{
			name:    "ignores -C with too few arguments",
			args:    []string{"-C", "/path/to/script"},
			wantOut: "",
			wantErr: "",
		},
		{
			name:    "ignores bare complete",
			args:    nil,
			wantOut: "",
			wantErr: "",
		},
		{
			name:           "unregisters completion with -r",
			initial:        map[string]string{"git": "/path/to/script"},
			args:           []string{"-r", "git"},
			wantRegistered: map[string]string{},
		},
		{
			name:           "ignores -r for unregistered command",
			args:           []string{"-r", "git"},
			wantRegistered: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			registry := completion.NewCompletionRegistry()
			for command, scriptPath := range tt.initial {
				registry.Register(command, scriptPath)
			}

			var stdout, stderr bytes.Buffer
			Complete(&stdout, &stderr, tt.args, registry)

			if diff := cmp.Diff(tt.wantOut, stdout.String()); diff != "" {
				t.Errorf("Complete(%v) stdout mismatch (-want +got):\n%s", tt.args, diff)
			}
			if diff := cmp.Diff(tt.wantErr, stderr.String()); diff != "" {
				t.Errorf("Complete(%v) stderr mismatch (-want +got):\n%s", tt.args, diff)
			}
			if tt.wantRegistered != nil {
				for command, wantPath := range tt.wantRegistered {
					gotPath, ok := registry.Lookup(command)
					if wantPath == "" {
						if ok {
							t.Errorf("registry.Lookup(%q) = %q, want missing", command, gotPath)
						}
						continue
					}
					if !ok || gotPath != wantPath {
						t.Errorf("registry.Lookup(%q) = (%q, %v), want %q", command, gotPath, ok, wantPath)
					}
				}
			}
		})
	}
}
