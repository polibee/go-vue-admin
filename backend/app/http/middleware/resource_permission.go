package middleware

import (
	"goravel/app/facades"
	"goravel/app/modules/admin/registry"
	rbacservices "goravel/app/services/rbac"

	"github.com/goravel/framework/contracts/http"
)

func RequireResourcePermission(action string) http.Middleware {
	return resourcePermissionMiddleware{action: action}
}

type resourcePermissionMiddleware struct{ action string }

func (m resourcePermissionMiddleware) Signature() string {
	return "admin:resource-permission:" + m.action
}

func (m resourcePermissionMiddleware) Handle(ctx http.Context) {
	if token := ctx.Request().Header("Authorization"); token == "" {
		ctx.Response().Status(401).Json(http.Json{"code": "AUTH_UNAUTHORIZED"}).Abort()
		return
	} else if _, err := facades.Auth(ctx).Parse(token); err != nil {
		ctx.Response().Status(401).Json(http.Json{"code": "AUTH_UNAUTHORIZED"}).Abort()
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
