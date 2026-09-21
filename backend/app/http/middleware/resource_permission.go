package middleware

import (
	"goravel/app/core/resource"
	"goravel/app/facades"
	"goravel/app/modules/admin/registry"
	rbacservices "goravel/app/services/rbac"

	"github.com/goravel/framework/contracts/http"
)

func RequireResourcePermission(action string) http.Middleware {
	return resourcePermissionMiddleware{action: action}
}

func RequireAnyResourcePermission() http.Middleware {
	return anyResourcePermissionMiddleware{}
}

type anyResourcePermissionMiddleware struct{}

func (m anyResourcePermissionMiddleware) Signature() string {
	return "admin:any-resource-permission"
}

func (m anyResourcePermissionMiddleware) Handle(ctx http.Context) {
	if !parseAuthenticatedRequest(ctx) {
		return
	}
	service := rbacservices.NewRBACService()
	for _, permission := range resourceViewPermissions(registry.AdminRegistry().All()) {
		allowed, err := service.UserHasPermission(ctx, permission)
		if err == nil && allowed {
			ctx.Request().Next()
			return
		}
	}
	ctx.Response().Status(403).Json(http.Json{"code": "RBAC_FORBIDDEN"}).Abort()
}

type resourcePermissionMiddleware struct{ action string }

func (m resourcePermissionMiddleware) Signature() string {
	return "admin:resource-permission:" + m.action
}

func (m resourcePermissionMiddleware) Handle(ctx http.Context) {
	if !parseAuthenticatedRequest(ctx) {
		return
	}
	manifest, err := registry.AdminRegistry().Find(ctx.Request().Route("resource"))
	if err != nil {
		ctx.Response().Status(404).Json(http.Json{"code": "RESOURCE_NOT_FOUND"}).Abort()
		return
	}
	permission := ""
	if m.action == "view" && len(manifest.Permissions) > 0 {
		permission = manifest.Permissions[0]
	}
	for _, action := range manifest.Actions {
		if action.Name == m.action {
			permission = action.Permission
			break
		}
	}
	if permission == "" {
		ctx.Response().Status(403).Json(http.Json{"code": "RBAC_FORBIDDEN"}).Abort()
		return
	}
	allowed, err := rbacservices.NewRBACService().UserHasPermission(ctx, permission)
	if err != nil {
		ctx.Response().Status(401).Json(http.Json{"code": "AUTH_UNAUTHORIZED"}).Abort()
		return
	}
	if !allowed {
		ctx.Response().Status(403).Json(http.Json{"code": "RBAC_FORBIDDEN"}).Abort()
		return
	}
	ctx.Request().Next()
}

func parseAuthenticatedRequest(ctx http.Context) bool {
	token := ctx.Request().Header("Authorization")
	if token == "" {
		ctx.Response().Status(401).Json(http.Json{"code": "AUTH_UNAUTHORIZED"}).Abort()
		return false
	}
	if _, err := facades.Auth(ctx).Parse(token); err != nil {
		ctx.Response().Status(401).Json(http.Json{"code": "AUTH_UNAUTHORIZED"}).Abort()
		return false
	}
	return true
}

func resourceViewPermissions(manifests []resource.Manifest) []string {
	permissions := make([]string, 0, len(manifests))
	seen := make(map[string]struct{})
	for _, manifest := range manifests {
		for _, permission := range manifest.Permissions {
			if permission == "" {
				continue
			}
			if _, exists := seen[permission]; exists {
				continue
			}
			seen[permission] = struct{}{}
			permissions = append(permissions, permission)
		}
	}
	return permissions
}
