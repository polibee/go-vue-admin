package controllers

import "testing"

func TestResourceRecordLoadedRequiresPositiveID(t *testing.T) {
	if resourceRecordLoaded(0) {
		t.Fatal("expected zero ID to represent a missing resource")
	}
	if !resourceRecordLoaded(1) {
		t.Fatal("expected positive ID to represent a loaded resource")
	}
}
