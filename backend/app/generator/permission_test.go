package generator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNormalizePermissionRejectsUnknownAction(t *testing.T) {
	if _, err := NormalizePermission(PermissionInput{Name: "orders", Actions: []string{"publish"}}); err == nil {
		t.Fatal("expected unknown permission action to fail")
	}
}

func TestRenderPermissionGolden(t *testing.T) {
	spec, err := NormalizePermission(PermissionInput{Name: "orders", Actions: []string{"view", "create", "update", "delete"}})
	if err != nil {
		t.Fatal(err)
	}
	artifacts, err := RenderPermission(spec)
	if err != nil {
		t.Fatal(err)
	}
	if len(artifacts) != 2 || artifacts[0].Path != "app/modules/orders/permissions/permissions.go" || artifacts[1].Path != "app/modules/orders/permissions/README.md" {
		t.Fatalf("artifacts = %+v", artifactPaths(artifacts))
	}
	for _, artifact := range artifacts {
		golden := filepath.Join("testdata", "permissions", filepath.Base(artifact.Path)+".golden")
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

func TestRenderPermissionUsesNamespace(t *testing.T) {
	spec, err := NormalizePermission(PermissionInput{Name: "orders", Namespace: "app", Actions: []string{"view"}})
	if err != nil {
		t.Fatal(err)
	}
	artifacts, err := RenderPermission(spec)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(artifacts[0].Content), `"app.orders.view"`) {
		t.Fatalf("generated permission does not use app namespace: %s", artifacts[0].Content)
	}
}
