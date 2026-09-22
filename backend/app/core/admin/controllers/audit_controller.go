package controllers

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/goravel/framework/contracts/http"

	"goravel/app/facades"
	auditservices "goravel/app/services/audit"
)

type AuditController struct{}

func NewAuditController() *AuditController { return &AuditController{} }

func (a *AuditController) Index(ctx http.Context) http.Response {
	query := resourceListQuery{
		Page:    positiveInt(ctx.Request().Query("page", "1"), 1),
		PerPage: positiveInt(ctx.Request().Query("per_page", "20"), 20),
	}
	if query.PerPage > 100 {
		query.PerPage = 100
	}
	action := strings.TrimSpace(ctx.Request().Query("action"))
	userID := strings.TrimSpace(ctx.Request().Query("user_id"))

	q := facades.Orm().Query().Table("audit_logs")
	if action != "" {
		q = q.Where("action = ?", action)
	}
	if userID != "" {
		q = q.Where("user_id = ?", userID)
	}
	q = q.OrderByDesc("id")

	var entries []map[string]any
	var total int64
	if err := q.Paginate(query.Page, query.PerPage, &entries, &total); err != nil {
		return ctx.Response().Status(500).Json(http.Json{"code": "INTERNAL_ERROR"})
	}
	return resourceListResponse(ctx, entries, query, total)
}

type auditCleanupPayload struct {
	RetentionDays int `json:"retention_days"`
}

func (a *AuditController) Cleanup(ctx http.Context) http.Response {
	var payload auditCleanupPayload
	if err := ctx.Request().Bind(&payload); err != nil {
		return ctx.Response().Status(422).Json(http.Json{"code": "AUDIT_RETENTION_INVALID"})
	}
	if payload.RetentionDays == 0 {
		payload.RetentionDays = auditservices.DefaultAuditRetentionDays
	}
	deleted, cutoff, err := auditservices.NewAuditService().Cleanup(payload.RetentionDays)
	if err != nil {
		if errors.Is(err, auditservices.ErrInvalidRetentionDays) {
			return ctx.Response().Status(422).Json(http.Json{"code": "AUDIT_RETENTION_INVALID"})
		}
		return ctx.Response().Status(500).Json(http.Json{"code": "AUDIT_CLEANUP_FAILED"})
	}
	if identity, identityErr := facades.Auth(ctx).ID(); identityErr == nil {
		userID, parseErr := strconv.ParseUint(identity, 10, 32)
		if parseErr == nil && userID > 0 {
			_ = auditservices.NewAuditService().Record(uint(userID), "audit.cleanup", map[string]any{
				"retention_days": payload.RetentionDays,
				"deleted":        deleted,
				"cutoff":         cutoff.Format(time.RFC3339),
			})
		}
	}
	return ctx.Response().Success().Json(http.Json{"data": http.Json{
		"deleted":        deleted,
		"retention_days": payload.RetentionDays,
		"cutoff":         cutoff.Format(time.RFC3339),
	}})
}
