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
