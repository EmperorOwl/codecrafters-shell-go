package testutils

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestWantLines(t *testing.T) {
	tests := []struct {
		name  string
		lines []string
		want  string
	}{
		{name: "empty", lines: nil, want: ""},
		{name: "single line", lines: []string{"hello"}, want: "hello\n"},
		{name: "multiple lines", lines: []string{"a", "b"}, want: "a\nb\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if diff := cmp.Diff(tt.want, WantLines(tt.lines)); diff != "" {
				t.Errorf("WantLines() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestWriteFileIn(t *testing.T) {
	dir := t.TempDir()
	path := WriteFileIn(t, dir, "file.txt", "hello\n")

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if diff := cmp.Diff("hello\n", string(got)); diff != "" {
		t.Errorf("WriteFileIn() content mismatch (-want +got):\n%s", diff)
	}
}

func TestWriteTempFile(t *testing.T) {
	path := WriteTempFile(t, "file.txt", "temp\n")

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if diff := cmp.Diff("temp\n", string(got)); diff != "" {
		t.Errorf("WriteTempFile() content mismatch (-want +got):\n%s", diff)
	}
}

func TestOutputLines(t *testing.T) {
	tests := []struct {
		name string
		text string
		want []string
	}{
		{name: "empty", text: "", want: nil},
		{name: "single line", text: "hello\n", want: []string{"hello"}},
		{name: "multiple lines", text: "a\nb\n", want: []string{"a", "b"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if diff := cmp.Diff(tt.want, OutputLines(tt.text), cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("OutputLines() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
