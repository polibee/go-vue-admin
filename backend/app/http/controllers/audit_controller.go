package controllers

import (
	"github.com/goravel/framework/contracts/http"

	"goravel/app/facades"
)

type AuditController struct{}

func NewAuditController() *AuditController { return &AuditController{} }

func (a *AuditController) Index(ctx http.Context) http.Response {
	var entries []map[string]any
	if err := facades.Orm().Query().Table("audit_logs").OrderByDesc("id").Limit(100).Get(&entries); err != nil {
		return ctx.Response().Status(500).Json(http.Json{"code": "INTERNAL_ERROR"})
	}
	return ctx.Response().Success().Json(http.Json{"data": entries})
}
