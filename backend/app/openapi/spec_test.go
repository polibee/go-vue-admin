package openapi

import "testing"

func TestSpecCoversImplementedAdminContracts(t *testing.T) {
	spec := Spec()
	if spec["openapi"] != "3.0.3" {
		t.Fatalf("unexpected OpenAPI version: %v", spec["openapi"])
	}
	paths := spec["paths"].(map[string]any)
	for _, path := range []string{"/auth/login", "/auth/refresh", "/auth/logout-all", "/admin/resources", "/admin/resources/{resource}", "/admin/overview", "/admin/audit-logs", "/admin/settings", "/admin/settings/{key}", "/admin/users/status"} {
		if _, ok := paths[path]; !ok {
			t.Fatalf("missing contract path %s", path)
		}
	}
	status := spec["components"].(map[string]any)["schemas"].(map[string]any)["UserStatus"]
	if status == nil {
		t.Fatal("missing UserStatus schema")
	}
}
