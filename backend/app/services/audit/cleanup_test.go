package auditservices

import "testing"

func TestValidateRetentionDays(t *testing.T) {
	for _, days := range []int{1, 30, DefaultAuditRetentionDays, MaxAuditRetentionDays} {
		if err := ValidateRetentionDays(days); err != nil {
			t.Fatalf("retention days %d rejected: %v", days, err)
		}
	}
	for _, days := range []int{0, -1, MaxAuditRetentionDays + 1} {
		if err := ValidateRetentionDays(days); err == nil {
			t.Fatalf("retention days %d unexpectedly accepted", days)
		}
	}
}

func TestValidateCleanupMode(t *testing.T) {
	for _, mode := range []string{CleanupModeRetention, CleanupModeAll} {
		if err := ValidateCleanupMode(mode); err != nil {
			t.Fatalf("cleanup mode %q rejected: %v", mode, err)
		}
	}
	if err := ValidateCleanupMode("unknown"); err == nil {
		t.Fatal("unknown cleanup mode unexpectedly accepted")
	}
}
