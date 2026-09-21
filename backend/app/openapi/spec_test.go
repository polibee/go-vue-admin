package openapi

import "testing"

func TestSpecCoversImplementedAdminContracts(t *testing.T) {
	spec := Spec()
	if spec["openapi"] != "3.0.3" {
		t.Fatalf("unexpected OpenAPI version: %v", spec["openapi"])
	}
	paths := spec["paths"].(map[string]any)
	for _, path := range []string{"/auth/login", "/auth/refresh", "/auth/logout-all", "/admin/resources", "/admin/resources/{resource}", "/admin/resources/{resource}/{id}", "/admin/overview", "/admin/audit-logs", "/admin/settings", "/admin/settings/{key}", "/admin/users/status"} {
		if _, ok := paths[path]; !ok {
			t.Fatalf("missing contract path %s", path)
		}
	}
	resourceCollection := paths["/admin/resources/{resource}"].(map[string]any)
	for _, method := range []string{"get", "post"} {
		if _, ok := resourceCollection[method]; !ok {
			t.Fatalf("missing resource collection method %s", method)
		}
	}
	resourceItem := paths["/admin/resources/{resource}/{id}"].(map[string]any)
	for _, method := range []string{"get", "put", "delete"} {
		if _, ok := resourceItem[method]; !ok {
			t.Fatalf("missing resource item method %s", method)
		}
	}
	schemas := spec["components"].(map[string]any)["schemas"].(map[string]any)
	for _, schema := range []string{"ResourceManifest", "ResourceField", "ResourceOption", "ResourceAction"} {
		if _, ok := schemas[schema]; !ok {
			t.Fatalf("missing resource schema %s", schema)
		}
	}
	status := spec["components"].(map[string]any)["schemas"].(map[string]any)["UserStatus"]
	if status == nil {
		t.Fatal("missing UserStatus schema")
	}
}
