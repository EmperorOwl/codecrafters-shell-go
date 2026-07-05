package jobs

import (
	"bytes"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestJobManagerAdd(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(*JobManager)
		pid        int
		command    string
		wantNumber int
		wantJobs   []Job
	}{
		{
			name:       "first job gets number 1",
			setup:      func(*JobManager) {},
			pid:        42,
			command:    "sleep 10 &",
			wantNumber: 1,
			wantJobs: []Job{{
				Number:  1,
				PID:     42,
				Command: "sleep 10 &",
				Status:  StatusRunning,
			}},
		},
		{
			name: "second job increments number",
			setup: func(t *JobManager) {
				t.Add(42, "sleep 10 &")
			},
			pid:        99,
			command:    "sleep 5 &",
			wantNumber: 2,
			wantJobs: []Job{
				{Number: 1, PID: 42, Command: "sleep 10 &", Status: StatusRunning},
				{Number: 2, PID: 99, Command: "sleep 5 &", Status: StatusRunning},
			},
		},
		{
			name: "reuses 1 after table is empty",
			setup: func(t *JobManager) {
				t.Add(1, "sleep 1 &")
				t.Add(2, "sleep 2 &")
				t.MarkDone(1)
				t.MarkDone(2)
				t.ReapDone()
			},
			pid:        99,
			command:    "sleep 10 &",
			wantNumber: 1,
			wantJobs: []Job{{
				Number:  1,
				PID:     99,
				Command: "sleep 10 &",
				Status:  StatusRunning,
			}},
		},
		{
			name: "reuses 2 when job 1 is still running",
			setup: func(t *JobManager) {
				t.Add(1, "sleep 100 &")
				t.Add(2, "sleep 1 &")
				t.MarkDone(2)
				t.ReapDone()
			},
			pid:        99,
			command:    "sleep 10 &",
			wantNumber: 2,
			wantJobs: []Job{
				{Number: 1, PID: 1, Command: "sleep 100 &", Status: StatusRunning},
				{Number: 2, PID: 99, Command: "sleep 10 &", Status: StatusRunning},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var table JobManager
			tt.setup(&table)

			got := table.Add(tt.pid, tt.command)
			if got != tt.wantNumber {
				t.Errorf("Add() number = %d, want %d", got, tt.wantNumber)
			}

			table.mu.Lock()
			if diff := cmp.Diff(tt.wantJobs, table.jobs, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("jobs mismatch (-want +got):\n%s", diff)
			}
			table.mu.Unlock()
		})
	}
}

func TestJobManagerMarkDone(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(*JobManager)
		jobNumber int
		wantJobs  []Job
	}{
		{
			name: "marks job done and strips background suffix",
			setup: func(t *JobManager) {
				t.Add(42, "sleep 1 &")
			},
			jobNumber: 1,
			wantJobs: []Job{{
				Number:  1,
				PID:     42,
				Command: "sleep 1",
				Status:  StatusDone,
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var table JobManager
			tt.setup(&table)

			table.MarkDone(tt.jobNumber)

			table.mu.Lock()
			if diff := cmp.Diff(tt.wantJobs, table.jobs, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("jobs mismatch (-want +got):\n%s", diff)
			}
			table.mu.Unlock()
		})
	}
}

func TestJobManagerReapDone(t *testing.T) {
	tests := []struct {
		name          string
		setup         func(*JobManager)
		wantDone      []Job
		wantRemaining []Job
	}{
		{
			name: "returns no done jobs",
			setup: func(t *JobManager) {
				t.Add(1, "sleep 10 &")
			},
			wantRemaining: []Job{{
				Number:  1,
				PID:     1,
				Command: "sleep 10 &",
				Status:  StatusRunning,
			}},
		},
		{
			name: "returns done jobs and removes them",
			setup: func(t *JobManager) {
				t.Add(1, "cat /path/to/fifo &")
				t.MarkDone(1)
			},
			wantDone: []Job{{
				Number:  1,
				PID:     1,
				Command: "cat /path/to/fifo",
				Status:  StatusDone,
			}},
		},
		{
			name: "only reaps done jobs",
			setup: func(t *JobManager) {
				t.Add(1, "sleep 500 &")
				t.Add(2, "cat /path/to/fifo &")
				t.MarkDone(2)
			},
			wantDone: []Job{{
				Number:  2,
				PID:     2,
				Command: "cat /path/to/fifo",
				Status:  StatusDone,
			}},
			wantRemaining: []Job{{
				Number:  1,
				PID:     1,
				Command: "sleep 500 &",
				Status:  StatusRunning,
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var table JobManager
			tt.setup(&table)

			done := table.ReapDone()
			if diff := cmp.Diff(tt.wantDone, done, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("done mismatch (-want +got):\n%s", diff)
			}

			table.mu.Lock()
			remaining := append([]Job(nil), table.jobs...)
			table.mu.Unlock()
			if diff := cmp.Diff(tt.wantRemaining, remaining, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("remaining mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestWriteAll(t *testing.T) {
	tests := []struct {
		name      string
		jobs      []Job
		wantLines []string
	}{
		{
			name: "no jobs",
			jobs: nil,
		},
		{
			name: "one running job",
			jobs: []Job{{
				Number:  1,
				Command: "sleep 10 &",
				Status:  StatusRunning,
			}},
			wantLines: []string{
				"[1]+  Running                 sleep 10 &",
			},
		},
		{
			name: "two jobs mark current and previous",
			jobs: []Job{
				{Number: 1, Command: "sleep 10 &", Status: StatusRunning},
				{Number: 2, Command: "sleep 20 &", Status: StatusRunning},
			},
			wantLines: []string{
				"[1]-  Running                 sleep 10 &",
				"[2]+  Running                 sleep 20 &",
			},
		},
		{
			name: "three jobs mark current, previous, and space",
			jobs: []Job{
				{Number: 1, Command: "sleep 10 &", Status: StatusRunning},
				{Number: 2, Command: "sleep 20 &", Status: StatusRunning},
				{Number: 3, Command: "sleep 30 &", Status: StatusRunning},
			},
			wantLines: []string{
				"[1]   Running                 sleep 10 &",
				"[2]-  Running                 sleep 20 &",
				"[3]+  Running                 sleep 30 &",
			},
		},
		{
			name: "done job",
			jobs: []Job{{
				Number:  1,
				Command: "cat /path/to/fifo",
				Status:  StatusDone,
			}},
			wantLines: []string{
				"[1]+  Done                    cat /path/to/fifo",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out bytes.Buffer
			WriteAll(&out, tt.jobs)

			want := ""
			if len(tt.wantLines) > 0 {
				want = strings.Join(tt.wantLines, "\n") + "\n"
			}
			if diff := cmp.Diff(want, out.String()); diff != "" {
				t.Errorf("WriteAll() output mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
