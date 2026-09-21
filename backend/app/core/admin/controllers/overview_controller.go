package controllers

import (
	"github.com/goravel/framework/contracts/http"

	"goravel/app/facades"
)

type OverviewController struct{}

func NewOverviewController() *OverviewController { return &OverviewController{} }

func (o *OverviewController) Index(ctx http.Context) http.Response {
	users, err := facades.Orm().Query().Table("users").Count()
	if err != nil {
		return ctx.Response().Status(500).Json(http.Json{"code": "INTERNAL_ERROR"})
	}
	roles, err := facades.Orm().Query().Table("roles").Count()
	if err != nil {
		return ctx.Response().Status(500).Json(http.Json{"code": "INTERNAL_ERROR"})
	}
	permissions, err := facades.Orm().Query().Table("permissions").Count()
	if err != nil {
		return ctx.Response().Status(500).Json(http.Json{"code": "INTERNAL_ERROR"})
	}
	return ctx.Response().Success().Json(http.Json{"data": http.Json{
		"users": users, "roles": roles, "permissions": permissions,
	}})
}
