package rbacservices

import (
	"errors"
	"testing"

	"goravel/app/core/resource"
)

func TestFieldPermissionServiceProjectsReadableAndExportFields(t *testing.T) {
	manifest := resource.Manifest{
		Name: "users", Label: "Users", Route: "/users",
		Fields: []resource.Field{
			{Name: "name", Label: "Name", Type: "text"},
			{Name: "password", Label: "Password", Type: "password", Visible: false, Readable: false, Writable: false, Sensitive: true, PolicyConfigured: true},
			{Name: "token", Label: "Token", Type: "text", Visible: true, Readable: true, Writable: true, Sensitive: true, PolicyConfigured: true},
		},
	}
	service := NewFieldPermissionService()
	policies := service.ManifestPolicies(manifest)

	readable := service.ReadableFields(manifest, policies, false)
	if len(readable) != 2 || readable[0].Name != "name" || readable[1].Name != "token" {
		t.Fatalf("unexpected readable fields: %+v", readable)
	}
	exportable := service.ReadableFields(manifest, policies, true)
	if len(exportable) != 1 || exportable[0].Name != "name" {
		t.Fatalf("unexpected export fields: %+v", exportable)
	}
}

func TestFieldPermissionServiceRejectsSensitiveQueriesAndReadOnlyWrites(t *testing.T) {
	manifest := resource.Manifest{
		Name: "users", Label: "Users", Route: "/users",
		Fields: []resource.Field{
			{Name: "name", Label: "Name", Type: "text"},
			{Name: "token", Label: "Token", Type: "text", Sensitive: true},
			{Name: "id", Label: "ID", Type: "integer", Visible: true, Readable: true, Writable: false, PolicyConfigured: true},
		},
	}
	service := NewFieldPermissionService()
	policies := service.ManifestPolicies(manifest)

	if err := service.ValidateQueryField(manifest, policies, "token", "search"); !errors.Is(err, ErrFieldQueryDenied) {
		t.Fatalf("expected sensitive query denial, got %v", err)
	}
	if err := service.ValidateQueryField(manifest, policies, "name", "search"); err != nil {
		t.Fatalf("expected name search to pass, got %v", err)
	}
	if _, err := service.ValidateWritablePayload(map[string]any{"id": 7}, manifest, policies); !errors.Is(err, ErrFieldPermissionDenied) {
		t.Fatalf("expected read-only write denial, got %v", err)
	}
}

func TestFieldPermissionServiceProjectsRecord(t *testing.T) {
	manifest := resource.Manifest{
		Name: "users", Label: "Users", Route: "/users",
		Fields: []resource.Field{
			{Name: "name", Label: "Name", Type: "text"},
			{Name: "password", Label: "Password", Type: "password", Readable: false, PolicyConfigured: true},
		},
	}
	service := NewFieldPermissionService()
	projected := service.ProjectRecord(map[string]any{"id": 1, "name": "Ada", "password": "secret"}, manifest, service.ManifestPolicies(manifest))
	if _, exists := projected["password"]; exists {
		t.Fatalf("unreadable field leaked: %#v", projected)
	}
	if projected["name"] != "Ada" {
		t.Fatalf("readable field missing: %#v", projected)
	}
}
