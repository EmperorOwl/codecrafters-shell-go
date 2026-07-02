package main

import (
	"os"

	"github.com/codecrafters-io/shell-starter-go/app/shell"
	"github.com/codecrafters-io/shell-starter-go/app/terminal"
)

func main() {
	s := shell.New()
	t := terminal.New(s, os.Stdin, os.Stdout, os.Stderr)
	if err := t.Run(); err != nil {
		os.Exit(1)
	}
}
