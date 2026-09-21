package routes

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/support"

	admincontrollers "goravel/app/core/admin/controllers"
	authcontrollers "goravel/app/core/auth/controllers"
	"goravel/app/facades"
	adminmiddleware "goravel/app/http/middleware"
	usercontrollers "goravel/app/modules/users/controllers"
	"goravel/app/openapi"
)

func Web() {
	facades.Route().Get("/", func(ctx http.Context) http.Response {
		return ctx.Response().View().Make("welcome.tmpl", map[string]any{
			"version": support.Version,
		})
	})

	facades.Route().Static("public", "./public")
	facades.Route().Get("/api/openapi.json", func(ctx http.Context) http.Response {
		return ctx.Response().Json(200, http.Json(openapi.Spec()))
	})
	if facades.Config().GetString("app.env", "production") != "production" {
		facades.Route().Get("/api/docs", func(ctx http.Context) http.Response {
			return ctx.Response().Header("Content-Type", "text/html; charset=utf-8").String(200, openapi.DocsHTML())
		})
	}

	userController := usercontrollers.NewUserController()
	facades.Route().Get("/users", userController.Index)

	authController := authcontrollers.NewAuthController()
	facades.Route().Post("/api/v1/auth/login", authController.Login)
	facades.Route().Get("/api/v1/auth/me", authController.Me)
	facades.Route().Post("/api/v1/auth/refresh", authController.Refresh)
	facades.Route().Post("/api/v1/auth/logout", authController.Logout)
	facades.Route().Post("/api/v1/auth/logout-all", authController.LogoutAll)

	rbacController := admincontrollers.NewRBACController()
	overviewController := admincontrollers.NewOverviewController()
	auditController := admincontrollers.NewAuditController()
	resourceController := admincontrollers.NewResourceController()
	settingsController := admincontrollers.NewSettingsController()
	facades.Route().Middleware(adminmiddleware.RequirePermission("admin.users.view")).Get("/api/v1/admin/registry", resourceController.Index)
	facades.Route().Middleware(adminmiddleware.RequirePermission("admin.users.view")).Get("/api/v1/admin/overview", overviewController.Index)
	facades.Route().Middleware(adminmiddleware.RequirePermission("admin.users.view")).Get("/api/v1/admin/audit-logs", auditController.Index)
	facades.Route().Middleware(adminmiddleware.RequirePermission("admin.settings.manage")).Get("/api/v1/admin/settings", settingsController.Index)
	facades.Route().Middleware(adminmiddleware.RequirePermission("admin.settings.manage")).Put("/api/v1/admin/settings/:key", settingsController.Upsert)
	facades.Route().Middleware(adminmiddleware.RequirePermission("admin.users.manage")).Put("/api/v1/admin/users/status", rbacController.BulkSetUserStatus)
	facades.Route().Middleware(adminmiddleware.RequirePermission("admin.roles.manage")).Put("/api/v1/admin/roles/:id/permissions", rbacController.ReplaceRolePermissions)
	facades.Route().Middleware(adminmiddleware.RequirePermission("admin.roles.manage")).Get("/api/v1/admin/users/:id/roles", rbacController.UserRoles)
	facades.Route().Middleware(adminmiddleware.RequirePermission("admin.roles.manage")).Put("/api/v1/admin/users/:id/roles", rbacController.ReplaceUserRoles)
	facades.Route().Middleware(adminmiddleware.RequireResourcePermission("view")).Get("/api/v1/admin/:resource", resourceController.List)
	facades.Route().Middleware(adminmiddleware.RequireResourcePermission("view")).Get("/api/v1/admin/:resource/export", resourceController.Export)
	facades.Route().Middleware(adminmiddleware.RequireResourcePermission("view")).Get("/api/v1/admin/:resource/:id", resourceController.Show)
	facades.Route().Middleware(adminmiddleware.RequireResourcePermission("create")).Post("/api/v1/admin/:resource", resourceController.Create)
	facades.Route().Middleware(adminmiddleware.RequireResourcePermission("update")).Put("/api/v1/admin/:resource/:id", resourceController.Update)
	facades.Route().Middleware(adminmiddleware.RequireResourcePermission("delete")).Delete("/api/v1/admin/:resource/:id", resourceController.Delete)
}
