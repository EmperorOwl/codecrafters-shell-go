package files

import (
	"os"
	"slices"
)

// ListInCurrentDir returns sorted filenames (not directories) in the current working directory.
func ListInCurrentDir() []string {
	return ListInDir("")
}

// ListInDir returns sorted filenames (not directories) in dir relative to the current working directory.
// An empty dir lists the current working directory.
func ListInDir(dir string) []string {
	if dir == "" {
		dir = "."
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	var result []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		result = append(result, entry.Name())
	}
	slices.Sort(result)
	return result
}
