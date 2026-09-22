package openapi

import "testing"

func TestSpecCoversImplementedAdminContracts(t *testing.T) {
	spec := Spec()
	if spec["openapi"] != "3.0.3" {
		t.Fatalf("unexpected OpenAPI version: %v", spec["openapi"])
	}
	paths := spec["paths"].(map[string]any)
	for _, path := range []string{"/auth/login", "/auth/refresh", "/auth/logout-all", "/admin/registry", "/admin/search", "/admin/{resource}", "/admin/{resource}/export", "/admin/{resource}/actions/{action}", "/admin/{resource}/relations/{relation}/options", "/admin/{resource}/{id}/relations/{relation}", "/admin/{resource}/{id}", "/admin/overview", "/admin/audit-logs", "/admin/audit-logs/cleanup", "/admin/settings", "/admin/settings/{key}", "/admin/roles/{id}/permissions"} {
		if _, ok := paths[path]; !ok {
			t.Fatalf("missing contract path %s", path)
		}
	}
	if _, ok := paths["/admin/users/status"]; ok {
		t.Fatal("legacy user status action remains in contract")
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
	for _, schema := range []string{"ResourceManifest", "ResourceField", "ResourceOption", "ResourceAction", "ActionPayloadField", "ResourceRelation", "ResourceFormGroup", "ResourceDetailSection", "ResourceFieldDependency", "RelationOption", "RelationOptionList", "ResourceListMeta", "ActionRequest", "ActionResponse", "GlobalSearchResult", "GlobalSearchResponse", "DataScope", "AuditCleanupMode", "AuditCleanupRequest", "AuditCleanupResponse"} {
		if _, ok := schemas[schema]; !ok {
			t.Fatalf("missing resource schema %s", schema)
		}
	}
	cleanup := paths["/admin/audit-logs/cleanup"].(map[string]any)["post"].(map[string]any)
	cleanupBody := cleanup["requestBody"].(map[string]any)["content"].(map[string]any)["application/json"].(map[string]any)["schema"].(map[string]any)
	if cleanupBody["$ref"] != "#/components/schemas/AuditCleanupRequest" {
		t.Fatalf("unexpected audit cleanup request schema: %v", cleanupBody)
	}
	cleanupProperties := schemas["AuditCleanupRequest"].(map[string]any)["properties"].(map[string]any)
	if _, ok := cleanupProperties["mode"]; !ok {
		t.Fatal("audit cleanup contract must expose mode")
	}
	actionContract := paths["/admin/{resource}/actions/{action}"].(map[string]any)["post"].(map[string]any)
	if actionContract["operationId"] != "executeResourceAction" {
		t.Fatalf("unexpected resource action operation: %v", actionContract["operationId"])
	}
	status := spec["components"].(map[string]any)["schemas"].(map[string]any)["UserStatus"]
	if status == nil {
		t.Fatal("missing UserStatus schema")
	}
	rolePermissions := paths["/admin/roles/{id}/permissions"].(map[string]any)["put"].(map[string]any)
	requestBody := rolePermissions["requestBody"].(map[string]any)
	content := requestBody["content"].(map[string]any)["application/json"].(map[string]any)
	bodySchema := content["schema"].(map[string]any)
	properties := bodySchema["properties"].(map[string]any)
	if _, ok := properties["scopes"]; !ok {
		t.Fatal("role permission contract must expose scopes")
	}
	if _, ok := properties["fields"]; !ok {
		t.Fatal("role permission contract must expose fields")
	}
}
