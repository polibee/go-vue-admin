package auditservices

import (
	"errors"
	"time"

	"goravel/app/facades"
)

const (
	CleanupModeRetention = "retention"
	CleanupModeAll       = "all"

	DefaultAuditRetentionDays = 365
	MaxAuditRetentionDays     = 3650
)

var ErrInvalidRetentionDays = errors.New("invalid audit retention days")
var ErrInvalidCleanupMode = errors.New("invalid audit cleanup mode")

func ValidateCleanupMode(mode string) error {
	if mode != CleanupModeRetention && mode != CleanupModeAll {
		return ErrInvalidCleanupMode
	}
	return nil
}

func ValidateRetentionDays(days int) error {
	if days < 1 || days > MaxAuditRetentionDays {
		return ErrInvalidRetentionDays
	}
	return nil
}

func (s *AuditService) Cleanup(retentionDays int) (int64, time.Time, error) {
	if err := ValidateRetentionDays(retentionDays); err != nil {
		return 0, time.Time{}, err
	}
	cutoff := time.Now().UTC().AddDate(0, 0, -retentionDays)
	deleted, err := s.CleanupBefore(cutoff)
	return deleted, cutoff, err
}

func (s *AuditService) CleanupBefore(cutoff time.Time) (int64, error) {
	result, err := facades.Orm().Query().Table("audit_logs").Where("created_at < ?", cutoff).Delete()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected, nil
}

func (s *AuditService) CleanupAll() (int64, error) {
	result, err := facades.Orm().Query().Table("audit_logs").Delete()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected, nil
}
