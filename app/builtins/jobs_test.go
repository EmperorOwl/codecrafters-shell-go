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
		name      string
		jobs      []jobs.Job
		wantLines []string
		wantJobs  []jobs.Job
	}{
		{
			name:     "no jobs",
			jobs:     nil,
			wantJobs: nil,
		},
		{
			name: "one running job",
			jobs: []jobs.Job{{
				Number:  1,
				Command: "sleep 10 &",
				Status:  jobs.StatusRunning,
			}},
			wantLines: []string{
				"[1]+  Running                 sleep 10 &",
			},
			wantJobs: []jobs.Job{{
				Number:  1,
				Command: "sleep 10 &",
				Status:  jobs.StatusRunning,
			}},
		},
		{
			name: "done job formatting",
			jobs: []jobs.Job{{
				Number:  1,
				PID:     1,
				Command: "cat /path/to/fifo",
				Status:  jobs.StatusDone,
			}},
			wantLines: []string{
				"[1]+  Done                    cat /path/to/fifo",
			},
			wantJobs: []jobs.Job{{
				Number:  1,
				PID:     1,
				Command: "cat /path/to/fifo",
				Status:  jobs.StatusDone,
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out bytes.Buffer
			jobList := tt.jobs
			Jobs(&out, jobList)

			want := ""
			if len(tt.wantLines) > 0 {
				want = strings.Join(tt.wantLines, "\n") + "\n"
			}
			if diff := cmp.Diff(want, out.String()); diff != "" {
				t.Errorf("Jobs() output mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tt.wantJobs, jobList, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("Jobs() job list mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
