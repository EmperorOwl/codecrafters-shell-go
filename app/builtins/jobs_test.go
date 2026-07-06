package builtins

import (
	"bytes"
	"strings"
	"testing"

	"github.com/codecrafters-io/shell-starter-go/app/jobs"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestJobs(t *testing.T) {
	tests := []struct {
		name          string
		setup         func(*jobs.JobTable)
		wantLines     []string
		wantRemaining []jobs.Job
	}{
		{
			name:  "no jobs",
			setup: func(*jobs.JobTable) {},
		},
		{
			name: "one running job",
			setup: func(t *jobs.JobTable) {
				t.Add(1, "sleep 10 &")
			},
			wantLines: []string{
				"[1]+  Running                 sleep 10 &",
			},
			wantRemaining: []jobs.Job{{
				Number:  1,
				PID:     1,
				Command: "sleep 10 &",
				Status:  jobs.StatusRunning,
			}},
		},
		{
			name: "done job is reaped",
			setup: func(t *jobs.JobTable) {
				t.Add(1, "cat /path/to/fifo &")
				t.MarkDone(1)
			},
			wantLines: []string{
				"[1]+  Done                    cat /path/to/fifo",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var table jobs.JobTable
			tt.setup(&table)

			var out bytes.Buffer
			Jobs(&out, &table)

			want := ""
			if len(tt.wantLines) > 0 {
				want = strings.Join(tt.wantLines, "\n") + "\n"
			}
			if diff := cmp.Diff(want, out.String()); diff != "" {
				t.Errorf("Jobs() output mismatch (-want +got):\n%s", diff)
			}

			remaining := table.List()
			if diff := cmp.Diff(tt.wantRemaining, remaining, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("Jobs() remaining jobs mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
