package builtins

import (
	"bytes"
	"maps"
	"testing"
)

func TestComplete(t *testing.T) {
	tests := []struct {
		name                     string
		registeredCompleters     map[string]string
		args                     []string
		wantOut                  string
		wantErr                  string
		wantRegisteredCompleters map[string]string
	}{
		{
			name:    "prints missing specification for -p",
			args:    []string{"-p", "git"},
			wantErr: "complete: git: no completion specification\n",
		},
		{
			name:                 "displays registered completion for -p",
			registeredCompleters: map[string]string{"git": "/path/to/script"},
			args:                 []string{"-p", "git"},
			wantOut:              "complete -C '/path/to/script' git\n",
		},
		{
			name:                 "displays docker completion for -p",
			registeredCompleters: map[string]string{"docker": "/path/to/docker/completer"},
			args:                 []string{"-p", "docker"},
			wantOut:              "complete -C '/path/to/docker/completer' docker\n",
		},
		{
			name:                     "registers completion with -C",
			args:                     []string{"-C", "/path/to/script", "git"},
			wantRegisteredCompleters: map[string]string{"git": "/path/to/script"},
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
			name:                     "unregisters completion with -r",
			registeredCompleters:     map[string]string{"git": "/path/to/script"},
			args:                     []string{"-r", "git"},
			wantRegisteredCompleters: map[string]string{},
		},
		{
			name:                     "ignores -r for unregistered command",
			args:                     []string{"-r", "git"},
			wantRegisteredCompleters: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			registeredCompleters := map[string]string{}
			maps.Copy(registeredCompleters, tt.registeredCompleters)

			var stdout, stderr bytes.Buffer
			Complete(&stdout, &stderr, tt.args, registeredCompleters)
			if got := stdout.String(); got != tt.wantOut {
				t.Errorf("Complete(%v) stdout = %q, want %q", tt.args, got, tt.wantOut)
			}
			if got := stderr.String(); got != tt.wantErr {
				t.Errorf("Complete(%v) stderr = %q, want %q", tt.args, got, tt.wantErr)
			}
			if tt.wantRegisteredCompleters != nil {
				if len(registeredCompleters) != len(tt.wantRegisteredCompleters) {
					t.Fatalf("registeredCompleters has %d entries, want %d", len(registeredCompleters), len(tt.wantRegisteredCompleters))
				}
				for command, want := range tt.wantRegisteredCompleters {
					got, ok := registeredCompleters[command]
					if !ok {
						t.Fatalf("registeredCompleters missing %q", command)
					}
					if got != want {
						t.Errorf("registeredCompleters[%q] = %q, want %q", command, got, want)
					}
				}
			}
		})
	}
}
