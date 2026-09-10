package openapi

import "testing"

func TestDocumentIncludesResourceContract(t *testing.T) {
	document := BuildDocument()
	paths, ok := document["paths"].(map[string]any)
	if !ok {
		t.Fatal("openapi document paths are missing")
	}

	for _, path := range []string{"/api/resources/users", "/api/resources/users/{id}"} {
		if _, ok := paths[path]; !ok {
			t.Fatalf("openapi document is missing path %s", path)
		}
	}
	for _, path := range []string{"/api/settings", "/api/settings/{namespace}/{key}"} {
		if _, ok := paths[path]; !ok {
			t.Fatalf("openapi document is missing settings path %s", path)
		}
	}
	for _, path := range []string{"/api/media", "/api/media/{id}/preview", "/api/media/{id}", "/api/audit", "/api/audit/{id}"} {
		if _, ok := paths[path]; !ok {
			t.Fatalf("openapi document is missing media path %s", path)
		}
	}
	if _, ok := paths["/api/auth/bootstrap"]; !ok {
		t.Fatal("openapi document is missing the API bootstrap endpoint")
	}
	users := paths["/api/resources/users"].(map[string]any)
	for _, method := range []string{"get", "post"} {
		if _, ok := users[method]; !ok {
			t.Fatalf("users resource is missing %s operation", method)
		}
	}
	if operationID := users["get"].(map[string]any)["operationId"]; operationID != "listUsers" {
		t.Fatalf("unexpected users list operation id: %v", operationID)
	}
	demo := paths["/api/resources/demo"].(map[string]any)
	if operationID := demo["get"].(map[string]any)["operationId"]; operationID != "listDemoResources" {
		t.Fatalf("unexpected demo list operation id: %v", operationID)
	}
	item := paths["/api/resources/users/{id}"].(map[string]any)
	for _, method := range []string{"get", "put", "delete"} {
		if _, ok := item[method]; !ok {
			t.Fatalf("users item resource is missing %s operation", method)
		}
	}
}

func TestSchemaDocumentsExposeCoreResourceModels(t *testing.T) {
	schemas := SchemaDocuments()
	for _, name := range []string{"demo-resource", "user-resource", "role-resource", "permission-resource", "setting-resource", "media-resource", "audit-resource", "resource-envelope"} {
		if _, ok := schemas[name]; !ok {
			t.Fatalf("schema %q is missing", name)
		}
	}
}
