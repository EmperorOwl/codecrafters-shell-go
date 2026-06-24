package shell

import (
	"os"
	"path/filepath"
)

func FindExecutableInPath(command, pathEnv string) (string, bool) {
	for _, dir := range filepath.SplitList(pathEnv) {
		if dir == "" {
			continue
		}

		candidate := filepath.Join(dir, command)
		info, err := os.Stat(candidate)
		if err != nil {
			continue
		}
		if info.IsDir() {
			continue
		}
		if info.Mode()&0111 == 0 {
			continue
		}

		return candidate, true
	}

	return "", false
}
