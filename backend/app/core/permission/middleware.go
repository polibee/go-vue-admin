package permission

import (
	contractshttp "github.com/goravel/framework/contracts/http"

	apierrors "goravel/app/core/shared/errors"
)

type PermissionResolver func(ctx contractshttp.Context) []string

type permissionMiddleware struct {
	required string
	resolve  PermissionResolver
}

func Middleware(required string, resolve PermissionResolver) contractshttp.Middleware {
	return &permissionMiddleware{required: required, resolve: resolve}
}

func (m *permissionMiddleware) Signature() string {
	return "go-vue-admin:permission:" + m.required
}

func (m *permissionMiddleware) Handle(ctx contractshttp.Context) {
	var granted []string
	if m.resolve != nil {
		granted = m.resolve(ctx)
	}
	if NewAuthorizer().Require(granted, m.required) != nil {
		ctx.Request().AbortWithStatusJson(403, apierrors.New("FORBIDDEN", "没有执行该操作的权限", map[string]string{"permission": m.required}))
		return
	}
	ctx.Request().Next()
}
