package builtins

import (
	"io"

	"github.com/codecrafters-io/shell-starter-go/app/history"
)

func init() {
	register("history", historyBuiltin)
}

func historyBuiltin(ctx *Context, args []string) (bool, error) {
	if ctx.State == nil {
		return false, nil
	}
	History(ctx.Stdout, ctx.State.History)
	return false, nil
}

// History prints all commands in the history table.
func History(out io.Writer, list *history.HistoryList) {
	history.WriteAll(out, list.List())
}
