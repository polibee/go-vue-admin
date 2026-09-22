package controllers

import (
	"testing"

	"goravel/app/core/resource"
)

func TestFindManifestActionRequiresDeclaredAction(t *testing.T) {
	manifest := resource.Manifest{Actions: []resource.Action{
		{Name: "set-status", Kind: "user-status", Permission: "admin.users.manage"},
	}}
	action, ok := findManifestAction(manifest, "set-status")
	if !ok || action.Kind != "user-status" {
		t.Fatalf("expected declared action, got %+v, %v", action, ok)
	}
	if _, ok := findManifestAction(manifest, "unknown"); ok {
		t.Fatal("expected undeclared action to be rejected")
	}
}

func TestMergeActionResultsIncludesScopeFailures(t *testing.T) {
	result := mergeActionResults("set-status", 3, []actionFailureInput{
		{ID: 9, Code: "RESOURCE_NOT_FOUND"},
	}, actionResultInput{Succeeded: 2})
	if result.Requested != 3 || result.Succeeded != 2 || result.Failed != 1 {
		t.Fatalf("unexpected result: %+v", result)
	}
	if len(result.Failures) != 1 || result.Failures[0].ID != 9 {
		t.Fatalf("unexpected failures: %+v", result.Failures)
	}
}
