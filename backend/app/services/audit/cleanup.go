package auditservices

import (
	"errors"
	"time"

	"goravel/app/facades"
)

const (
	CleanupModeRetention = "retention"
	CleanupModeAll       = "all"
	CleanupModeSelected  = "selected"
	CleanupModeFiltered  = "filtered"
	CleanupConfirmation  = "DELETE"

	DefaultAuditRetentionDays = 365
	MaxAuditRetentionDays     = 3650
)

var ErrInvalidRetentionDays = errors.New("invalid audit retention days")
var ErrInvalidCleanupMode = errors.New("invalid audit cleanup mode")

func ValidateCleanupMode(mode string) error {
	if mode != CleanupModeRetention && mode != CleanupModeAll && mode != CleanupModeSelected && mode != CleanupModeFiltered {
		return ErrInvalidCleanupMode
	}
	return nil
}

func ValidateCleanupConfirmation(mode, confirmation string) error {
	if mode == CleanupModeRetention {
		return nil
	}
	if confirmation != CleanupConfirmation {
		return errors.New("invalid audit cleanup confirmation")
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

func (s *AuditService) CleanupSelected(ids []int64) (int64, error) {
	if len(ids) == 0 {
		return 0, errors.New("audit cleanup requires at least one selected id")
	}
	values := make([]any, 0, len(ids))
	for _, id := range ids {
		if id < 1 {
			return 0, errors.New("audit cleanup contains an invalid id")
		}
		values = append(values, id)
	}
	result, err := facades.Orm().Query().Table("audit_logs").WhereIn("id", values).Delete()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected, nil
}

func (s *AuditService) CleanupFiltered(action, userID string) (int64, error) {
	query := facades.Orm().Query().Table("audit_logs")
	if action != "" && action != "all" {
		query = query.Where("action = ?", action)
	}
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	result, err := query.Delete()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected, nil
}
