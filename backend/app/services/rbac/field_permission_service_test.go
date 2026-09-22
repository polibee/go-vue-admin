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

func TestValidateFieldOverridesRejectsUnknownAndManifestExpansion(t *testing.T) {
	manifest := resource.Manifest{Fields: []resource.Field{
		{Name: "name", Label: "Name", Type: "text"},
		{Name: "secret", Label: "Secret", Type: "text", Visible: true, Readable: false, Writable: false, PolicyConfigured: true},
	}}
	service := NewFieldPermissionService()
	if err := service.ValidateFieldOverrides(manifest, map[string]FieldOverride{"missing": {Readable: false}}); !errors.Is(err, ErrFieldNotFound) {
		t.Fatalf("expected unknown field rejection, got %v", err)
	}
	if err := service.ValidateFieldOverrides(manifest, map[string]FieldOverride{"secret": {Readable: true, Writable: true}}); !errors.Is(err, ErrFieldPolicyExpansion) {
		t.Fatalf("expected manifest expansion rejection, got %v", err)
	}
	if err := service.ValidateFieldOverrides(manifest, map[string]FieldOverride{"name": {Readable: false, Writable: false}}); err != nil {
		t.Fatalf("expected restrictive override to pass, got %v", err)
	}
}

func TestMergeFieldPoliciesPreservesManifestBoundsAcrossRoles(t *testing.T) {
	base := map[string]FieldPolicy{
		"email":  {Visible: true, Readable: true, Writable: true},
		"secret": {Visible: false, Readable: false, Writable: false, Sensitive: true},
	}
	merged := mergeFieldPolicies(base, []map[string]FieldPolicy{
		{"email": {Visible: true, Readable: false, Writable: false}, "secret": {Visible: false, Readable: true, Writable: true, Sensitive: true}},
		{"email": {Visible: true, Readable: true, Writable: true}},
	})
	if !merged["email"].Readable || !merged["email"].Writable {
		t.Fatalf("expected any-role allow to merge: %+v", merged["email"])
	}
	if merged["secret"].Visible || merged["secret"].Readable || merged["secret"].Writable {
		t.Fatalf("manifest bounds were expanded: %+v", merged["secret"])
	}
}
