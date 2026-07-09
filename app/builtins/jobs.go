package builtins

import (
	"io"

	"github.com/codecrafters-io/shell-starter-go/app/jobs"
)

func init() {
	register("jobs", jobsBuiltin)
}

func jobsBuiltin(ctx *Context, args []string) (bool, error) {
	if ctx.State == nil {
		return false, nil
	}
	Jobs(ctx.Stdout, ctx.State.Jobs)
	return false, nil
}

// Jobs prints all jobs in the table, then reaps finished ones.
func Jobs(out io.Writer, table *jobs.JobTable) {
	display := table.List()
	table.ReapDone()
	jobs.WriteAll(out, display)
}
