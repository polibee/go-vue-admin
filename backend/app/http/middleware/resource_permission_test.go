package middleware

import (
	"testing"

	"goravel/app/core/resource"
)

func TestResourceViewPermissionsAreUnique(t *testing.T) {
	permissions := resourceViewPermissions([]resource.Manifest{
		{Permissions: []string{"admin.orders.view", "admin.shared.view"}},
		{Permissions: []string{"admin.orders.view", "admin.users.view"}},
	})
	if len(permissions) != 3 || permissions[0] != "admin.orders.view" || permissions[1] != "admin.shared.view" || permissions[2] != "admin.users.view" {
		t.Fatalf("unexpected resource permissions: %#v", permissions)
	}
}
