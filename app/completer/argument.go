package completer

import (
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/completion"
)

func (c *Completer) completeArgument(buffer string, candidates []string) (string, []string) {
	lastSpace := strings.LastIndex(buffer, " ")
	if lastSpace < 0 {
		return buffer, nil
	}

	bufferPrefix := buffer[:lastSpace+1]
	token := buffer[lastSpace+1:]
	dirPath, prefix := splitDirToken(token)

	token, listings, unique := completion.Complete(prefix, candidates)
	if len(listings) > 0 {
		return buffer, listings
	}
	if token != prefix {
		suffix := ""
		if unique {
			suffix = completionSuffix(token)
		}
		return bufferPrefix + dirPath + token + suffix, nil
	}
	return buffer, nil
}

func splitDirToken(token string) (dirPath, prefix string) {
	dirPath = ""
	prefix = token

	if idx := strings.LastIndex(token, "/"); idx >= 0 {
		dirPath = token[:idx+1]
		prefix = token[idx+1:]
	}
	return dirPath, prefix
}

func completionSuffix(token string) string {
	if strings.HasSuffix(token, "/") {
		return ""
	}
	return " "
}

// CompleteArgument completes the last argument token for tests and direct use.
func (c *Completer) CompleteArgument(buffer string, candidates []string) (string, []string) {
	return c.completeArgument(buffer, candidates)
}
