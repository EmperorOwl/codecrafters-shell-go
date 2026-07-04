//go:build windows

package external

import (
	"os"
	"path/filepath"
	"strings"
)

var defaultPATHEXT = ".COM;.EXE;.BAT;.CMD;.VBS;.VBE;.JS;.JSE;.WSF;.WSH;.MSC"

func isExecutable(path string) (string, bool) {
	if ext := filepath.Ext(path); ext != "" {
		if !isPathext(ext) {
			return "", false
		}
		return statFile(path)
	}

	for _, ext := range pathextList() {
		if candidate, ok := statFile(path + strings.ToLower(ext)); ok {
			return candidate, true
		}
	}

	return "", false
}

func statFile(path string) (string, bool) {
	info, err := os.Stat(path)
	if err != nil || info.IsDir() {
		return "", false
	}
	return path, true
}

func isPathext(ext string) bool {
	for _, allowed := range pathextList() {
		if strings.EqualFold(ext, allowed) {
			return true
		}
	}
	return false
}

func pathextList() []string {
	pathext := os.Getenv("PATHEXT")
	if pathext == "" {
		pathext = defaultPATHEXT
	}
	return filepath.SplitList(pathext)
}
