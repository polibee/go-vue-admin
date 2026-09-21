package generator

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestNormalizeModuleRejectsUnsafeName(t *testing.T) {
	if _, err := NormalizeModule("../billing"); err == nil {
		t.Fatal("expected unsafe module name to fail")
	}
}

func TestRenderModuleGolden(t *testing.T) {
	spec, err := NormalizeModule("billing")
	if err != nil {
		t.Fatalf("normalize module: %v", err)
	}
	artifacts, err := RenderModule(spec)
	if err != nil {
		t.Fatalf("render module: %v", err)
	}
	wantPaths := []string{
		"app/modules/billing/module.go",
		"app/modules/billing/model.go",
		"app/modules/billing/request.go",
		"app/modules/billing/repository.go",
		"app/modules/billing/service.go",
		"app/modules/billing/controller.go",
		"app/modules/billing/routes.go",
		"app/modules/billing/resource.go",
		"app/modules/billing/permissions.go",
		"app/modules/billing/events.go",
		"app/modules/billing/tests/module_test.go",
		"app/modules/billing/README.md",
	}
	if got := artifactPaths(artifacts); !reflect.DeepEqual(got, wantPaths) {
		t.Fatalf("paths = %v, want %v", got, wantPaths)
	}
	for _, artifact := range artifacts {
		golden := filepath.Join("testdata", "billing", filepath.FromSlash(artifact.Path))
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
			t.Fatalf("read golden %s: %v", golden, err)
		}
		if string(artifact.Content) != string(want) {
			t.Errorf("content mismatch for %s", artifact.Path)
		}
	}
}
