package builtins

import (
	"io"

	"github.com/codecrafters-io/shell-starter-go/app/jobs"
)

func Jobs(out io.Writer, jobList []jobs.Job) {
	jobs.WriteAll(out, jobList)
}
