package jobs

import "strings"

func reapJobs(jobList []Job, hasExited func(int) bool) (display []Job, remaining []Job) {
	for _, job := range jobList {
		if job.Status == "Running" && hasExited(job.PID) {
			doneJob := job
			doneJob.Status = "Done"
			doneJob.Command = strings.TrimSuffix(job.Command, " &")
			display = append(display, doneJob)
			continue
		}

		display = append(display, job)
		remaining = append(remaining, job)
	}
	return display, remaining
}
