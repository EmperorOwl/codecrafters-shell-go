package shell

import "path/filepath"

func FindExecutableInPath(command, pathEnv string) (string, bool) {
	for _, dir := range filepath.SplitList(pathEnv) {
		if dir == "" {
			continue
		}

		candidate := filepath.Join(dir, command)
		if path, ok := isExecutable(candidate); ok {
			return path, true
		}
	}

	return "", false
}
