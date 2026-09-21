package generator

import (
	"strings"
	"testing"
)

func TestRenderFrontendArtifactsFromResourceSpec(t *testing.T) {
	spec, err := Normalize(Input{Name: "orders", Label: "Orders", Route: "/admin/orders", Icon: "shopping-cart", Fields: []string{"number:text:required", "paid:boolean"}})
	if err != nil {
		t.Fatal(err)
	}
	artifacts, err := RenderFrontend(spec)
	if err != nil {
		t.Fatal(err)
	}
	if len(artifacts) != 8 {
		t.Fatalf("artifact count = %d, want 8", len(artifacts))
	}
	paths := artifactPaths(artifacts)
	for _, want := range []string{
		"admin/src/modules/orders/resource.ts",
		"admin/src/modules/orders/api.ts",
		"admin/src/modules/orders/menu.ts",
		"admin/src/modules/orders/routes.ts",
		"admin/src/modules/orders/pages/OrdersListPage.vue",
		"admin/src/modules/orders/pages/OrdersFormPage.vue",
		"admin/src/modules/orders/pages/OrdersDetailPage.vue",
		"admin/src/modules/orders/orders.test.ts",
	} {
		found := false
		for _, path := range paths {
			if path == want {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("missing frontend artifact %q in %v", want, paths)
		}
	}
	for _, artifact := range artifacts {
		if len(artifact.Content) == 0 {
			t.Fatalf("empty frontend artifact %s", artifact.Path)
		}
	}
	for _, artifact := range artifacts {
		if artifact.Path == "admin/src/modules/orders/api.ts" {
			content := string(artifact.Content)
			for _, fragment := range []string{"list", "show", "create", "update", "remove", "/api/v1/admin/orders"} {
				if !strings.Contains(content, fragment) {
					t.Fatalf("api contract missing %q", fragment)
				}
			}
		}
		if artifact.Path == "admin/src/modules/orders/pages/OrdersFormPage.vue" && !strings.Contains(string(artifact.Content), "ResourceFormView") {
			t.Fatal("generated form page does not use ResourceFormView")
		}
	}
}
