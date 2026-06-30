package shell

import (
	"io"
	"runtime"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
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
			cmd := newExternalCommand(tt.fields, tt.executablePath, io.Discard, io.Discard)
			if diff := cmp.Diff(tt.wantProgramPath, cmd.Path); diff != "" {
				t.Errorf("Path mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tt.wantArgs, cmd.Args, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("Args mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestStartExternalProgram(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("sleep is not available on Windows")
	}

	executed, pid, cmd, err := StartExternalProgram([]string{"sleep", "30"}, io.Discard, io.Discard)
	if err != nil {
		t.Fatalf("StartExternalProgram() error = %v", err)
	}
	if !executed {
		t.Fatal("StartExternalProgram() executed = false, want true")
	}
	if pid <= 0 {
		t.Fatalf("StartExternalProgram() pid = %d, want > 0", pid)
	}
	if cmd == nil {
		t.Fatal("StartExternalProgram() cmd = nil, want non-nil")
	}
	_ = cmd.Process.Kill()

	executed, pid, cmd, err = StartExternalProgram([]string{"missing_command_xyz"}, io.Discard, io.Discard)
	if err != nil {
		t.Fatalf("StartExternalProgram() error = %v", err)
	}
	if executed {
		t.Fatalf("StartExternalProgram() executed = true, want false")
	}
	if pid != 0 {
		t.Fatalf("StartExternalProgram() pid = %d, want 0", pid)
	}
}
