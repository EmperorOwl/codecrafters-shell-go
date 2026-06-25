package parser

import (
	"reflect"
	"testing"
)

func TestParseRedirect(t *testing.T) {
	tests := []struct {
		name           string
		tokens         []string
		wantFields     []string
		wantStdoutPath string
		wantStderrPath string
	}{
		{
			name:       "no redirect",
			tokens:     []string{"echo", "hello"},
			wantFields: []string{"echo", "hello"},
		},
		{
			name:           "stdout redirect",
			tokens:         []string{"echo", "hello", ">", "output.txt"},
			wantFields:     []string{"echo", "hello"},
			wantStdoutPath: "output.txt",
		},
		{
			name:           "explicit stdout redirect",
			tokens:         []string{"echo", "Hello", "James", "1>", "/tmp/foo/foo.md"},
			wantFields:     []string{"echo", "Hello", "James"},
			wantStdoutPath: "/tmp/foo/foo.md",
		},
		{
			name:           "external command with redirect",
			tokens:         []string{"ls", "/tmp/baz", ">", "/tmp/foo/baz.md"},
			wantFields:     []string{"ls", "/tmp/baz"},
			wantStdoutPath: "/tmp/foo/baz.md",
		},
		{
			name:           "redirect with multiple arguments",
			tokens:         []string{"cat", "/tmp/baz/blueberry", "nonexistent", "1>", "/tmp/foo/quz.md"},
			wantFields:     []string{"cat", "/tmp/baz/blueberry", "nonexistent"},
			wantStdoutPath: "/tmp/foo/quz.md",
		},
		{
			name:       "greater than inside token is not redirect",
			tokens:     []string{"echo", "a>b"},
			wantFields: []string{"echo", "a>b"},
		},
		{
			name:           "stderr redirect",
			tokens:         []string{"ls", "nonexistent", "2>", "/tmp/quz/baz.md"},
			wantFields:     []string{"ls", "nonexistent"},
			wantStderrPath: "/tmp/quz/baz.md",
		},
		{
			name:           "stderr redirect with stdout output",
			tokens:         []string{"cat", "/tmp/bar/pear", "nonexistent", "2>", "/tmp/quz/quz.md"},
			wantFields:     []string{"cat", "/tmp/bar/pear", "nonexistent"},
			wantStderrPath: "/tmp/quz/quz.md",
		},
		{
			name:           "echo with stderr redirect leaves args intact",
			tokens:         []string{"echo", "Maria", "file", "cannot", "be", "found", "2>", "/tmp/quz/foo.md"},
			wantFields:     []string{"echo", "Maria", "file", "cannot", "be", "found"},
			wantStderrPath: "/tmp/quz/foo.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFields, gotStdoutPath, gotStderrPath := ParseRedirect(tt.tokens)
			if !reflect.DeepEqual(gotFields, tt.wantFields) {
				t.Errorf("ParseRedirect() fields = %v, want %v", gotFields, tt.wantFields)
			}
			if gotStdoutPath != tt.wantStdoutPath {
				t.Errorf("ParseRedirect() stdoutPath = %q, want %q", gotStdoutPath, tt.wantStdoutPath)
			}
			if gotStderrPath != tt.wantStderrPath {
				t.Errorf("ParseRedirect() stderrPath = %q, want %q", gotStderrPath, tt.wantStderrPath)
			}
		})
	}
}
