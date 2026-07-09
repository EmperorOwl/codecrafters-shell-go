package files

import (
	"bufio"
	"os"
	"path/filepath"
	"slices"
	"strings"
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

// ReadLines returns the newline-delimited contents of path, one string per line.
func ReadLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

// WriteLines writes lines to path, one per line, with a trailing newline.
func WriteLines(path string, lines []string) error {
	content := strings.Join(lines, "\n") + "\n"
	return os.WriteFile(path, []byte(content), 0o644)
}
