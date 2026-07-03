package shell

import (
	"bytes"
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

func (s *Shell) runBuiltin(fields []string, stdout, stderr io.Writer) (shouldExit bool) {
	_, shouldExit = TryBuiltin(fields, stdout, stderr, s.completers, &s.jobs)
	return shouldExit
}

func (s *Shell) runBuiltinDrainingStdin(fields []string, stdin io.Reader, stdout, stderr io.Writer) (shouldExit bool) {
	drainDone := make(chan struct{})
	go func() {
		_, _ = io.Copy(io.Discard, stdin)
		close(drainDone)
	}()

	shouldExit = s.runBuiltin(fields, stdout, stderr)
	<-drainDone
	return shouldExit
}

func (s *Shell) executeExternalBuiltinPipeline(fields0, fields1 []string, stdout, stderr io.Writer) (bool, string, error) {
	pr, pw := io.Pipe()

	path, _ := findExecutable(fields0)
	cmd0 := newExternalCommand(fields0, path, pw, stderr)
	cmd0.Stdin = bytes.NewReader(nil)
	if err := cmd0.Start(); err != nil {
		_ = pw.Close()
		_ = pr.Close()
		return true, "", err
	}

	drainDone := make(chan struct{})
	go func() {
		_, _ = io.Copy(io.Discard, pr)
		close(drainDone)
	}()

	if s.runBuiltin(fields1, stdout, stderr) {
		_ = cmd0.Process.Kill()
		_ = cmd0.Wait()
		_ = pw.Close()
		<-drainDone
		return true, "", nil
	}

	_ = cmd0.Wait()
	_ = pw.Close()
	<-drainDone
	return true, "", nil
}

func (s *Shell) ExecutePipeline(fieldsList [2][]string, stdout, stderr io.Writer) (executed bool, notFound string, err error) {
	fields0, fields1 := fieldsList[0], fieldsList[1]
	if len(fields0) == 0 || len(fields1) == 0 {
		return false, "", nil
	}

	builtin0 := IsShellBuiltin(fields0[0])
	builtin1 := IsShellBuiltin(fields1[0])

	if !builtin0 {
		if _, ok := findExecutable(fields0); !ok {
			return false, fields0[0], nil
		}
	}
	if !builtin1 {
		if _, ok := findExecutable(fields1); !ok {
			return false, fields1[0], nil
		}
	}

	switch {
	case builtin0 && builtin1:
		return s.executeBuiltinBuiltinPipeline(fields0, fields1, stdout, stderr)
	case builtin0 && !builtin1:
		return s.executeBuiltinExternalPipeline(fields0, fields1, stdout, stderr)
	case !builtin0 && builtin1:
		return s.executeExternalBuiltinPipeline(fields0, fields1, stdout, stderr)
	default:
		return s.executeExternalExternalPipeline(fields0, fields1, stdout, stderr)
	}
}

func (s *Shell) executeBuiltinBuiltinPipeline(fields0, fields1 []string, stdout, stderr io.Writer) (bool, string, error) {
	pr, pw := io.Pipe()

	done := make(chan struct{})
	go func() {
		s.runBuiltinDrainingStdin(fields1, pr, stdout, stderr)
		close(done)
	}()

	if s.runBuiltin(fields0, pw, stderr) {
		_ = pw.Close()
		<-done
		return true, "", nil
	}
	_ = pw.Close()
	<-done
	return true, "", nil
}

func (s *Shell) executeBuiltinExternalPipeline(fields0, fields1 []string, stdout, stderr io.Writer) (bool, string, error) {
	pr, pw := io.Pipe()

	path, _ := findExecutable(fields1)
	cmd1 := newExternalCommand(fields1, path, stdout, stderr)
	cmd1.Stdin = pr

	errCh := make(chan error, 1)
	go func() {
		errCh <- cmd1.Run()
	}()

	if s.runBuiltin(fields0, pw, stderr) {
		_ = pw.Close()
		<-errCh
		return true, "", nil
	}
	_ = pw.Close()

	err := <-errCh
	return true, "", err
}

func (s *Shell) executeExternalExternalPipeline(fields0, fields1 []string, stdout, stderr io.Writer) (bool, string, error) {
	path0, _ := findExecutable(fields0)
	path1, _ := findExecutable(fields1)

	cmd0 := newExternalCommand(fields0, path0, nil, stderr)
	cmd1 := newExternalCommand(fields1, path1, stdout, stderr)

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
