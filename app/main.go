package main

import (
	"os"

	"github.com/codecrafters-io/shell-starter-go/app/shell"
)

func main() {
	shell.WritePrompt(os.Stdout)
}
