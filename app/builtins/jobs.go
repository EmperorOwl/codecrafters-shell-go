package builtins

import (
	"io"

	"github.com/codecrafters-io/shell-starter-go/app/jobs"
)

func Jobs(out io.Writer, jobList *[]jobs.Job, hasExited func(int) bool) {
	jobs.WriteAllWithChecker(out, jobList, hasExited)
}
