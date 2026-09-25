package routes

import (
	"github.com/goravel/framework/contracts/http"
	"goravel/app/facades"
	"goravel/app/openapi"
)

func registerOpenAPIRoutes() {
	adminContract := func(ctx http.Context) http.Response {
		return ctx.Response().Json(200, http.Json(openapi.AdminSpec()))
	}

	// The explicit Admin contract is the only OpenAPI contract exposed by this backend.
	facades.Route().Get("/api/openapi/admin.json", adminContract)
	if facades.Config().GetString("app.env", "production") != "production" {
		facades.Route().Get("/api/docs", func(ctx http.Context) http.Response {
			return ctx.Response().Header("Content-Type", "text/html; charset=utf-8").String(200, openapi.DocsHTML())
		})
	}
}
