package generator

import (
	"reflect"
	"testing"
)

func TestNormalizeCRUDReusesResourceAndModuleIdentity(t *testing.T) {
	spec, err := NormalizeCRUD(Input{Name: "orders", Fields: []string{"number:text:required"}})
	if err != nil {
		t.Fatalf("normalize CRUD: %v", err)
	}
	if spec.Resource.Name != "orders" || spec.Module.Name != "orders" || spec.Resource.GoName != "Orders" || spec.Module.GoName != "Orders" {
		t.Fatalf("spec = %+v", spec)
	}
}

func TestRenderCRUDCombinesResourceAndModuleArtifacts(t *testing.T) {
	spec, err := NormalizeCRUD(Input{Name: "orders", Fields: []string{"number:text:required"}})
	if err != nil {
		t.Fatal(err)
	}
	artifacts, err := RenderCRUD(spec, "20260921000000")
	if err != nil {
		t.Fatalf("render CRUD: %v", err)
	}
	if len(artifacts) != 23 {
		t.Fatalf("artifact count = %d, want 23", len(artifacts))
	}
	if !reflect.DeepEqual(artifactPaths(artifacts[:12]), artifactPathsMust(RenderModule(spec.Module))) {
		t.Fatal("module artifacts were not composed first")
	}
	if artifacts[12].Path != "app/generated/resources/orders/manifest.go" {
		t.Fatalf("resource artifacts start at %q", artifacts[12].Path)
	}
}

func artifactPathsMust(artifacts []Artifact, err error) []string {
	if err != nil {
		panic(err)
	}
	return artifactPaths(artifacts)
}
