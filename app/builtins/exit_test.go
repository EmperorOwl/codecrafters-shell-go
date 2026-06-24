package builtins

import "testing"

func TestExit(t *testing.T) {
	if !Exit() {
		t.Error("Exit() = false, want true")
	}
}
