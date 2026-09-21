package controllers

import "testing"

func TestNormalizeUserStatusFilter(t *testing.T) {
	if got := normalizeUserStatusFilter("disabled"); got != "disabled" {
		t.Fatalf("got %q, want disabled", got)
	}
	if got := normalizeUserStatusFilter("pending"); got != "" {
		t.Fatalf("invalid status filter = %q, want empty", got)
	}
}
