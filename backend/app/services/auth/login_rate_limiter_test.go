package authservices

import "testing"

func TestLoginRateLimitKeyNormalizesIdentity(t *testing.T) {
	first := loginRateLimitKey(" Admin@Example.com ", "127.0.0.1")
	second := loginRateLimitKey("admin@example.com", "127.0.0.1")
	if first != second {
		t.Fatalf("expected normalized identities to share a key, got %q and %q", first, second)
	}
	if first == "admin@example.com:127.0.0.1" {
		t.Fatal("rate-limit key must not store the raw email and IP")
	}
}
