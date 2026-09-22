package rbacservices

import (
	"testing"

	"goravel/app/core/resource"
)

func TestMergeDataScopesPrefersAll(t *testing.T) {
	if got := mergeDataScopes([]resource.DataScope{resource.DataScopeOwn, resource.DataScopeAll}); got != resource.DataScopeAll {
		t.Fatalf("expected all scope, got %q", got)
	}
	if got := mergeDataScopes([]resource.DataScope{resource.DataScopeOwn}); got != resource.DataScopeOwn {
		t.Fatalf("expected own scope, got %q", got)
	}
}
