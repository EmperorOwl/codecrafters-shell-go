package builtins

import (
	"fmt"
	"io"
	"os"
)

func init() {
	register("cd", cdHandler)
}

func cdHandler(ctx *Context, args []string) (bool, error) {
	directory := ""
	if len(args) > 0 {
		directory = args[0]
	}
	cdBuiltin(ctx.Stderr, directory)
	return false, nil
}

func cdBuiltin(stderr io.Writer, directory string) {
	if directory == "" {
		return
	}

	target := directory
	if directory == "~" {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintln(stderr, cdErrorMessage(directory))
			return
		}
		target = home
	}
	if err := os.Chdir(target); err != nil {
		fmt.Fprintln(stderr, cdErrorMessage(directory))
	}
}

func cdErrorMessage(directory string) string {
	return "cd: " + directory + ": No such file or directory"
}
