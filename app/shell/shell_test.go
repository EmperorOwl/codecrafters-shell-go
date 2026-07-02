package shell

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestPrintReapedJobs(t *testing.T) {
	shell := New()
	shell.jobs.Add(1, "cat /path/to/fifo &")
	shell.jobs.MarkDone(1)

	var out bytes.Buffer
	shell.PrintReapedJobs(&out)

	want := "[1]+  Done                    cat /path/to/fifo\n"
	if diff := cmp.Diff(want, out.String()); diff != "" {
		t.Errorf("PrintReapedJobs() output mismatch (-want +got):\n%s", diff)
	}

	if done := shell.jobs.ReapDone(); len(done) != 0 {
		t.Errorf("PrintReapedJobs() left %d done jobs in table", len(done))
	}
}
