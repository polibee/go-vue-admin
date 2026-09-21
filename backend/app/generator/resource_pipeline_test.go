package generator

import "testing"

func TestRenderResourcePipelineComposesAllResourceArtifacts(t *testing.T) {
	spec, err := Normalize(Input{Name: "orders", Fields: []string{"number:text:required"}})
	if err != nil {
		t.Fatal(err)
	}
	artifacts, err := RenderResourcePipeline(spec, "20260921000000")
	if err != nil {
		t.Fatalf("render resource pipeline: %v", err)
	}
	if len(artifacts) != 23 {
		t.Fatalf("artifact count = %d, want 23", len(artifacts))
	}
	if artifacts[0].Path != "app/generated/resources/orders/manifest.go" || artifacts[11].Path != "app/generated/permissions/orders/permissions.go" || artifacts[13].Path != "app/generated/menus/orders/menu.go" || artifacts[15].Path != "admin/src/generated/resources/orders/resource.ts" {
		t.Fatalf("unexpected pipeline boundaries: %v, %v, %v, %v", artifacts[0].Path, artifacts[11].Path, artifacts[13].Path, artifacts[15].Path)
	}
}
