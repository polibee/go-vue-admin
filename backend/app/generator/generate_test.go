package generator

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGenerateResourceWritesPipelineAndDiscovery(t *testing.T) {
	root := t.TempDir()
	artifacts, err := GenerateResource(root, Input{
		Name:   "billing",
		Fields: []string{"title:text:required", "status:select:required:active=Active|disabled=Disabled"},
	}, "20260921000000")
	if err != nil {
		t.Fatalf("generate resource: %v", err)
	}
	if len(artifacts) == 0 {
		t.Fatal("expected generated artifacts")
	}
	for _, path := range []string{
		"app/modules/billing/resource/manifest.go",
		"admin/src/modules/billing/resource.ts",
		"app/modules/admin/registry/generated_resources.go",
		"admin/src/core/resource/generated.ts",
	} {
		if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(path))); err != nil {
			t.Fatalf("missing generated file %s: %v", path, err)
		}
	}
}

func TestGenerateResourceUsesRepositoryAdminDirectoryWhenRunFromBackend(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "backend"), 0o755); err != nil {
		t.Fatalf("create backend root: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(root, "admin", "src", "modules"), 0o755); err != nil {
		t.Fatalf("create admin root: %v", err)
	}
	backendRoot := filepath.Join(root, "backend")
	if _, err := GenerateResource(backendRoot, Input{Name: "billing", Fields: []string{"title:text:required"}}, "20260921000000"); err != nil {
		t.Fatalf("generate resource from backend root: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "admin", "src", "modules", "billing", "resource.ts")); err != nil {
		t.Fatalf("frontend resource was not written to repository admin directory: %v", err)
	}
	if _, err := os.Stat(filepath.Join(backendRoot, "admin")); !os.IsNotExist(err) {
		t.Fatalf("backend admin directory should not be created, err=%v", err)
	}
}
