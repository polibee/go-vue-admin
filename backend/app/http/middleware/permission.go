package middleware

import (
	"goravel/app/facades"
	"goravel/app/services"

	"github.com/goravel/framework/contracts/http"
)

type permissionMiddleware struct {
	permission string
}

func (m permissionMiddleware) Signature() string {
	return "admin:permission:" + m.permission
}

func (m permissionMiddleware) Handle(ctx http.Context) {
	if token := ctx.Request().Header("Authorization"); token == "" {
		ctx.Response().Status(401).Json(http.Json{"code": "AUTH_UNAUTHORIZED"}).Abort()
		return
	} else if _, err := facades.Auth(ctx).Parse(token); err != nil {
		ctx.Response().Status(401).Json(http.Json{"code": "AUTH_UNAUTHORIZED"}).Abort()
		return
	}

	allowed, err := services.NewRBACService().UserHasPermission(ctx, m.permission)
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

func RequirePermission(permission string) http.Middleware {
	return permissionMiddleware{permission: permission}
}
