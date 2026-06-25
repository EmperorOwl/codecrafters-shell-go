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
