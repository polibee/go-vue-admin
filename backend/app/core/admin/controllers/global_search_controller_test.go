package controllers

import (
	"testing"

	"goravel/app/core/resource"
)

func TestNormalizeGlobalSearchQuery(t *testing.T) {
	if query, ok := normalizeGlobalSearchQuery("  users  "); !ok || query != "users" {
		t.Fatalf("unexpected normalized query: %q, %v", query, ok)
	}
	if query, ok := normalizeGlobalSearchQuery("   "); ok || query != "" {
		t.Fatalf("blank query should be ignored: %q, %v", query, ok)
	}
	long := "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz"
	if query, ok := normalizeGlobalSearchQuery(long); !ok || len(query) != globalSearchMaxQuery {
		t.Fatalf("query should be bounded to %d characters, got %d", globalSearchMaxQuery, len(query))
	}
}

func TestGlobalSearchFieldsExcludeSensitiveValues(t *testing.T) {
	fields := globalSearchFields(resource.Manifest{Fields: []resource.Field{
		{Name: "name", Type: "text"},
		{Name: "email", Type: "email"},
		{Name: "password", Type: "password"},
		{Name: "api_token", Type: "text"},
		{Name: "notes", Type: "textarea"},
	}})
	if len(fields) != 2 || fields[0].Name != "name" || fields[1].Name != "email" {
		t.Fatalf("unexpected searchable fields: %#v", fields)
	}
}

func TestFrontendSearchRoute(t *testing.T) {
	if got := frontendSearchRoute("/admin/orders", int64(42)); got != "/admin/orders/42" {
		t.Fatalf("unexpected frontend route: %s", got)
	}
}
