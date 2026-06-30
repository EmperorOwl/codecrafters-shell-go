//go:build windows

package jobs

import (
	"os/exec"
	"testing"
	"time"
)

func TestProcessExited(t *testing.T) {
	cmd := exec.Command("cmd", "/c", "exit", "0")
	if err := cmd.Start(); err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if processExited(cmd.Process.Pid) {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}

	t.Fatalf("processExited(%d) = false, want true after child exited", cmd.Process.Pid)
}
