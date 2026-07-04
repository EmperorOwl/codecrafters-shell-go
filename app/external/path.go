package external

import (
	"os"
	"path/filepath"
	"slices"
)

// FindExecutableInPath searches PATH for an executable named fullName.
// It returns the resolved file path and true when found; missing PATH directories are skipped.
func FindExecutableInPath(fullName string) (string, bool) {
	for _, dir := range filepath.SplitList(os.Getenv("PATH")) {
		if dir == "" {
			continue
		}

		candidate := filepath.Join(dir, fullName)
		if path, ok := isExecutable(candidate); ok {
			return path, true
		}
	}

	return "", false
}

// FindAllExecutablesInPath returns sorted executable basenames from PATH.
// Missing PATH directories are skipped.
func FindAllExecutablesInPath() []string {
	seen := make(map[string]struct{})
	var matches []string

	for _, dir := range filepath.SplitList(os.Getenv("PATH")) {
		if dir == "" {
			continue
		}

		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			name := entry.Name()
			if _, ok := isExecutable(filepath.Join(dir, name)); !ok {
				continue
			}

			if _, ok := seen[name]; ok {
				continue
			}
			seen[name] = struct{}{}
			matches = append(matches, name)
		}
	}

	slices.Sort(matches)
	return matches
}
