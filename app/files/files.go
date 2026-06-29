package files

import (
	"os"
	"slices"
)

// ListInCurrentDir returns sorted file and directory names in the current working directory.
// Directory names include a trailing slash.
func ListInCurrentDir() []string {
	return ListInDir("")
}

// ListInDir returns sorted file and directory names in dir relative to the current working directory.
// An empty dir lists the current working directory. Directory names include a trailing slash.
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
		name := entry.Name()
		if entry.IsDir() {
			name += "/"
		}
		result = append(result, name)
	}
	slices.Sort(result)
	return result
}
