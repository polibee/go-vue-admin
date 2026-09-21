package routes

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/support"

	"goravel/app/facades"
	"goravel/app/http/controllers"
	adminmiddleware "goravel/app/http/middleware"
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

	userController := controllers.NewUserController()
	facades.Route().Get("/users", userController.Index)

	authController := controllers.NewAuthController()
	facades.Route().Post("/api/v1/auth/login", authController.Login)
	facades.Route().Get("/api/v1/auth/me", authController.Me)
	facades.Route().Post("/api/v1/auth/refresh", authController.Refresh)
	facades.Route().Post("/api/v1/auth/logout", authController.Logout)
	facades.Route().Post("/api/v1/auth/logout-all", authController.LogoutAll)

	rbacController := controllers.NewRBACController()
	overviewController := controllers.NewOverviewController()
	resourceController := controllers.NewResourceController()
	facades.Route().Middleware(adminmiddleware.RequirePermission("admin.users.view")).Get("/api/v1/admin/resources", resourceController.Index)
	facades.Route().Middleware(adminmiddleware.RequirePermission("admin.users.view")).Get("/api/v1/admin/overview", overviewController.Index)
	facades.Route().Middleware(adminmiddleware.RequirePermission("admin.users.view")).Get("/api/v1/admin/resources/:resource", resourceController.List)
	facades.Route().Middleware(adminmiddleware.RequirePermission("admin.users.view")).Get("/api/v1/admin/resources/:resource/:id", resourceController.Show)
	facades.Route().Middleware(adminmiddleware.RequirePermission("admin.users.view")).Get("/api/v1/admin/users", rbacController.Users)
	facades.Route().Middleware(adminmiddleware.RequirePermission("admin.users.manage")).Post("/api/v1/admin/users", rbacController.CreateUser)
	facades.Route().Middleware(adminmiddleware.RequirePermission("admin.users.manage")).Put("/api/v1/admin/users/status", rbacController.BulkSetUserStatus)
	facades.Route().Middleware(adminmiddleware.RequirePermission("admin.users.manage")).Put("/api/v1/admin/users/:id", rbacController.UpdateUser)
	facades.Route().Middleware(adminmiddleware.RequirePermission("admin.users.manage")).Delete("/api/v1/admin/users/:id", rbacController.DeleteUser)
	facades.Route().Middleware(adminmiddleware.RequirePermission("admin.roles.manage")).Get("/api/v1/admin/roles", rbacController.Roles)
	facades.Route().Middleware(adminmiddleware.RequirePermission("admin.permissions.manage")).Get("/api/v1/admin/permissions", rbacController.Permissions)
	facades.Route().Middleware(adminmiddleware.RequirePermission("admin.roles.manage")).Post("/api/v1/admin/roles", rbacController.CreateRole)
	facades.Route().Middleware(adminmiddleware.RequirePermission("admin.roles.manage")).Get("/api/v1/admin/roles/:id", rbacController.ShowRole)
	facades.Route().Middleware(adminmiddleware.RequirePermission("admin.roles.manage")).Put("/api/v1/admin/roles/:id", rbacController.UpdateRole)
	facades.Route().Middleware(adminmiddleware.RequirePermission("admin.roles.manage")).Delete("/api/v1/admin/roles/:id", rbacController.DeleteRole)
	facades.Route().Middleware(adminmiddleware.RequirePermission("admin.roles.manage")).Put("/api/v1/admin/roles/:id/permissions", rbacController.ReplaceRolePermissions)
	facades.Route().Middleware(adminmiddleware.RequirePermission("admin.roles.manage")).Get("/api/v1/admin/users/:id/roles", rbacController.UserRoles)
	facades.Route().Middleware(adminmiddleware.RequirePermission("admin.roles.manage")).Put("/api/v1/admin/users/:id/roles", rbacController.ReplaceUserRoles)
}
