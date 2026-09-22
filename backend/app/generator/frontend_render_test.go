package generator

import (
	"strings"
	"testing"
)

func TestRenderFrontendArtifactsFromResourceSpec(t *testing.T) {
	spec, err := Normalize(Input{Name: "orders", Label: "Orders", Route: "/admin/orders", Icon: "shopping-cart", Fields: []string{"number:text:required", "paid:boolean", "state:select:required:open=Open|closed=Closed"}})
	if err != nil {
		t.Fatal(err)
	}
	artifacts, err := RenderFrontend(spec)
	if err != nil {
		t.Fatal(err)
	}
	if len(artifacts) != 3 {
		t.Fatalf("artifact count = %d, want 3", len(artifacts))
	}
	paths := artifactPaths(artifacts)
	for _, want := range []string{
		"admin/src/modules/orders/resource.ts",
		"admin/src/modules/orders/api.ts",
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
			for _, fragment := range []string{"generatedApi", "list", "show", "create", "update", "remove"} {
				if !strings.Contains(content, fragment) {
					t.Fatalf("api contract missing %q", fragment)
				}
			}
		}
		if artifact.Path == "admin/src/modules/orders/resource.ts" {
			content := string(artifact.Content)
			for _, fragment := range []string{"pageMode: \"generic\"", "name: \"state\"", "value: \"open\"", "label: \"Closed\""} {
				if !strings.Contains(content, fragment) {
					t.Fatalf("generated resource metadata missing %q", fragment)
				}
			}
		}
	}
}

func TestRenderFrontendOwnScopeMetadata(t *testing.T) {
	spec, err := Normalize(Input{Name: "orders", Scope: "own", OwnerField: "owner_id", Fields: []string{"owner_id:integer", "number:text"}})
	if err != nil {
		t.Fatal(err)
	}
	artifacts, err := RenderFrontend(spec)
	if err != nil {
		t.Fatal(err)
	}
	for _, artifact := range artifacts {
		if artifact.Path == "admin/src/modules/orders/resource.ts" {
			content := string(artifact.Content)
			if !strings.Contains(content, `dataScope: "own"`) || !strings.Contains(content, `ownerField: "owner_id"`) {
				t.Fatalf("own scope metadata missing: %s", content)
			}
			return
		}
	}
	t.Fatal("resource metadata artifact not found")
}

func TestRenderFrontendCustomResourceIncludesDedicatedPages(t *testing.T) {
	spec, err := Normalize(Input{Name: "orders", PageMode: "custom", Fields: []string{"number:text"}})
	if err != nil {
		t.Fatal(err)
	}
	artifacts, err := RenderFrontend(spec)
	if err != nil {
		t.Fatal(err)
	}
	paths := artifactPaths(artifacts)
	for _, want := range []string{
		"admin/src/modules/orders/pages/OrdersListPage.vue",
		"admin/src/modules/orders/pages/OrdersFormPage.vue",
		"admin/src/modules/orders/pages/OrdersDetailPage.vue",
	} {
		found := false
		for _, path := range paths {
			if path == want {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("missing custom page artifact %q in %v", want, paths)
		}
	}
}

func TestRenderFrontendFieldPolicies(t *testing.T) {
	spec, err := Normalize(Input{Name: "users", Fields: []string{"email:email", "password:text:sensitive", "id:integer:readonly"}})
	if err != nil {
		t.Fatal(err)
	}
	artifacts, err := RenderFrontend(spec)
	if err != nil {
		t.Fatal(err)
	}
	for _, artifact := range artifacts {
		if artifact.Path == "admin/src/modules/users/resource.ts" {
			content := string(artifact.Content)
			for _, fragment := range []string{"sensitive: true", "writable: false", "readable: true", "visible: true"} {
				if !strings.Contains(content, fragment) {
					t.Fatalf("generated field policy missing %q: %s", fragment, content)
				}
			}
			return
		}
	}
	t.Fatal("resource metadata artifact not found")
}
