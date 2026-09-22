package corspolicy

import "testing"

func TestAllowedOrigin(t *testing.T) {
	if !AllowedOrigin("http://127.0.0.1:4174") {
		t.Fatal("expected the local admin preview origin to be allowed")
	}
	if !AllowedOrigin("http://127.0.0.1:5182") {
		t.Fatal("expected the isolated local admin verification origin to be allowed")
	}
	if AllowedOrigin("https://example.com") {
		t.Fatal("unexpectedly allowed an untrusted origin")
	}
}
