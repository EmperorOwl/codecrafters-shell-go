package jobs

import (
	"fmt"
	"io"
	"strings"
)

type Job struct {
	Number  int
	PID     int
	Command string
	Status  string
}

func AddJob(jobs *[]Job, nextID *int, pid int, command string) int {
	*nextID++
	job := Job{
		Number:  *nextID,
		PID:     pid,
		Command: command,
		Status:  "Running",
	}
	*jobs = append(*jobs, job)
	return job.Number
}

func WriteAll(out io.Writer, jobs []Job) {
	for i, job := range jobs {
		fmt.Fprintln(out, formatLine(job, i == len(jobs)-1))
	}
}

func formatLine(job Job, isCurrent bool) string {
	marker := " "
	if isCurrent {
		marker = "+"
	}

	status := job.Status
	if len(status) < 24 {
		status += strings.Repeat(" ", 24-len(status))
	}

	return fmt.Sprintf("[%d]%s  %s%s", job.Number, marker, status, job.Command)
}
