package main

import (
	"os"

	"github.com/codecrafters-io/shell-starter-go/app/shell"
)

func main() {
	s := shell.New()
	if err := s.Run(os.Stdin, os.Stdout, os.Stderr); err != nil {
		os.Exit(1)
	}
}
