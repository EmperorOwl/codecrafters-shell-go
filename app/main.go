package main

import (
	"os"

	"github.com/codecrafters-io/shell-starter-go/app/shell"
)

func main() {
	s := shell.New(os.Stdin, os.Stdout, os.Stderr)
	if err := s.Run(); err != nil {
		os.Exit(1)
	}
}
