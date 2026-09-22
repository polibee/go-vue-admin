package controllers

import (
	"strconv"

	"github.com/goravel/framework/contracts/http"

	"goravel/app/facades"
	"goravel/app/models"
	"goravel/app/modules/admin/registry"
)

func (r *ResourceController) Show(ctx http.Context) http.Response {
	id, err := strconv.ParseInt(ctx.Request().Route("id"), 10, 64)
	if err != nil || id < 1 {
		return ctx.Response().Status(404).Json(http.Json{"code": "RESOURCE_NOT_FOUND"})
	}

	resourceName := ctx.Request().Route("resource")
	manifest, manifestErr := registry.AdminRegistry().Find(resourceName)
	if manifestErr != nil || manifest.Table == "" {
		return ctx.Response().Status(404).Json(http.Json{"code": "RESOURCE_NOT_FOUND"})
	}
	switch resourceName {
	case "users":
		var user models.User
		q, scopeErr := applyResourceScope(ctx, facades.Orm().Query(), manifest, "view")
		if scopeErr != nil {
			return resourceScopeError(ctx, scopeErr)
		}
		if err := q.Where("id", id).First(&user); err != nil {
			return ctx.Response().Status(404).Json(http.Json{"code": "RESOURCE_NOT_FOUND"})
		}
		return ctx.Response().Success().Json(http.Json{"data": user.Public()})
	case "roles":
		var role models.Role
		q, scopeErr := applyResourceScope(ctx, facades.Orm().Query(), manifest, "view")
		if scopeErr != nil {
			return resourceScopeError(ctx, scopeErr)
		}
		if err := q.Where("id", id).First(&role); err != nil {
			return ctx.Response().Status(404).Json(http.Json{"code": "RESOURCE_NOT_FOUND"})
		}
		return ctx.Response().Success().Json(http.Json{"data": role})
	case "permissions":
		var permission models.Permission
		q, scopeErr := applyResourceScope(ctx, facades.Orm().Query(), manifest, "view")
		if scopeErr != nil {
			return resourceScopeError(ctx, scopeErr)
		}
		if err := q.Where("id", id).First(&permission); err != nil {
			return ctx.Response().Status(404).Json(http.Json{"code": "RESOURCE_NOT_FOUND"})
		}
		return ctx.Response().Success().Json(http.Json{"data": permission})
	default:
		var rows []map[string]any
		q, scopeErr := applyResourceScope(ctx, facades.Orm().Query().Table(manifest.Table), manifest, "view")
		if scopeErr != nil {
			return resourceScopeError(ctx, scopeErr)
		}
		if err := q.Where("id = ?", id).Get(&rows); err != nil || len(rows) == 0 {
			return ctx.Response().Status(404).Json(http.Json{"code": "RESOURCE_NOT_FOUND"})
		}
		return ctx.Response().Success().Json(http.Json{"data": rows[0]})
	}
}
