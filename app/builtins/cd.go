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

func Cd(out io.Writer, directory string) {
	if directory == "" {
		return
	}
	if err := ChangeDirectory(directory); err != nil {
		fmt.Fprintln(out, CdErrorMessage(directory))
	}
}
