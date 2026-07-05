package shell

import (
	"io"
	"strings"
	"testing"

	"github.com/codecrafters-io/shell-starter-go/app/jobs"
	"github.com/google/go-cmp/cmp"
)

func TestReapDoneJobs(t *testing.T) {
	shell := New(strings.NewReader(""), io.Discard, io.Discard)
	shell.jobManager.Add(1, "cat /path/to/fifo &")
	shell.jobManager.MarkDone(1)

	got := jobs.FormatLines(shell.jobManager.ReapDone())
	want := []string{"[1]+  Done                    cat /path/to/fifo"}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("FormatLines(ReapDone()) mismatch (-want +got):\n%s", diff)
	}

	if done := shell.jobManager.ReapDone(); len(done) != 0 {
		t.Errorf("ReapDone() left %d done jobs in table", len(done))
	}
}
