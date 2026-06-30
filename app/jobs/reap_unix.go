//go:build !windows

package jobs

import "syscall"

func processExited(pid int) bool {
	var status syscall.WaitStatus
	reaped, err := syscall.Wait4(pid, &status, syscall.WNOHANG, nil)
	if err != nil || reaped == 0 {
		return false
	}
	return status.Exited()
}
