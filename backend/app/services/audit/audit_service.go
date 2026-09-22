package auditservices

import (
	"strings"
	"time"

	"goravel/app/facades"
)

type AuditService struct{}

func NewAuditService() *AuditService { return &AuditService{} }

func (s *AuditService) Record(userID uint, action string, metadata map[string]any) error {
	return s.record(&userID, action, metadata)
}

func (s *AuditService) RecordHTTP(userID *uint, input HTTPAuditInput) error {
	return s.record(userID, "http."+strings.ToUpper(input.Method), BuildHTTPAuditMetadata(input))
}

func (s *AuditService) record(userID *uint, action string, metadata map[string]any) error {
	encoded, _, err := MarshalBounded(metadata)
	if err != nil {
		return err
	}
	values := map[string]any{
		"action":     action,
		"metadata":   string(encoded),
		"created_at": time.Now(),
		"updated_at": time.Now(),
	}
	if userID != nil {
		values["user_id"] = *userID
	}
	return facades.Orm().Query().Table("audit_logs").Create(&values)
}
