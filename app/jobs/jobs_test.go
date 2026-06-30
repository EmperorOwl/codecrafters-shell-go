package jobs

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestAddJob(t *testing.T) {
	var jobList []Job
	nextID := 0

	number := AddJob(&jobList, &nextID, 42, "sleep 10 &")
	if number != 1 {
		t.Fatalf("AddJob() number = %d, want 1", number)
	}
	if nextID != 1 {
		t.Fatalf("nextID = %d, want 1", nextID)
	}
	want := []Job{{
		Number:  1,
		PID:     42,
		Command: "sleep 10 &",
		Status:  "Running",
	}}
	if diff := cmp.Diff(want, jobList); diff != "" {
		t.Errorf("job list mismatch (-want +got):\n%s", diff)
	}

	AddJob(&jobList, &nextID, 99, "sleep 5 &")
	if diff := cmp.Diff(2, jobList[1].Number); diff != "" {
		t.Errorf("second job number mismatch (-want +got):\n%s", diff)
	}
}

func TestWriteAll(t *testing.T) {
	tests := []struct {
		name string
		jobs []Job
		want string
	}{
		{
			name: "no jobs",
			jobs: nil,
			want: "",
		},
		{
			name: "one running job",
			jobs: []Job{{
				Number:  1,
				Command: "sleep 10 &",
				Status:  "Running",
			}},
			want: "[1]+  Running                 sleep 10 &\n",
		},
		{
			name: "most recent job gets plus marker",
			jobs: []Job{
				{Number: 1, Command: "sleep 5 &", Status: "Running"},
				{Number: 2, Command: "sleep 10 &", Status: "Running"},
			},
			want: "[1]   Running                 sleep 5 &\n[2]+  Running                 sleep 10 &\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out bytes.Buffer
			WriteAll(&out, tt.jobs)
			if diff := cmp.Diff(tt.want, out.String()); diff != "" {
				t.Errorf("WriteAll() output mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
