package builtins

import (
	"fmt"
	"io"
	"os"
)

func CdErrorMessage(directory string) string {
	return "cd: " + directory + ": No such file or directory"
}

func ChangeDirectory(directory string) error {
	return os.Chdir(directory)
}

func ResolveDirectory(directory string) string {
	if directory == "~" {
		return os.Getenv("HOME")
	}
	return directory
}

func Cd(stderr io.Writer, directory string) {
	if directory == "" {
		return
	}
	target := ResolveDirectory(directory)
	if target == "" {
		fmt.Fprintln(stderr, CdErrorMessage(directory))
		return
	}
	if err := ChangeDirectory(target); err != nil {
		fmt.Fprintln(stderr, CdErrorMessage(directory))
	}
}

func init() {
	register("cd", cdBuiltin)
}

func cdBuiltin(ctx *Context, args []string) (bool, error) {
	directory := ""
	if len(args) > 0 {
		directory = args[0]
	}
	Cd(ctx.Stderr, directory)
	return false, nil
}
