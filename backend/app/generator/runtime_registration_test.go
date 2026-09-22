package generator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRenderRuntimeRegistrationDiscoversExistingAndNewResource(t *testing.T) {
	root := t.TempDir()
	for _, path := range []string{
		"app/modules/orders/resource/manifest.go",
		"admin/src/modules/orders/resource.ts",
	} {
		fullPath := filepath.Join(root, filepath.FromSlash(path))
		if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(fullPath, []byte("generated"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	spec, err := Normalize(Input{Name: "posts", Fields: []string{"title:text"}})
	if err != nil {
		t.Fatal(err)
	}
	artifacts, err := RenderRuntimeRegistration(root, spec)
	if err != nil {
		t.Fatal(err)
	}
	if len(artifacts) != 2 || !artifacts[0].AllowOverwrite || !artifacts[1].AllowOverwrite {
		t.Fatalf("unexpected runtime artifacts: %+v", artifacts)
	}
	backend := string(artifacts[0].Content)
	for _, fragment := range []string{"ordersresource", "postsresource", "registerGenerated"} {
		if !strings.Contains(backend, fragment) {
			t.Fatalf("backend registry missing %q: %s", fragment, backend)
		}
	}
	frontend := string(artifacts[1].Content)
	for _, fragment := range []string{"ResourceListPage", "ResourceFormPage", "ResourceDetailPage", "ordersResource", "postsResource", "generatedResourceRoutes"} {
		if !strings.Contains(frontend, fragment) {
			t.Fatalf("frontend registry missing %q: %s", fragment, frontend)
		}
	}
	for _, forbidden := range []string{"OrdersListPage.vue", "PostsListPage.vue", "OrdersFormPage.vue", "PostsFormPage.vue"} {
		if strings.Contains(frontend, forbidden) {
			t.Fatalf("generic frontend registry imported dedicated page %q: %s", forbidden, frontend)
		}
	}
}
