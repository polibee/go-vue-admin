package rbac

import "testing"

func TestHasAnyPermission(t *testing.T) {
	permissions := map[string]struct{}{
		"admin.users.view":   {},
		"admin.roles.manage": {},
	}
	if !HasAnyPermission(permissions, "admin.roles.manage", "admin.permissions.manage") {
		t.Fatal("expected one matching permission to authorize the request")
	}
	if HasAnyPermission(permissions, "admin.permissions.manage") {
		t.Fatal("unexpectedly authorized a missing permission")
	}
}

func TestValidateRoleIDs(t *testing.T) {
	if err := ValidateRoleIDs([]int64{1, 2}); err != nil {
		t.Fatalf("expected valid role ids, got %v", err)
	}
	if err := ValidateRoleIDs([]int64{1, 1}); err == nil {
		t.Fatal("expected duplicate role ids to be rejected")
	}
}

func TestValidateRoleInput(t *testing.T) {
	if err := ValidateRoleInput("", "Administrator"); err == nil {
		t.Fatal("expected an empty role name to be rejected")
	}
	if err := ValidateRoleInput("admin", "Administrator"); err != nil {
		t.Fatalf("expected a valid role input, got %v", err)
	}
}
