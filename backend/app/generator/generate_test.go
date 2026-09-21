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
