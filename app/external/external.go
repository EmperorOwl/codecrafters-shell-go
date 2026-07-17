package external

import (
	"io"
	"os"
	"os/exec"
)

// ExternalProgram describes an external command to run.
type ExternalProgram struct {
	Name   string
	Path   string
	Args   []string
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

// New resolves fields[0] against PATH and returns an ExternalProgram.
// The second return value is false when fields is empty or the command is not found.
func New(fields []string, stdout, stderr io.Writer) (*ExternalProgram, bool) {
	if len(fields) == 0 {
		return nil, false
	}

	path, ok := FindExecutableInPath(fields[0])
	if !ok {
		return nil, false
	}

	return &ExternalProgram{
		Name:   fields[0],
		Path:   path,
		Args:   fields[1:],
		Stdout: stdout,
		Stderr: stderr,
	}, true
}

func (p *ExternalProgram) command() *exec.Cmd {
	cmd := exec.Command(p.Path, p.Args...)
	cmd.Args = append([]string{p.Name}, p.Args...)
	stdin := p.Stdin
	if stdin == nil {
		stdin = os.Stdin
	}
	cmd.Stdin = stdin
	cmd.Stdout = p.Stdout
	cmd.Stderr = p.Stderr
	return cmd
}

// Run starts the program and waits for it to exit.
func (p *ExternalProgram) Run() error {
	return p.command().Run()
}

// RunInBackground starts the program and reaps it in a goroutine.
// onStarted is called synchronously after Start with the child PID; onExit is
// called after the process exits. Either callback may be nil.
func (p *ExternalProgram) RunInBackground(onStarted func(int), onExit func()) (pid int, err error) {
	cmd := p.command()
	if err := cmd.Start(); err != nil {
		return 0, err
	}

	pid = cmd.Process.Pid
	if onStarted != nil {
		onStarted(pid)
	}
	go func() {
		_ = cmd.Wait()
		if onExit != nil {
			onExit()
		}
	}()
	return pid, nil
}
