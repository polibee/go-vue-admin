package routes

import (
	admincontrollers "goravel/app/core/admin/controllers"
	"goravel/app/facades"
	adminmiddleware "goravel/app/http/middleware"
	"goravel/app/modules/admin/registry"
)

func registerAdminRoutes() {
	rbacController := admincontrollers.NewRBACController()
	overviewController := admincontrollers.NewOverviewController()
	auditController := admincontrollers.NewAuditController()
	resourceController := admincontrollers.NewResourceController()
	globalSearchController := admincontrollers.NewGlobalSearchController()
	settingsController := admincontrollers.NewSettingsController()

	facades.Route().Middleware(adminmiddleware.RequireAnyResourcePermission()).Get(adminAPIBase+"/registry", resourceController.Index)
	facades.Route().Middleware(adminmiddleware.RequireAnyResourcePermission()).Get(adminAPIBase+"/search", globalSearchController.Index)
	facades.Route().Middleware(adminmiddleware.RequirePermission("admin.users.view")).Get(adminAPIBase+"/overview", overviewController.Index)
	facades.Route().Middleware(adminmiddleware.RequirePermission("admin.users.view")).Get(adminAPIBase+"/audit-logs", auditController.Index)
	facades.Route().Middleware(adminmiddleware.RequirePermission("admin.settings.manage")).Post(adminAPIBase+"/audit-logs/cleanup", auditController.Cleanup)
	facades.Route().Middleware(adminmiddleware.RequirePermission("admin.settings.manage")).Get(adminAPIBase+"/settings", settingsController.Index)
	facades.Route().Middleware(adminmiddleware.RequirePermission("admin.settings.manage")).Put(adminAPIBase+"/settings/{key}", settingsController.Upsert)
	facades.Route().Middleware(adminmiddleware.RequirePermission("admin.roles.manage")).Put(adminAPIBase+"/roles/{id}/permissions", rbacController.ReplaceRolePermissions)
	facades.Route().Middleware(adminmiddleware.RequirePermission("admin.roles.manage")).Get(adminAPIBase+"/users/{id}/roles", rbacController.UserRoles)
	facades.Route().Middleware(adminmiddleware.RequirePermission("admin.roles.manage")).Put(adminAPIBase+"/users/{id}/roles", rbacController.ReplaceUserRoles)
	facades.Route().Middleware(adminmiddleware.RequireResourcePermission("view")).Get(adminAPIBase+"/{resource}", resourceController.List)
	facades.Route().Middleware(adminmiddleware.RequireResourcePermission("view")).Get(adminAPIBase+"/{resource}/export", resourceController.Export)
	facades.Route().Middleware(adminmiddleware.RequireAuthentication()).Post(adminAPIBase+"/{resource}/actions/{action}", resourceController.Action)
	facades.Route().Middleware(adminmiddleware.RequireAuthentication()).Get(adminAPIBase+"/{resource}/relations/{relation}/options", resourceController.RelationOptions)
	facades.Route().Middleware(adminmiddleware.RequireAuthentication()).Get(adminAPIBase+"/{resource}/{id}/relations/{relation}", resourceController.RelationRecords)
	facades.Route().Middleware(adminmiddleware.RequireResourcePermission("view")).Get(adminAPIBase+"/{resource}/{id}", resourceController.Show)
	facades.Route().Middleware(adminmiddleware.RequireResourcePermission("create")).Post(adminAPIBase+"/{resource}", resourceController.Create)
	facades.Route().Middleware(adminmiddleware.RequireResourcePermission("update")).Put(adminAPIBase+"/{resource}/{id}", resourceController.Update)
	facades.Route().Middleware(adminmiddleware.RequireResourcePermission("delete")).Delete(adminAPIBase+"/{resource}/{id}", resourceController.Delete)

	for _, manifest := range registry.AdminRegistry().All() {
		base := adminAPIBase + "/" + manifest.Name
		facades.Route().Middleware(adminmiddleware.RequireResourcePermission("view")).Get(base+"/export", resourceController.Export)
		facades.Route().Middleware(adminmiddleware.RequireResourcePermission("view")).Get(base+"/{id}", resourceController.Show)
		facades.Route().Middleware(adminmiddleware.RequireResourcePermission("update")).Put(base+"/{id}", resourceController.Update)
		facades.Route().Middleware(adminmiddleware.RequireResourcePermission("delete")).Delete(base+"/{id}", resourceController.Delete)
		facades.Route().Middleware(adminmiddleware.RequireAuthentication()).Post(base+"/actions/{action}", resourceController.Action)
		facades.Route().Middleware(adminmiddleware.RequireAuthentication()).Get(base+"/relations/{relation}/options", resourceController.RelationOptions)
		facades.Route().Middleware(adminmiddleware.RequireAuthentication()).Get(base+"/{id}/relations/{relation}", resourceController.RelationRecords)
	}
}
