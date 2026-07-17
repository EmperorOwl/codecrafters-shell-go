package external

import (
	"bytes"
	"io"
	"testing"
	"time"

	"github.com/codecrafters-io/shell-starter-go/app/testutils"
	"github.com/google/go-cmp/cmp"
)

func TestExternalProgram_Run(t *testing.T) {
	dir := t.TempDir()
	name, programPath := testutils.WriteMockProgram(t, dir)

	tests := []struct {
		name    string
		args    []string
		wantOut string
	}{
		{
			name:    "prints greeting with argument",
			args:    []string{"world"},
			wantOut: "hello world\n",
		},
		{
			name:    "prints greeting without argument",
			args:    nil,
			wantOut: "hello\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out bytes.Buffer
			prog := &ExternalProgram{
				Name:   name,
				Path:   programPath,
				Args:   tt.args,
				Stdout: &out,
				Stderr: io.Discard,
			}

			if err := prog.Run(); err != nil {
				t.Fatalf("Run() error = %v", err)
			}
			if diff := cmp.Diff(tt.wantOut, out.String()); diff != "" {
				t.Errorf("stdout mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestExternalProgram_RunInBackground(t *testing.T) {
	dir := t.TempDir()
	name, programPath := testutils.WriteMockProgram(t, dir)

	exited := make(chan struct{})
	prog := &ExternalProgram{
		Name:   name,
		Path:   programPath,
		Stdout: io.Discard,
		Stderr: io.Discard,
	}

	pid, err := prog.RunInBackground(nil, func() { close(exited) })
	if err != nil {
		t.Fatalf("RunInBackground() error = %v", err)
	}
	if pid <= 0 {
		t.Fatalf("RunInBackground() pid = %d, want > 0", pid)
	}

	select {
	case <-exited:
	case <-time.After(5 * time.Second):
		t.Fatal("onExit was not called within 5s")
	}
}
