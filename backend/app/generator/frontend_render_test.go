package generator

import "testing"

func TestRenderFrontendArtifactsFromResourceSpec(t *testing.T) {
	spec, err := Normalize(Input{Name: "orders", Label: "Orders", Route: "/admin/orders", Icon: "shopping-cart", Fields: []string{"number:text:required", "paid:boolean"}})
	if err != nil {
		t.Fatal(err)
	}
	artifacts, err := RenderFrontend(spec)
	if err != nil {
		t.Fatal(err)
	}
	if len(artifacts) != 7 {
		t.Fatalf("artifact count = %d, want 7", len(artifacts))
	}
	paths := artifactPaths(artifacts)
	for _, want := range []string{
		"admin/src/generated/resources/orders/resource.ts",
		"admin/src/generated/resources/orders/menu.ts",
		"admin/src/generated/resources/orders/routes.ts",
		"admin/src/generated/resources/orders/pages/OrdersListPage.vue",
		"admin/src/generated/resources/orders/pages/OrdersFormPage.vue",
		"admin/src/generated/resources/orders/pages/OrdersDetailPage.vue",
		"admin/src/generated/resources/orders/orders.test.ts",
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
}
