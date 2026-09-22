package auditservices

import (
	"strings"
	"testing"
)

func TestShouldAuditHTTP(t *testing.T) {
	for _, path := range []string{"/api/v1/auth/login", "/api/v1/admin/users"} {
		if !ShouldAuditHTTP("POST", path) {
			t.Fatalf("expected API path to be audited: %s", path)
		}
	}
	for _, path := range []string{"/api/v1/admin/audit-logs", "/api/docs", "/"} {
		if ShouldAuditHTTP("GET", path) {
			t.Fatalf("did not expect path to be audited: %s", path)
		}
	}
	if ShouldAuditHTTP("OPTIONS", "/api/v1/admin/users") {
		t.Fatal("OPTIONS must not be audited")
	}
}

func TestBuildHTTPAuditMetadataContainsRequestAndResponseSummary(t *testing.T) {
	metadata := BuildHTTPAuditMetadata(HTTPAuditInput{
		Method:       "POST",
		Path:         "/api/v1/admin/users",
		Query:        map[string]string{"page": "1"},
		RequestBody:  map[string]any{"email": "admin@example.com", "password": "secret"},
		Status:       201,
		ContentType:  "application/json",
		ResponseBody: []byte(`{"data":{"id":1,"access_token":"secret-token"}}`),
	})

	request, ok := metadata["request"].(map[string]any)
	if !ok || request["method"] != "POST" || request["path"] != "/api/v1/admin/users" {
		t.Fatalf("unexpected request summary: %#v", metadata["request"])
	}
	response, ok := metadata["response"].(map[string]any)
	if !ok || response["status"] != 201 {
		t.Fatalf("unexpected response summary: %#v", metadata["response"])
	}
	encoded, _, err := MarshalBounded(metadata)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(encoded), "secret-token") || strings.Contains(string(encoded), "secret") {
		t.Fatalf("HTTP audit metadata leaked sensitive data: %s", encoded)
	}
}
