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
	Mode          string  `json:"mode"`
	RetentionDays int     `json:"retention_days"`
	IDs           []int64 `json:"ids"`
	Action        string  `json:"action"`
	UserID        string  `json:"user_id"`
	Confirmation  string  `json:"confirmation"`
}

func (a *AuditController) Cleanup(ctx http.Context) http.Response {
	var payload auditCleanupPayload
	if err := ctx.Request().Bind(&payload); err != nil {
		return ctx.Response().Status(422).Json(http.Json{"code": "AUDIT_RETENTION_INVALID"})
	}
	if payload.Mode == "" {
		payload.Mode = auditservices.CleanupModeRetention
	}
	if err := auditservices.ValidateCleanupMode(payload.Mode); err != nil {
		return ctx.Response().Status(422).Json(http.Json{"code": "AUDIT_CLEANUP_MODE_INVALID"})
	}
	if err := auditservices.ValidateCleanupConfirmation(payload.Mode, payload.Confirmation); err != nil {
		return ctx.Response().Status(422).Json(http.Json{"code": "AUDIT_CLEANUP_CONFIRMATION_REQUIRED"})
	}
	if payload.Mode == auditservices.CleanupModeSelected && len(payload.IDs) == 0 {
		return ctx.Response().Status(422).Json(http.Json{"code": "AUDIT_CLEANUP_SELECTION_REQUIRED"})
	}
	if payload.Mode == auditservices.CleanupModeFiltered && strings.TrimSpace(payload.Action) == "" && strings.TrimSpace(payload.UserID) == "" {
		return ctx.Response().Status(422).Json(http.Json{"code": "AUDIT_CLEANUP_FILTER_REQUIRED"})
	}
	if len(payload.IDs) > 1000 {
		return ctx.Response().Status(422).Json(http.Json{"code": "AUDIT_CLEANUP_SELECTION_TOO_LARGE"})
	}
	var deleted int64
	var cutoff time.Time
	var err error
	if payload.Mode == auditservices.CleanupModeAll {
		deleted, err = auditservices.NewAuditService().CleanupAll()
	} else if payload.Mode == auditservices.CleanupModeSelected {
		deleted, err = auditservices.NewAuditService().CleanupSelected(payload.IDs)
	} else if payload.Mode == auditservices.CleanupModeFiltered {
		deleted, err = auditservices.NewAuditService().CleanupFiltered(strings.TrimSpace(payload.Action), strings.TrimSpace(payload.UserID))
	} else {
		if payload.RetentionDays == 0 {
			payload.RetentionDays = auditservices.DefaultAuditRetentionDays
		}
		deleted, cutoff, err = auditservices.NewAuditService().Cleanup(payload.RetentionDays)
		if errors.Is(err, auditservices.ErrInvalidRetentionDays) {
			return ctx.Response().Status(422).Json(http.Json{"code": "AUDIT_RETENTION_INVALID"})
		}
	}
	if err != nil {
		return ctx.Response().Status(500).Json(http.Json{"code": "AUDIT_CLEANUP_FAILED"})
	}
	if identity, identityErr := facades.Auth(ctx).ID(); identityErr == nil {
		userID, parseErr := strconv.ParseUint(identity, 10, 32)
		if parseErr == nil && userID > 0 {
			metadata := map[string]any{
				"mode":           payload.Mode,
				"retention_days": payload.RetentionDays,
				"deleted":        deleted,
			}
			if payload.Mode == auditservices.CleanupModeSelected {
				metadata["selected_count"] = len(payload.IDs)
			}
			if payload.Mode == auditservices.CleanupModeFiltered {
				metadata["action"] = payload.Action
				metadata["user_id"] = payload.UserID
			}
			if payload.Mode == auditservices.CleanupModeRetention {
				metadata["cutoff"] = cutoff.Format(time.RFC3339)
			}
			_ = auditservices.NewAuditService().Record(uint(userID), "audit.cleanup", metadata)
		}
	}
	response := http.Json{"deleted": deleted, "mode": payload.Mode}
	if payload.Mode == auditservices.CleanupModeRetention {
		response["retention_days"] = payload.RetentionDays
		response["cutoff"] = cutoff.Format(time.RFC3339)
	} else {
		response["cutoff"] = nil
	}
	return ctx.Response().Success().Json(http.Json{"data": response})
}
