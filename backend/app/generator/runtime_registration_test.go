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

func TestRenderRuntimeRegistrationUsesExplicitCustomPageOverrides(t *testing.T) {
	root := t.TempDir()
	resourceRoot := filepath.Join(root, "admin", "src", "modules", "orders")
	if err := os.MkdirAll(filepath.Join(resourceRoot, "pages"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(resourceRoot, "resource.ts"), []byte("export const resourceDefinition = { name: \"orders\", pageMode: \"custom\" }"), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"OrdersListPage.vue", "OrdersFormPage.vue", "OrdersDetailPage.vue"} {
		if err := os.WriteFile(filepath.Join(resourceRoot, "pages", name), []byte("custom"), 0o644); err != nil {
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
	frontend := string(artifacts[1].Content)
	for _, fragment := range []string{"OrdersListPage.vue", "OrdersFormPage.vue", "OrdersDetailPage.vue"} {
		if !strings.Contains(frontend, fragment) {
			t.Fatalf("custom frontend registry missing %q: %s", fragment, frontend)
		}
	}
}

func TestRenderRuntimeRegistrationRejectsMissingCustomPageOverride(t *testing.T) {
	root := t.TempDir()
	resourceRoot := filepath.Join(root, "admin", "src", "modules", "orders")
	if err := os.MkdirAll(resourceRoot, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(resourceRoot, "resource.ts"), []byte("export const resourceDefinition = { name: \"orders\", pageMode: \"custom\" }"), 0o644); err != nil {
		t.Fatal(err)
	}
	spec, err := Normalize(Input{Name: "posts", Fields: []string{"title:text"}})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := RenderRuntimeRegistration(root, spec); err == nil {
		t.Fatal("expected missing custom page override to fail")
	}
}

func TestRenderRuntimeRegistrationAcceptsCustomPageModeWithSingleQuotesAndSpacing(t *testing.T) {
	root := t.TempDir()
	resourceRoot := filepath.Join(root, "admin", "src", "modules", "orders")
	if err := os.MkdirAll(filepath.Join(resourceRoot, "pages"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(resourceRoot, "resource.ts"), []byte("export const resourceDefinition = { name: 'orders', pageMode : 'custom' }"), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"OrdersListPage.vue", "OrdersFormPage.vue", "OrdersDetailPage.vue"} {
		if err := os.WriteFile(filepath.Join(resourceRoot, "pages", name), []byte("custom"), 0o644); err != nil {
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
	frontend := string(artifacts[1].Content)
	if !strings.Contains(frontend, "OrdersListPage.vue") {
		t.Fatalf("custom frontend registry did not recognize single-quoted page mode: %s", frontend)
	}
}
