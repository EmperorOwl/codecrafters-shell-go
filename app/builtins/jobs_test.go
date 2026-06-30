package builtins

import (
	"bytes"
	"testing"

	"github.com/codecrafters-io/shell-starter-go/app/jobs"
	"github.com/google/go-cmp/cmp"
)

func TestJobs(t *testing.T) {
	tests := []struct {
		name string
		jobs []jobs.Job
		want string
	}{
		{
			name: "no jobs",
			jobs: nil,
			want: "",
		},
		{
			name: "one running job",
			jobs: []jobs.Job{{
				Number:  1,
				Command: "sleep 10 &",
				Status:  "Running",
			}},
			want: "[1]+  Running                 sleep 10 &\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out bytes.Buffer
			Jobs(&out, tt.jobs)
			if diff := cmp.Diff(tt.want, out.String()); diff != "" {
				t.Errorf("Jobs() output mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
