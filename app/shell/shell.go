package shell

import (
	"io"
	"os"

	"github.com/codecrafters-io/shell-starter-go/app/completion"
	"github.com/codecrafters-io/shell-starter-go/app/files"
	"github.com/codecrafters-io/shell-starter-go/app/jobs"
)

type Shell struct {
	jobs       jobs.JobTable
	completers map[string]string
}

func New() *Shell {
	return &Shell{
		completers: make(map[string]string),
	}
}

func CommandNotFoundMessage(command string) string {
	return command + ": command not found"
}

func (s *Shell) PrintReapedJobs(out io.Writer) {
	done := s.jobs.ReapDone()
	if len(done) > 0 {
		jobs.WriteAll(out, done)
	}
}

func (s *Shell) listFiles(dir string) []string {
	cwd, err := os.Getwd()
	if err != nil {
		return nil
	}
	return files.ListInDir(cwd, dir)
}

func (s *Shell) complete(opts completion.CompleterFuncOptions) []string {
	return completion.CompleteCommand(s.completers, opts)
}
