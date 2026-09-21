package controllers

import (
	"testing"

	"goravel/app/core/resource"
)

func TestResourceListCapabilitiesDeriveFromManifestFields(t *testing.T) {
	manifest := resource.Manifest{Fields: []resource.Field{
		{Name: "title", Type: "text"},
		{Name: "email", Type: "email"},
		{Name: "status", Type: "select"},
		{Name: "published", Type: "boolean"},
		{Name: "secret", Type: "password"},
	}}

	search := resourceSearchFields(manifest)
	if got, want := joinFieldNames(search), "title,email"; got != want {
		t.Fatalf("search fields = %q, want %q", got, want)
	}

	filters := resourceFilterFields(manifest)
	if got, want := joinFieldNames(filters), "status,published"; got != want {
		t.Fatalf("filter fields = %q, want %q", got, want)
	}
}

func joinFieldNames(fields []resource.Field) string {
	result := ""
	for index, field := range fields {
		if index > 0 {
			result += ","
		}
		result += field.Name
	}
	return result
}
