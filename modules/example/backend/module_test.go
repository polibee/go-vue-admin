package example

import "testing"

func TestModuleIdentity(t *testing.T) {
	if got := (Module{}).Name(); got != "example" {
		t.Fatalf("module name = %q", got)
	}
}
