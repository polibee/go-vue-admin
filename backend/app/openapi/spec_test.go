package openapi

import "testing"

func TestSpecCoversImplementedAdminContracts(t *testing.T) {
	spec := Spec()
	if spec["openapi"] != "3.0.3" {
		t.Fatalf("unexpected OpenAPI version: %v", spec["openapi"])
	}
	paths := spec["paths"].(map[string]any)
	for _, path := range []string{"/auth/login", "/auth/refresh", "/auth/logout-all", "/admin/registry", "/admin/search", "/admin/{resource}", "/admin/{resource}/export", "/admin/{resource}/{id}", "/admin/overview", "/admin/audit-logs", "/admin/settings", "/admin/settings/{key}", "/admin/users/status"} {
		if _, ok := paths[path]; !ok {
			t.Fatalf("missing contract path %s", path)
		}
	}
	if _, ok := paths["/admin/resources/{resource}"]; ok {
		t.Fatal("legacy resources route remains in contract")
	}
	resourceCollection := paths["/admin/{resource}"].(map[string]any)
	for _, method := range []string{"get", "post"} {
		if _, ok := resourceCollection[method]; !ok {
			t.Fatalf("missing resource collection method %s", method)
		}
	}
	exportContract := paths["/admin/{resource}/export"].(map[string]any)
	exportResponse := exportContract["responses"].(map[string]any)["200"].(map[string]any)
	if _, ok := exportResponse["content"].(map[string]any)["text/csv"]; !ok {
		t.Fatal("resource export must return text/csv")
	}
	resourceItem := paths["/admin/{resource}/{id}"].(map[string]any)
	for _, method := range []string{"get", "put", "delete"} {
		if _, ok := resourceItem[method]; !ok {
			t.Fatalf("missing resource item method %s", method)
		}
	}
	schemas := spec["components"].(map[string]any)["schemas"].(map[string]any)
	for _, schema := range []string{"ResourceManifest", "ResourceField", "ResourceOption", "ResourceAction", "GlobalSearchResult", "GlobalSearchResponse"} {
		if _, ok := schemas[schema]; !ok {
			t.Fatalf("missing resource schema %s", schema)
		}
	}
	status := spec["components"].(map[string]any)["schemas"].(map[string]any)["UserStatus"]
	if status == nil {
		t.Fatal("missing UserStatus schema")
	}
}
