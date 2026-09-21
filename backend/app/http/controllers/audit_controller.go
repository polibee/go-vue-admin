package controllers

import (
	"strings"

	"github.com/goravel/framework/contracts/http"

	"goravel/app/facades"
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
