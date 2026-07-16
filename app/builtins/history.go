package builtins

import (
	"fmt"
	"io"
	"strconv"

	"github.com/codecrafters-io/shell-starter-go/app/history"
)

func init() {
	register("history", historyHandler)
}

func historyHandler(ctx *Context, args []string) (bool, error) {
	if ctx.State == nil {
		return false, nil
	}
	historyBuiltin(ctx.Stdout, ctx.Stderr, args, ctx.State.History)
	return false, nil
}

func historyBuiltin(stdout, stderr io.Writer, args []string, list *history.List) {
	if len(args) >= 2 && args[0] == "-r" {
		if err := list.ReadFromFile(args[1]); err != nil {
			fmt.Fprintln(stderr, historyErrorMessage(err))
		}
		return
	}
	if len(args) >= 2 && args[0] == "-w" {
		if err := list.WriteToFile(args[1]); err != nil {
			fmt.Fprintln(stderr, historyErrorMessage(err))
		}
		return
	}
	if len(args) >= 2 && args[0] == "-a" {
		if err := list.AppendToFile(args[1]); err != nil {
			fmt.Fprintln(stderr, historyErrorMessage(err))
		}
		return
	}

	limit := parseHistoryLimit(args)
	var entries []history.Entry
	if limit > 0 {
		entries = list.ListLast(limit)
	} else {
		entries = list.List()
	}
	history.WriteAll(stdout, entries)
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

func historyErrorMessage(err error) string {
	return fmt.Sprintf("history: %v", err)
}
