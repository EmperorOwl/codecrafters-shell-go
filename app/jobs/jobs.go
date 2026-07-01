package jobs

import (
	"fmt"
	"io"
	"strings"
	"sync"
)

const (
	StatusRunning = "Running"
	StatusDone    = "Done"
)

type Job struct {
	Number  int
	PID     int
	Command string
	Status  string
}

type JobTable struct {
	mu   sync.Mutex
	jobs []Job
}

// nextJobNumberLocked returns the next job number to assign. The table must
// already be locked. Empty table yields 1; otherwise max existing number + 1.
func (t *JobTable) nextJobNumberLocked() int {
	if len(t.jobs) == 0 {
		return 1
	}
	maxNumber := 0
	for _, job := range t.jobs {
		if job.Number > maxNumber {
			maxNumber = job.Number
		}
	}
	return maxNumber + 1
}

// Add registers a new background job and returns its job number.
func (t *JobTable) Add(pid int, command string) int {
	t.mu.Lock()
	defer t.mu.Unlock()

	jobNumber := t.nextJobNumberLocked()
	job := Job{
		Number:  jobNumber,
		PID:     pid,
		Command: command,
		Status:  StatusRunning,
	}
	t.jobs = append(t.jobs, job)
	return job.Number
}

// MarkDone marks the given job as finished and strips the trailing " &".
func (t *JobTable) MarkDone(jobNumber int) {
	t.mu.Lock()
	defer t.mu.Unlock()

	for i, job := range t.jobs {
		if job.Number != jobNumber {
			continue
		}
		t.jobs[i].Status = StatusDone
		t.jobs[i].Command = strings.TrimSuffix(job.Command, " &")
		return
	}
}

// ReapDone removes finished jobs from the table and returns them.
func (t *JobTable) ReapDone() []Job {
	t.mu.Lock()
	defer t.mu.Unlock()

	var done []Job
	remaining := t.jobs[:0]
	for _, job := range t.jobs {
		if job.Status == StatusDone {
			done = append(done, job)
		} else {
			remaining = append(remaining, job)
		}
	}
	t.jobs = remaining
	return done
}

// ListForDisplay returns a snapshot of all jobs currently in the table.
func (t *JobTable) ListForDisplay() []Job {
	t.mu.Lock()
	defer t.mu.Unlock()

	display := make([]Job, len(t.jobs))
	copy(display, t.jobs)
	return display
}

// WriteAll prints each job on its own line using bash-style formatting.
func WriteAll(out io.Writer, jobList []Job) {
	for i, job := range jobList {
		fmt.Fprintln(out, formatLine(job, i, len(jobList)))
	}
}

// markerForIndex returns the job marker for the given position in a listing.
func markerForIndex(index, count int) string {
	switch {
	case index == count-1:
		return "+" // current (most recently started) job
	case index == count-2:
		return "-" // previous job
	default:
		return " " // older job
	}
}

// formatLine builds a single jobs-listing line for the given job.
func formatLine(job Job, index, count int) string {
	marker := markerForIndex(index, count)

	status := job.Status
	if len(status) < 24 {
		status += strings.Repeat(" ", 24-len(status))
	}

	return fmt.Sprintf("[%d]%s  %s%s", job.Number, marker, status, job.Command)
}
