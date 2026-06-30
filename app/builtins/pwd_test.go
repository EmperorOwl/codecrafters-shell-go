package builtins

import (
	"bytes"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestPwd(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd() error = %v", err)
	}

	var out bytes.Buffer
	Pwd(&out)
	want := cwd + "\n"
	if diff := cmp.Diff(want, out.String()); diff != "" {
		t.Errorf("Pwd() output mismatch (-want +got):\n%s", diff)
	}
}
