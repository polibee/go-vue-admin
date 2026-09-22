package registry

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
	if len(manifest.Actions) != 7 {
		t.Fatalf("expected standard CRUD and user status actions, got %+v", manifest.Actions)
	}
	for _, action := range manifest.Actions[:4] {
		if action.Permission != "admin.users.manage" || action.Batch {
			t.Fatalf("action %q has permission %q, want admin.users.manage", action.Name, action.Permission)
		}
	}
	if manifest.Actions[6].Name != "set-status" || manifest.Actions[6].Kind != "user-status" || !manifest.Actions[6].Batch || manifest.Actions[6].Payload != "user-status" {
		t.Fatalf("expected declarative user status action, got %+v", manifest.Actions[4])
	}
}
