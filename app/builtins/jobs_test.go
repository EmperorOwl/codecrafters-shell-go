package builtins

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestJobs(t *testing.T) {
	var out bytes.Buffer
	Jobs(&out)
	if diff := cmp.Diff("", out.String()); diff != "" {
		t.Errorf("Jobs() output mismatch (-want +got):\n%s", diff)
	}
}
