package generator

import "testing"

func TestParseFieldAcceptsRequiredSuffix(t *testing.T) {
	field, err := ParseField("title:text:required")
	if err != nil {
		t.Fatalf("parse field: %v", err)
	}
	if field.Name != "title" || field.Type != "text" || !field.Required {
		t.Fatalf("field = %+v, want title/text/required", field)
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
