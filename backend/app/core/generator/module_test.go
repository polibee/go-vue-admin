package generator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateModuleCreatesStableModuleContract(t *testing.T) {
	root := t.TempDir()

	err := GenerateModule(ModuleOptions{RootDir: root, Name: "catalog"})
	if err != nil {
		t.Fatalf("GenerateModule() error = %v", err)
	}

	files := []string{
		"README.md",
		"acceptance.md",
		"admin/locales/README.md",
		"admin/module.ts",
		"admin/package.json",
		"admin/resources/README.md",
		"admin/routes/README.md",
		"backend/README.md",
		"backend/go.mod",
		"backend/module.go",
		"backend/module_test.go",
		"module.yaml",
	}
	for _, file := range files {
		if _, err := os.Stat(filepath.Join(root, "modules", "catalog", file)); err != nil {
			t.Errorf("generated file %s: %v", file, err)
		}
	}

	manifest, err := os.ReadFile(filepath.Join(root, "modules", "catalog", "module.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	for _, expected := range []string{"id: catalog", "name: Catalog Module", "runtime: compiled", "dependencies: []"} {
		if !contains(string(manifest), expected) {
			t.Errorf("module manifest does not contain %q", expected)
		}
	}
}

func TestGenerateModuleRejectsDuplicateModule(t *testing.T) {
	root := t.TempDir()
	if err := GenerateModule(ModuleOptions{RootDir: root, Name: "catalog"}); err != nil {
		t.Fatal(err)
	}
	if err := GenerateModule(ModuleOptions{RootDir: root, Name: "catalog"}); err == nil {
		t.Fatal("GenerateModule() accepted a duplicate module")
	}
}

func contains(value, expected string) bool {
	return len(value) >= len(expected) && strings.Contains(value, expected)
}
