package builtins

import (
	"bytes"
	"os"
	"testing"
)

func TestPwd(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd() error = %v", err)
	}

	var out bytes.Buffer
	Pwd(&out)
	want := cwd + "\n"
	if got := out.String(); got != want {
		t.Errorf("Pwd() output = %q, want %q", got, want)
	}
}
