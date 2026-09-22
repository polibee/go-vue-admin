package notifications

import "testing"

func TestValidateNotificationURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		ok   bool
	}{
		{name: "empty", url: "", ok: true},
		{name: "internal path", url: "/users/1", ok: true},
		{name: "query path", url: "/audit-logs?unread=1", ok: true},
		{name: "absolute URL", url: "https://example.com", ok: false},
		{name: "protocol relative", url: "//example.com", ok: false},
		{name: "javascript URL", url: "javascript:alert(1)", ok: false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := validateNotificationURL(test.url); (got == nil) != test.ok {
				t.Fatalf("validateNotificationURL(%q) error = %v, want valid = %v", test.url, got, test.ok)
			}
		})
	}
}

func TestNotificationInputRejectsEmptyUser(t *testing.T) {
	if err := validateNotificationInput(0, NotificationInput{Title: "Hello", Body: "Body"}); err == nil {
		t.Fatal("expected zero user id to be rejected")
	}
}
