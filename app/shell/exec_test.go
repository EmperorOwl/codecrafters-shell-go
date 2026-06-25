package shell

import (
	"io"
	"reflect"
	"testing"
)

func TestNewExternalCommand(t *testing.T) {
	tests := []struct {
		name            string
		fields          []string
		executablePath  string
		wantArgs        []string
		wantProgramPath string
	}{
		{
			name:            "program with one argument",
			fields:          []string{"custom_exe", "alice"},
			executablePath:  "/usr/local/bin/custom_exe",
			wantArgs:        []string{"custom_exe", "alice"},
			wantProgramPath: "/usr/local/bin/custom_exe",
		},
		{
			name:            "program without arguments",
			fields:          []string{"custom_exe"},
			executablePath:  "/usr/local/bin/custom_exe",
			wantArgs:        []string{"custom_exe"},
			wantProgramPath: "/usr/local/bin/custom_exe",
		},
		{
			name:            "quoted program name with spaces",
			fields:          []string{`exe with "quotes"`, "file"},
			executablePath:  "/tmp/cow/exe with \"quotes\"",
			wantArgs:        []string{`exe with "quotes"`, "file"},
			wantProgramPath: `/tmp/cow/exe with "quotes"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := newExternalCommand(tt.fields, tt.executablePath, io.Discard)
			if cmd.Path != tt.wantProgramPath {
				t.Errorf("Path = %q, want %q", cmd.Path, tt.wantProgramPath)
			}
			if !reflect.DeepEqual(cmd.Args, tt.wantArgs) {
				t.Errorf("Args = %v, want %v", cmd.Args, tt.wantArgs)
			}
		})
	}
}
