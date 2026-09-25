package generator

import "testing"

func TestNormalizeNamespaceDefaultsToAdmin(t *testing.T) {
	namespace, err := normalizeNamespace("")
	if err != nil || namespace != "admin" {
		t.Fatalf("normalizeNamespace(\"\") = %q, %v; want admin", namespace, err)
	}
}

func TestNormalizeNamespaceAcceptsApp(t *testing.T) {
	namespace, err := normalizeNamespace("app")
	if err != nil || namespace != "app" {
		t.Fatalf("normalizeNamespace(\"app\") = %q, %v; want app", namespace, err)
	}
}

func TestNormalizeNamespaceRejectsInvalidValue(t *testing.T) {
	for _, value := range []string{"Admin", "public", "app.admin", "../app"} {
		if _, err := normalizeNamespace(value); err == nil {
			t.Fatalf("normalizeNamespace(%q) accepted invalid namespace", value)
		}
	}
}

func TestPermissionNameUsesNamespace(t *testing.T) {
	if got := permissionName("app", "orders", "view"); got != "app.orders.view" {
		t.Fatalf("permissionName() = %q; want app.orders.view", got)
	}
}
