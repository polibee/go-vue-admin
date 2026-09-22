package controllers

import (
	"testing"

	adminactions "goravel/app/core/admin/actions"
	"goravel/app/core/resource"
)

func TestFindManifestActionRequiresDeclaredAction(t *testing.T) {
	manifest := resource.Manifest{Actions: []resource.Action{
		{Name: "set-status", Kind: "user-status", Permission: "admin.users.manage", Batch: true, Payload: "user-status"},
	}}
	action, ok := findManifestAction(manifest, "set-status")
	if !ok || action.Kind != "user-status" {
		t.Fatalf("expected declared action, got %+v, %v", action, ok)
	}
	if _, ok := findManifestAction(manifest, "unknown"); ok {
		t.Fatal("expected undeclared action to be rejected")
	}
}

func TestMergeActionResultsIncludesSkippedIDs(t *testing.T) {
	result := mergeActionResults("set-status", 3, []actionFailureInput{{ID: 9, Code: "OUT_OF_SCOPE"}}, actionResultInput{Succeeded: 1, Failures: []adminactions.Failure{{ID: 8, Code: "RESOURCE_NOT_FOUND"}}})
	if result.Requested != 3 || result.Succeeded != 1 || result.Failed != 1 || result.Skipped != 1 {
		t.Fatalf("unexpected result: %+v", result)
	}
	if len(result.Skips) != 1 || result.Skips[0].Code != "OUT_OF_SCOPE" {
		t.Fatalf("unexpected skips: %+v", result.Skips)
	}
}

func TestMergeActionResultsIncludesHandlerFailures(t *testing.T) {
	result := mergeActionResults("set-status", 3, nil, actionResultInput{Succeeded: 2, Failures: []adminactions.Failure{{ID: 9, Code: "RESOURCE_NOT_FOUND"}}})
	if result.Requested != 3 || result.Succeeded != 2 || result.Failed != 1 {
		t.Fatalf("unexpected result: %+v", result)
	}
	if len(result.Failures) != 1 || result.Failures[0].ID != 9 {
		t.Fatalf("unexpected failures: %+v", result.Failures)
	}
}
