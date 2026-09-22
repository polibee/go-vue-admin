package generator

import (
	"reflect"
	"testing"
)

func TestParseFieldAcceptsRequiredSuffix(t *testing.T) {
	field, err := ParseField("title:text:required")
	if err != nil {
		t.Fatalf("parse field: %v", err)
	}
	if field.Name != "title" || field.Type != "text" || !field.Required {
		t.Fatalf("field = %+v, want title/text/required", field)
	}
}

func TestParseFieldAcceptsSelectOptions(t *testing.T) {
	field, err := ParseField("status:select:required:active=Active|disabled=Disabled")
	if err != nil {
		t.Fatalf("parse select field: %v", err)
	}
	if len(field.Options) != 2 || field.Options[0].Value != "active" || field.Options[0].Label != "Active" || !field.Required {
		t.Fatalf("field = %+v, want required select options", field)
	}
}

func TestNormalizeRejectsUnsafeResourceName(t *testing.T) {
	_, err := Normalize(Input{Name: "../posts", Fields: []string{"title:text"}})
	if err == nil {
		t.Fatal("expected unsafe resource name to fail")
	}
}

func TestNormalizeRejectsDuplicateFields(t *testing.T) {
	_, err := Normalize(Input{Name: "posts", Fields: []string{"title:text", "title:email"}})
	if err == nil {
		t.Fatal("expected duplicate fields to fail")
	}
}

func TestNormalizeAppliesResourceDefaults(t *testing.T) {
	spec, err := Normalize(Input{Name: "blog-posts", Fields: []string{"title:text"}})
	if err != nil {
		t.Fatalf("normalize: %v", err)
	}
	if spec.GoName != "BlogPosts" || spec.Label != "Blog Posts" || spec.Route != "/admin/blog-posts" || spec.Permission != "admin.blog-posts.view" {
		t.Fatalf("spec defaults = %+v", spec)
	}
}

func TestNormalizeBuildsSharedMenuAndPageMetadata(t *testing.T) {
	spec, err := Normalize(Input{Name: "orders", Route: "/admin/orders", Icon: "shopping-cart", Fields: []string{"number:text"}})
	if err != nil {
		t.Fatalf("normalize: %v", err)
	}
	if spec.Icon != "shopping-cart" || spec.FrontendRoute != "/orders" {
		t.Fatalf("shared navigation metadata = %+v", spec)
	}
	wantActions := []string{"view", "create", "update", "delete"}
	if !reflect.DeepEqual(spec.Actions, wantActions) {
		t.Fatalf("actions = %v, want %v", spec.Actions, wantActions)
	}
}

func TestNormalizeSupportsOwnDataScopeWithDeclaredOwnerField(t *testing.T) {
	spec, err := Normalize(Input{
		Name: "orders", Scope: "own", OwnerField: "owner_id",
		Fields: []string{"owner_id:integer", "number:text"},
	})
	if err != nil {
		t.Fatalf("normalize own scope: %v", err)
	}
	if spec.DataScope != "own" || spec.OwnerField != "owner_id" {
		t.Fatalf("scope metadata = %+v", spec)
	}

	if _, err := Normalize(Input{Name: "invalid-orders", Scope: "own", Fields: []string{"number:text"}}); err == nil {
		t.Fatal("expected own scope without owner field to fail")
	}
	if _, err := Normalize(Input{Name: "invalid-type", Scope: "own", OwnerField: "owner_id", Fields: []string{"owner_id:text"}}); err == nil {
		t.Fatal("expected non-integer owner field to fail")
	}
}
