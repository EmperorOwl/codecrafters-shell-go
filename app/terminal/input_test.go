package terminal

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestReadLineRaw_upArrowRecallsHistory(t *testing.T) {
	handler := stubHistoryHandler{commands: []string{"echo hello", "echo world"}}
	input := bytes.NewBufferString("\x1b[A\r")
	var out bytes.Buffer

	got, eof, err := readLineRaw(bufio.NewReader(input), &out, nil, handler)
	if err != nil {
		t.Fatalf("readLineRaw() error = %v", err)
	}
	if eof {
		t.Fatal("readLineRaw() eof = true, want false")
	}
	if diff := cmp.Diff("echo world", got); diff != "" {
		t.Errorf("readLineRaw() line mismatch (-want +got):\n%s", diff)
	}
}

func TestReadLineRaw_upArrowTwiceRecallsEarlierHistory(t *testing.T) {
	handler := stubHistoryHandler{commands: []string{"echo hello", "echo world"}}
	input := bytes.NewBufferString("\x1b[A\x1b[A\r")
	var out bytes.Buffer

	got, _, err := readLineRaw(bufio.NewReader(input), &out, nil, handler)
	if err != nil {
		t.Fatalf("readLineRaw() error = %v", err)
	}
	if diff := cmp.Diff("echo hello", got); diff != "" {
		t.Errorf("readLineRaw() line mismatch (-want +got):\n%s", diff)
	}
}

func TestReadLineRaw_downArrowAfterUpRecallsNewerHistory(t *testing.T) {
	handler := stubHistoryHandler{commands: []string{"echo hello", "echo world"}}
	input := bytes.NewBufferString("\x1b[A\x1b[A\x1b[B\r")
	var out bytes.Buffer

	got, _, err := readLineRaw(bufio.NewReader(input), &out, nil, handler)
	if err != nil {
		t.Fatalf("readLineRaw() error = %v", err)
	}
	if diff := cmp.Diff("echo world", got); diff != "" {
		t.Errorf("readLineRaw() line mismatch (-want +got):\n%s", diff)
	}
}
