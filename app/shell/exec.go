package shell

import (
	"os"
	"os/exec"

	shellpath "github.com/codecrafters-io/shell-starter-go/app/path"
)

func newExternalCommand(fields []string, executablePath string) *exec.Cmd {
	cmd := exec.Command(executablePath, fields[1:]...)
	cmd.Args = append([]string{fields[0]}, fields[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}

func ExecuteExternalProgram(fields []string) (executed bool, err error) {
	if len(fields) == 0 {
		return false, nil
	}

	path, ok := shellpath.FindExecutableInPath(fields[0], os.Getenv("PATH"))
	if !ok {
		return false, nil
	}

	return true, newExternalCommand(fields, path).Run()
}
