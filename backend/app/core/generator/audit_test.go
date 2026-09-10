package generator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAuditAPIFindsUndocumentedRoute(t *testing.T) {
	root := auditFixture(t, `package routes
import "goravel/app/facades"
func Web() { facades.Route().Get("/api/ghost", func() {}) }
`)
	findings, err := AuditAPI(APIAuditOptions{RootDir: root})
	if err != nil {
		t.Fatal(err)
	}
	if !hasAuditFinding(findings, "/api/ghost", "openapi") {
		t.Fatalf("expected undocumented route finding, got %#v", findings)
	}
}

func TestAuditAPIFindsMissingGeneratedOperation(t *testing.T) {
	root := auditFixture(t, `package routes
import "goravel/app/facades"
func Web() { facades.Route().Get("/api/resources/demo", func() {}) }
`)
	clientPath := filepath.Join(root, "admin", "src", "generated", "api", "client.ts")
	client, err := os.ReadFile(clientPath)
	if err != nil {
		t.Fatal(err)
	}
	client = []byte(strings.ReplaceAll(string(client), "listDemoResources", "listDemoResources_removed"))
	if err := os.WriteFile(clientPath, client, 0o644); err != nil {
		t.Fatal(err)
	}
	findings, err := AuditAPI(APIAuditOptions{RootDir: root})
	if err != nil {
		t.Fatal(err)
	}
	if !hasAuditFinding(findings, "/api/resources/demo", "generated-client") {
		t.Fatalf("expected missing client finding, got %#v", findings)
	}
}

func auditFixture(t *testing.T, routes string) string {
	t.Helper()
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "backend", "routes"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "backend", "routes", "web.go"), []byte(routes), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := GenerateAPI(APIOptions{RootDir: root}); err != nil {
		t.Fatal(err)
	}
	return root
}

func hasAuditFinding(findings []APIAuditFinding, path, layer string) bool {
	for _, finding := range findings {
		if finding.Path == path && finding.Layer == layer {
			return true
		}
	}
	return false
}
