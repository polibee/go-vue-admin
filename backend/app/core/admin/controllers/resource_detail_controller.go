package controllers

import (
	"github.com/goravel/framework/contracts/http"

	"goravel/app/facades"
	"goravel/app/models"
	"goravel/app/modules/admin/registry"
	rbacservices "goravel/app/services/rbac"
)

func resourceRecordLoaded(id uint) bool { return id > 0 }

func (r *ResourceController) Show(ctx http.Context) http.Response {
	id := resourceID(ctx)
	if id < 1 {
		return ctx.Response().Status(404).Json(http.Json{"code": "RESOURCE_NOT_FOUND"})
	}

	resourceName := resourceName(ctx)
	manifest, manifestErr := registry.AdminRegistry().Find(resourceName)
	if manifestErr != nil || manifest.Table == "" {
		return ctx.Response().Status(404).Json(http.Json{"code": "RESOURCE_NOT_FOUND"})
	}
	fieldPolicies, fieldErr := resourceFieldPoliciesFor(ctx, manifest, "view")
	if fieldErr != nil {
		return resourceScopeError(ctx, fieldErr)
	}
	switch resourceName {
	case "users":
		var user models.User
		q, scopeErr := applyResourceScope(ctx, facades.Orm().Query(), manifest, "view")
		if scopeErr != nil {
			return resourceScopeError(ctx, scopeErr)
		}
		if err := q.Where("id = ?", id).First(&user); err != nil || !resourceRecordLoaded(user.ID) {
			return ctx.Response().Status(404).Json(http.Json{"code": "RESOURCE_NOT_FOUND"})
		}
		return ctx.Response().Success().Json(http.Json{"data": projectResourceValueWithPolicies(user.Public(), manifest, fieldPolicies, false)})
	case "roles":
		var role models.Role
		q, scopeErr := applyResourceScope(ctx, facades.Orm().Query(), manifest, "view")
		if scopeErr != nil {
			return resourceScopeError(ctx, scopeErr)
		}
		if err := q.Where("id = ?", id).First(&role); err != nil || !resourceRecordLoaded(role.ID) {
			return ctx.Response().Status(404).Json(http.Json{"code": "RESOURCE_NOT_FOUND"})
		}
		value := projectResourceValueWithPolicies(role, manifest, fieldPolicies, false)
		permissions, permissionErr := rbacservices.NewRoleService().PermissionAssignments(int64(role.ID))
		if permissionErr != nil {
			return ctx.Response().Status(500).Json(http.Json{"code": "RBAC_PERMISSIONS_ERROR"})
		}
		value["permissions"] = permissions
		return ctx.Response().Success().Json(http.Json{"data": value})
	case "permissions":
		var permission models.Permission
		q, scopeErr := applyResourceScope(ctx, facades.Orm().Query(), manifest, "view")
		if scopeErr != nil {
			return resourceScopeError(ctx, scopeErr)
		}
		if err := q.Where("id = ?", id).First(&permission); err != nil || !resourceRecordLoaded(permission.ID) {
			return ctx.Response().Status(404).Json(http.Json{"code": "RESOURCE_NOT_FOUND"})
		}
		return ctx.Response().Success().Json(http.Json{"data": projectResourceValueWithPolicies(permission, manifest, fieldPolicies, false)})
	default:
		var rows []map[string]any
		q, scopeErr := applyResourceScope(ctx, facades.Orm().Query().Table(manifest.Table), manifest, "view")
		if scopeErr != nil {
			return resourceScopeError(ctx, scopeErr)
		}
		if err := q.Where("id = ?", id).Get(&rows); err != nil || len(rows) == 0 {
			return ctx.Response().Status(404).Json(http.Json{"code": "RESOURCE_NOT_FOUND"})
		}
		return ctx.Response().Success().Json(http.Json{"data": projectResourceRecordWithPolicies(rows[0], manifest, fieldPolicies, false)})
	}
}
