package controllers

import "testing"

func TestLoginAllowedByStatus(t *testing.T) {
	if !loginAllowedForStatus("active") {
		t.Fatal("active users should be allowed to log in")
	}
	if loginAllowedForStatus("disabled") || loginAllowedForStatus("locked") {
		t.Fatal("disabled and locked users should not be allowed to log in")
	}
}
