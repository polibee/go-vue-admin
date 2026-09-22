package auditservices

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestRedactValueRedactsSensitiveKeysRecursively(t *testing.T) {
	input := map[string]any{
		"email":    "admin@example.com",
		"password": "secret",
		"nested": map[string]any{
			"Authorization": "Bearer token",
			"items":         []any{map[string]any{"refresh_token": "refresh-secret"}},
		},
	}

	got := RedactValue(input)
	encoded, err := json.Marshal(got)
	if err != nil {
		t.Fatal(err)
	}

	text := string(encoded)
	if strings.Contains(text, "secret") || strings.Contains(text, "Bearer token") || strings.Contains(text, "refresh-secret") {
		t.Fatalf("sensitive values leaked: %s", text)
	}
	if !strings.Contains(text, "[REDACTED]") || !strings.Contains(text, "admin@example.com") {
		t.Fatalf("redaction removed expected safe data: %s", text)
	}
}

func TestMarshalBoundedMarksOversizedPayload(t *testing.T) {
	encoded, truncated, err := MarshalBounded(map[string]any{"value": strings.Repeat("x", auditPayloadLimitBytes)})
	if err != nil {
		t.Fatal(err)
	}
	if !truncated {
		t.Fatal("expected oversized payload to be marked truncated")
	}
	if len(encoded) >= auditPayloadLimitBytes {
		t.Fatalf("truncated payload is not bounded: %d", len(encoded))
	}
	if !strings.Contains(string(encoded), "truncated") {
		t.Fatalf("truncated marker missing: %s", encoded)
	}
}

func TestBoundedBytesCopiesOnlyTheConfiguredLimit(t *testing.T) {
	body, truncated := BoundedBytes([]byte(strings.Repeat("x", auditPayloadLimitBytes+10)))
	if !truncated || len(body) != auditPayloadLimitBytes {
		t.Fatalf("unexpected bounded body: truncated=%v len=%d", truncated, len(body))
	}
}
