package resource

import (
	"errors"
	"testing"
)

func TestRegistryRejectsInvalidAndDuplicateManifests(t *testing.T) {
	registry := NewRegistry()
	if err := registry.Register(Manifest{}); !errors.Is(err, ErrInvalidManifest) {
		t.Fatalf("expected invalid manifest, got %v", err)
	}
	manifest := Manifest{Name: "users", Label: "Users", Route: "/users"}
	if err := registry.Register(manifest); err != nil {
		t.Fatalf("register: %v", err)
	}
	if err := registry.Register(manifest); !errors.Is(err, ErrDuplicate) {
		t.Fatalf("expected duplicate, got %v", err)
	}
}

func TestRegistryReturnsStableSortedManifests(t *testing.T) {
	registry := NewRegistry()
	for _, name := range []string{"roles", "users", "permissions"} {
		if err := registry.Register(Manifest{Name: name, Label: name, Route: "/" + name}); err != nil {
			t.Fatalf("register %s: %v", name, err)
		}
	}
	all := registry.All()
	if got := all[0].Name + "," + all[1].Name + "," + all[2].Name; got != "permissions,roles,users" {
		t.Fatalf("unexpected order: %s", got)
	}
	if _, err := registry.Find("missing"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected not found, got %v", err)
	}
}

func TestRegistryValidatesDataScopeMetadata(t *testing.T) {
	own := Manifest{
		Name: "orders", Label: "Orders", Route: "/orders",
		DataScope: DataScopeOwn, OwnerField: "owner_id",
		Fields: []Field{{Name: "owner_id", Label: "Owner", Type: "integer"}},
	}
	if err := NewRegistry().Register(own); err != nil {
		t.Fatalf("expected own scope manifest to register: %v", err)
	}

	unknown := own
	unknown.Name = "unknown-scope"
	unknown.DataScope = DataScope("department")
	if err := NewRegistry().Register(unknown); !errors.Is(err, ErrInvalidManifest) {
		t.Fatalf("expected unknown scope to be rejected, got %v", err)
	}

	missingOwner := own
	missingOwner.Name = "missing-owner"
	missingOwner.OwnerField = "missing_id"
	if err := NewRegistry().Register(missingOwner); !errors.Is(err, ErrInvalidManifest) {
		t.Fatalf("expected missing owner field to be rejected, got %v", err)
	}
}

func TestRegistryValidatesRelationsAndFormExtensions(t *testing.T) {
	registry := NewRegistry()
	manifest := Manifest{
		Name: "orders", Label: "Orders", Route: "/admin/orders",
		Fields:       []Field{{Name: "customer_id"}, {Name: "status"}},
		Relations:    []Relation{{Name: "customer", Kind: "belongsTo", Resource: "customers", Field: "customer_id", ForeignField: "id", LabelField: "name", Selectable: true}},
		FormGroups:   []FormGroup{{Name: "main", Label: "Main", Columns: 2, Fields: []string{"customer_id", "status"}}},
		Details:      []DetailSection{{Name: "summary", Label: "Summary", Fields: []string{"status"}}},
		Dependencies: []FieldDependency{{Field: "customer_id", On: "status", Value: "active"}},
	}
	if err := registry.Register(manifest); err != nil {
		t.Fatalf("register valid extensions: %v", err)
	}
	for _, invalid := range []Manifest{
		{Name: "invalid-kind", Label: "Invalid", Route: "/admin/invalid-kind", Fields: []Field{{Name: "id"}}, Relations: []Relation{{Name: "x", Kind: "hasOne", Resource: "other", Field: "id", ForeignField: "id", LabelField: "name"}}},
		{Name: "invalid-field", Label: "Invalid", Route: "/admin/invalid-field", Fields: []Field{{Name: "id"}}, Relations: []Relation{{Name: "x", Kind: "belongsTo", Resource: "other", Field: "missing", ForeignField: "id", LabelField: "name"}}},
		{Name: "invalid-group", Label: "Invalid", Route: "/admin/invalid-group", Fields: []Field{{Name: "id"}}, FormGroups: []FormGroup{{Name: "main", Label: "Main", Columns: 5, Fields: []string{"id"}}}},
	} {
		if err := NewRegistry().Register(invalid); err != ErrInvalidManifest {
			t.Fatalf("expected invalid extension manifest, got %v", err)
		}
	}
}

func TestFieldPolicyDefaultsPreserveExistingManifestBehavior(t *testing.T) {
	manifest := Manifest{
		Name: "posts", Label: "Posts", Route: "/posts",
		Fields: []Field{{Name: "title", Label: "Title", Type: "text"}},
	}

	normalized := NormalizeManifestFields(manifest)
	field := normalized.Fields[0]
	if !field.Visible || !field.Readable || !field.Writable || field.Sensitive {
		t.Fatalf("unexpected default field policy: %+v", field)
	}
}

func TestFieldPolicyPreservesExplicitRestrictions(t *testing.T) {
	manifest := Manifest{
		Name: "users", Label: "Users", Route: "/users",
		Fields: []Field{{Name: "password", Label: "Password", Type: "password", Visible: false, Readable: false, Writable: false, Sensitive: true}},
	}

	normalized := NormalizeManifestFields(manifest)
	field := normalized.Fields[0]
	if field.Visible || field.Readable || field.Writable || !field.Sensitive {
		t.Fatalf("explicit field policy was not preserved: %+v", field)
	}
}
