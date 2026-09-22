package rbacservices

import (
	"testing"

	"goravel/app/core/resource"
)

func TestOwnerPredicateRequiresOwnScopeMetadata(t *testing.T) {
	manifest := resource.Manifest{
		Name: "orders", DataScope: resource.DataScopeOwn, OwnerField: "owner_id",
		Fields: []resource.Field{{Name: "owner_id", Type: "integer"}},
	}
	field, value, err := OwnerPredicate(manifest, resource.DataScopeOwn, 42)
	if err != nil {
		t.Fatalf("owner predicate error: %v", err)
	}
	if field != "owner_id" || value != int64(42) {
		t.Fatalf("unexpected owner predicate: %q %#v", field, value)
	}

	manifest.OwnerField = ""
	if _, _, err := OwnerPredicate(manifest, resource.DataScopeOwn, 42); err == nil {
		t.Fatal("expected own scope without owner metadata to fail")
	}
}
