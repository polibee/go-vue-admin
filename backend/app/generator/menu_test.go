package generator

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNormalizeMenuRejectsInvalidRoute(t *testing.T) {
	if _, err := NormalizeMenu(MenuInput{Name: "orders", Label: "Orders", Route: "orders", Permission: "admin.orders.view"}); err == nil {
		t.Fatal("expected route without leading slash to fail")
	}
}

func TestRenderMenuProducesStableArtifacts(t *testing.T) {
	spec, err := NormalizeMenu(MenuInput{Name: "orders", Label: "Orders", Route: "/orders", Permission: "admin.orders.view", Icon: "shopping-cart"})
	if err != nil {
		t.Fatal(err)
	}
	artifacts, err := RenderMenu(spec)
	if err != nil {
		t.Fatal(err)
	}
	if len(artifacts) != 2 || artifacts[0].Path != "app/modules/orders/menu/menu.go" || artifacts[1].Path != "app/modules/orders/menu/README.md" {
		t.Fatalf("artifacts = %+v", artifactPaths(artifacts))
	}
	for _, artifact := range artifacts {
		golden := filepath.Join("testdata", "menus", filepath.Base(artifact.Path)+".golden")
		if os.Getenv("UPDATE_GOLDEN") == "1" {
			if err := os.MkdirAll(filepath.Dir(golden), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(golden, artifact.Content, 0o644); err != nil {
				t.Fatal(err)
			}
			continue
		}
		want, err := os.ReadFile(golden)
		if err != nil {
			t.Fatalf("read golden: %v", err)
		}
		if string(artifact.Content) != string(want) {
			t.Errorf("content mismatch for %s", artifact.Path)
		}
	}
}
