package files

import (
	"os"
	"path/filepath"
	"slices"
)

// ListInDir returns sorted file and directory names in dir relative to base.
// An empty dir lists base. Directory names include a trailing slash.
func ListInDir(base, dir string) []string {
	path := base
	if dir != "" {
		path = filepath.Join(base, filepath.FromSlash(dir))
	}

	entries, err := os.ReadDir(path)
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
