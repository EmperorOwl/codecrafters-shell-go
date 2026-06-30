package shell

import (
	"io"
	"os"
	"os/exec"

	shellpath "github.com/codecrafters-io/shell-starter-go/app/path"
)

func newExternalCommand(fields []string, executablePath string, stdout, stderr io.Writer) *exec.Cmd {
	cmd := exec.Command(executablePath, fields[1:]...)
	cmd.Args = append([]string{fields[0]}, fields[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	return cmd
}

func ExecuteExternalProgram(fields []string, stdout, stderr io.Writer) (executed bool, err error) {
	if len(fields) == 0 {
		return false, nil
	}

	path, ok := shellpath.FindExecutableInPath(fields[0])
	if !ok {
		return false, nil
	}

	return true, newExternalCommand(fields, path, stdout, stderr).Run()
}

func StartExternalProgram(fields []string, stdout, stderr io.Writer) (executed bool, pid int, err error) {
	if len(fields) == 0 {
		return false, 0, nil
	}

	path, ok := shellpath.FindExecutableInPath(fields[0])
	if !ok {
		return false, 0, nil
	}

	cmd := newExternalCommand(fields, path, stdout, stderr)
	if err := cmd.Start(); err != nil {
		return true, 0, err
	}
	return true, cmd.Process.Pid, nil
}
