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
