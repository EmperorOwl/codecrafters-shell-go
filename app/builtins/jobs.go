package builtins

import (
	"io"

	"github.com/codecrafters-io/shell-starter-go/app/jobs"
)

// Jobs prints all jobs in the table, then reaps finished ones.
func Jobs(out io.Writer, table *jobs.JobTable) {
	display := table.ListForDisplay()
	table.ReapDone()
	jobs.WriteAll(out, display)
}

func jobsBuiltin(ctx *Context, args []string) (bool, error) {
	Jobs(ctx.Stdout, ctx.Jobs)
	return false, nil
}
