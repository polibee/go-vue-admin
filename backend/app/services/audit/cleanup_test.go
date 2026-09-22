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
	for _, mode := range []string{CleanupModeRetention, CleanupModeAll, CleanupModeSelected, CleanupModeFiltered} {
		if err := ValidateCleanupMode(mode); err != nil {
			t.Fatalf("cleanup mode %q rejected: %v", mode, err)
		}
	}
	if err := ValidateCleanupMode("unknown"); err == nil {
		t.Fatal("unknown cleanup mode unexpectedly accepted")
	}
}

func TestValidateCleanupConfirmation(t *testing.T) {
	if err := ValidateCleanupConfirmation(CleanupModeRetention, ""); err != nil {
		t.Fatalf("retention confirmation should be optional: %v", err)
	}
	if err := ValidateCleanupConfirmation(CleanupModeAll, CleanupConfirmation); err != nil {
		t.Fatalf("expected valid destructive confirmation: %v", err)
	}
	if err := ValidateCleanupConfirmation(CleanupModeFiltered, "delete"); err == nil {
		t.Fatal("expected exact confirmation text to be required")
	}
}
