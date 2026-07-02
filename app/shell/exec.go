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

func StartExternalProgram(fields []string, stdout, stderr io.Writer) (executed bool, pid int, cmd *exec.Cmd, err error) {
	if len(fields) == 0 {
		return false, 0, nil, nil
	}

	path, ok := shellpath.FindExecutableInPath(fields[0])
	if !ok {
		return false, 0, nil, nil
	}

	cmd = newExternalCommand(fields, path, stdout, stderr)
	if err := cmd.Start(); err != nil {
		return true, 0, nil, err
	}
	return true, cmd.Process.Pid, cmd, nil
}

func findExecutable(fields []string) (path string, ok bool) {
	if len(fields) == 0 {
		return "", false
	}
	return shellpath.FindExecutableInPath(fields[0])
}

func ExecutePipeline(fieldsList [2][]string, stdout, stderr io.Writer) (executed bool, notFound string, err error) {
	for _, fields := range fieldsList {
		if len(fields) == 0 {
			return false, "", nil
		}
		if _, ok := findExecutable(fields); !ok {
			return false, fields[0], nil
		}
	}

	path0, _ := findExecutable(fieldsList[0])
	path1, _ := findExecutable(fieldsList[1])

	cmd0 := newExternalCommand(fieldsList[0], path0, nil, stderr)
	cmd1 := newExternalCommand(fieldsList[1], path1, stdout, stderr)

	pipeReader, err := cmd0.StdoutPipe()
	if err != nil {
		return true, "", err
	}
	cmd1.Stdin = pipeReader

	if err := cmd0.Start(); err != nil {
		return true, "", err
	}
	if err := cmd1.Start(); err != nil {
		_ = cmd0.Process.Kill()
		_ = cmd0.Wait()
		return true, "", err
	}

	err1 := cmd1.Wait()
	_ = cmd0.Wait()
	return true, "", err1
}

func startBackgroundWait(cmd *exec.Cmd, onExit func()) {
	if cmd == nil || onExit == nil {
		return
	}
	go func() {
		_ = cmd.Wait()
		onExit()
	}()
}
