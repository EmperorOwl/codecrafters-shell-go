package builtins

import (
	"io"

	"github.com/codecrafters-io/shell-starter-go/app/jobs"
)

func init() {
	register("jobs", jobsHandler)
}

func jobsHandler(ctx *Context, args []string) (bool, error) {
	if ctx.Session == nil {
		return false, nil
	}
	jobsBuiltin(ctx.Stdout, ctx.Session.Jobs)
	return false, nil
}

func jobsBuiltin(stdout io.Writer, table *jobs.Table) {
	display := table.List()
	table.ReapDone()
	jobs.WriteAll(stdout, display)
}
