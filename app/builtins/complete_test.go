package builtins

import (
	"bytes"
	"maps"
	"testing"
)

func TestComplete(t *testing.T) {
	runCompleter := CompleterFunc(func(opts CompleterFuncOptions) ([]string, error) {
		return []string{opts.ScriptPath}, nil
	})

	tests := []struct {
		name                     string
		registeredCompleters     map[string]Completer
		args                     []string
		wantOut                  string
		wantErr                  string
		wantRegisteredCompleters map[string]Completer
	}{
		{
			name:    "prints missing specification for -p",
			args:    []string{"-p", "git"},
			wantErr: "complete: git: no completion specification\n",
		},
		{
			name:                 "displays registered completion for -p",
			registeredCompleters: map[string]Completer{"git": {Path: "/path/to/script"}},
			args:                 []string{"-p", "git"},
			wantOut:              "complete -C '/path/to/script' git\n",
		},
		{
			name:                 "displays docker completion for -p",
			registeredCompleters: map[string]Completer{"docker": {Path: "/path/to/docker/completer"}},
			args:                 []string{"-p", "docker"},
			wantOut:              "complete -C '/path/to/docker/completer' docker\n",
		},
		{
			name:                     "registers completion with -C",
			args:                     []string{"-C", "/path/to/script", "git"},
			wantRegisteredCompleters: map[string]Completer{"git": {Path: "/path/to/script", Func: runCompleter}},
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
			registeredCompleters:     map[string]Completer{"git": {Path: "/path/to/script"}},
			args:                     []string{"-r", "git"},
			wantRegisteredCompleters: map[string]Completer{},
		},
		{
			name:                     "ignores -r for unregistered command",
			args:                     []string{"-r", "git"},
			wantRegisteredCompleters: map[string]Completer{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			registeredCompleters := map[string]Completer{}
			maps.Copy(registeredCompleters, tt.registeredCompleters)

			var stdout, stderr bytes.Buffer
			Complete(&stdout, &stderr, tt.args, registeredCompleters, runCompleter)
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
					if got.Path != want.Path {
						t.Errorf("registeredCompleters[%q].Path = %q, want %q", command, got.Path, want.Path)
					}
					if (got.Func == nil) != (want.Func == nil) {
						t.Errorf("registeredCompleters[%q].Func set = %v, want %v", command, got.Func != nil, want.Func != nil)
					}
				}
			}
		})
	}
}
