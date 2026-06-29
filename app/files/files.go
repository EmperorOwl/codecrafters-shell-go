package files

import (
	"os"
	"slices"
)

// ListInCurrentDir returns sorted filenames (not directories) in the current working directory.
func ListInCurrentDir() []string {
	entries, err := os.ReadDir(".")
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
