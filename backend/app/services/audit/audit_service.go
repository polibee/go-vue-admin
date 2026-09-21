package auditservices

import (
	"encoding/json"
	"time"

	"goravel/app/facades"
)

type AuditService struct{}

func NewAuditService() *AuditService { return &AuditService{} }

func (s *AuditService) Record(userID uint, action string, metadata map[string]any) error {
	encoded, err := json.Marshal(metadata)
	if err != nil {
		return err
	}
	return facades.Orm().Query().Table("audit_logs").Create(&map[string]any{
		"user_id":    userID,
		"action":     action,
		"metadata":   string(encoded),
		"created_at": time.Now(),
		"updated_at": time.Now(),
	})
}
