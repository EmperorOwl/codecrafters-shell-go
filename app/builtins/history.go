package builtins

import (
	"fmt"
	"io"
	"strconv"

	"github.com/codecrafters-io/shell-starter-go/app/history"
)

func init() {
	register("history", historyBuiltin)
}

func historyBuiltin(ctx *Context, args []string) (bool, error) {
	if ctx.State == nil {
		return false, nil
	}
	if len(args) >= 2 && args[0] == "-r" {
		if err := ctx.State.History.ReadFromFile(args[1]); err != nil {
			fmt.Fprintf(ctx.Stderr, "history: %v\n", err)
		}
		return false, nil
	}
	if len(args) >= 2 && args[0] == "-w" {
		if err := ctx.State.History.WriteToFile(args[1]); err != nil {
			fmt.Fprintf(ctx.Stderr, "history: %v\n", err)
		}
		return false, nil
	}
	if len(args) >= 2 && args[0] == "-a" {
		if err := ctx.State.History.AppendToFile(args[1]); err != nil {
			fmt.Fprintf(ctx.Stderr, "history: %v\n", err)
		}
		return false, nil
	}
	History(ctx.Stdout, ctx.State.History, parseHistoryLimit(args))
	return false, nil
}

// History prints commands from the history list. A positive limit shows only
// the last n entries; zero shows the full history.
func History(out io.Writer, list *history.HistoryList, limit int) {
	var entries []history.Entry
	if limit > 0 {
		entries = list.ListLast(limit)
	} else {
		entries = list.List()
	}
	history.WriteAll(out, entries)
}

func parseHistoryLimit(args []string) int {
	if len(args) == 0 {
		return 0
	}
	limit, err := strconv.Atoi(args[0])
	if err != nil || limit <= 0 {
		return 0
	}
	return limit
}
