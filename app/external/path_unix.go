//go:build !windows

package external

import "os"

func isExecutable(path string) (string, bool) {
	info, err := os.Stat(path)
	if err != nil {
		return "", false
	}
	if info.IsDir() {
		return "", false
	}
	if info.Mode()&0111 == 0 {
		return "", false
	}

	return path, true
}
