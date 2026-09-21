package resources

import "testing"

func TestAdminRegistryExposesFormFields(t *testing.T) {
	manifest, err := AdminRegistry().Find("users")
	if err != nil {
		t.Fatalf("find users: %v", err)
	}
	if len(manifest.Fields) != 5 {
		t.Fatalf("expected five user fields, got %d", len(manifest.Fields))
	}
	for index, want := range []struct{ name, fieldType string }{
		{"name", "text"},
		{"email", "email"},
		{"password", "password"},
		{"locale", "text"},
		{"status", "select"},
	} {
		field := manifest.Fields[index]
		if field.Name != want.name || field.Type != want.fieldType {
			t.Fatalf("field %d = %+v, want %s/%s", index, field, want.name, want.fieldType)
		}
	}
}
