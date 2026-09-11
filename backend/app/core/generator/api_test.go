package generator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateAPIWritesDeterministicContractSchemasAndClient(t *testing.T) {
	root := t.TempDir()
	options := APIOptions{RootDir: root}

	if err := GenerateAPI(options); err != nil {
		t.Fatalf("GenerateAPI() error = %v", err)
	}

	for _, name := range []string{
		"contracts/openapi/openapi.json",
		"contracts/schemas/demo-resource.json",
		"contracts/schemas/user-resource.json",
		"admin/src/generated/api/models.ts",
		"admin/src/generated/api/client.ts",
		"admin/src/generated/api/index.ts",
	} {
		if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(name))); err != nil {
			t.Fatalf("generated file %s: %v", name, err)
		}
	}

	client, err := os.ReadFile(filepath.Join(root, "admin/src/generated/api/client.ts"))
	if err != nil {
		t.Fatal(err)
	}
	clientText := string(client)
	for _, marker := range []string{"export interface GeneratedApiClient", "listDemoResources", "bulkDeleteDemoResources", "listModules", "getModule", "listPlugins", "getPlugin", "setPluginState"} {
		if !strings.Contains(clientText, marker) {
			t.Errorf("client is missing %q", marker)
		}
	}

	first := snapshotGeneratedAPI(t, root)
	if err := GenerateAPI(options); err != nil {
		t.Fatalf("second GenerateAPI() error = %v", err)
	}
	second := snapshotGeneratedAPI(t, root)
	if first != second {
		t.Fatal("GenerateAPI() output is not deterministic")
	}
}

func snapshotGeneratedAPI(t *testing.T, root string) string {
	t.Helper()
	var result strings.Builder
	for _, name := range []string{
		"contracts/openapi/openapi.json",
		"contracts/schemas/demo-resource.json",
		"admin/src/generated/api/models.ts",
		"admin/src/generated/api/client.ts",
		"admin/src/generated/api/index.ts",
	} {
		data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(name)))
		if err != nil {
			t.Fatal(err)
		}
		result.Write(data)
	}
	return result.String()
}
