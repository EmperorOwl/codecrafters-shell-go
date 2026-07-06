package completer

import (
	"os"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/files"
)

func (c *Completer) completeFile(buffer string) (string, []string) {
	return c.completeArgument(buffer, c.fileCandidates(buffer))
}

func (c *Completer) fileCandidates(buffer string) []string {
	lastSpace := strings.LastIndex(buffer, " ")
	if lastSpace < 0 {
		return nil
	}

	token := buffer[lastSpace+1:]
	dirPath, _ := splitDirToken(token)

	cwd, err := os.Getwd()
	if err != nil {
		return nil
	}
	return files.ListInDir(cwd, dirPath)
}
