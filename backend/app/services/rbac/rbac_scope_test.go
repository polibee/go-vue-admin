package rbacservices

import (
	"testing"

	"goravel/app/core/resource"
)

func TestValidatePermissionScope(t *testing.T) {
	if err := validatePermissionScope(resource.DataScopeAll); err != nil {
		t.Fatalf("all scope should be valid: %v", err)
	}
	if err := validatePermissionScope(resource.DataScopeOwn); err != nil {
		t.Fatalf("own scope should be valid: %v", err)
	}
	if err := validatePermissionScope(resource.DataScope("department")); err == nil {
		t.Fatal("expected unknown scope to be rejected")
	}
}
