package routes

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/support"

	permissionresources "goravel/app/core/permission/resources"
	"goravel/app/core/resource"
	roleresources "goravel/app/core/role/resources"
	userresources "goravel/app/core/user/resources"
	"goravel/app/facades"
	"goravel/app/http/controllers"
)

func Web() {
	facades.Route().Get("/", func(ctx http.Context) http.Response {
		return ctx.Response().View().Make("welcome.tmpl", map[string]any{
			"version": support.Version,
		})
	})

	facades.Route().Static("public", "./public")

	userController := controllers.NewUserController()
	facades.Route().Get("/users", userController.Index)

	authController := controllers.NewAuthController()
	menuController := controllers.NewMenuController(authController)
	resourceController := controllers.NewConfiguredResourceController(authController)
	userResourceController := controllers.NewCoreResourceController(authController, resource.NewCoreResourceService(resource.NewMemoryCoreResourceRepository(userresources.Seed()...)), "users")
	roleResourceController := controllers.NewCoreResourceController(authController, resource.NewCoreResourceService(resource.NewMemoryCoreResourceRepository(roleresources.Seed()...)), "roles")
	permissionResourceController := controllers.NewCoreResourceController(authController, resource.NewCoreResourceService(resource.NewMemoryCoreResourceRepository(permissionresources.Seed()...)), "permissions")
	settingController := controllers.NewConfiguredSettingController(authController)
	mediaController := controllers.NewConfiguredMediaController(authController)
	for _, prefix := range []string{"", "/api"} {
		facades.Route().Get(prefix+"/csrf", authController.CSRF)
		facades.Route().Get(prefix+"/auth/bootstrap", authController.Bootstrap)
		facades.Route().Post(prefix+"/login", authController.Login)
		facades.Route().Post(prefix+"/logout", authController.Logout)
		facades.Route().Get(prefix+"/me", authController.Me)
	}
	facades.Route().Get("/menu", menuController.Index)
	facades.Route().Get("/api/menu", menuController.Index)
	facades.Route().Get("/api/settings", settingController.Index)
	facades.Route().Post("/api/settings", settingController.Store)
	facades.Route().Put("/api/settings/:namespace/:key", settingController.Update)
	facades.Route().Delete("/api/settings/:namespace/:key", settingController.Destroy)
	facades.Route().Get("/api/media", mediaController.Index)
	facades.Route().Post("/api/media", mediaController.Store)
	facades.Route().Get("/api/media/:id/preview", mediaController.Preview)
	facades.Route().Delete("/api/media/:id", mediaController.Destroy)
	facades.Route().Get("/api/resources/demo", resourceController.Index)
	facades.Route().Get("/api/resources/demo/:id", resourceController.Show)
	facades.Route().Post("/api/resources/demo", resourceController.Store)
	facades.Route().Put("/api/resources/demo/:id", resourceController.Update)
	facades.Route().Delete("/api/resources/demo/:id", resourceController.Destroy)
	facades.Route().Post("/api/resources/demo/bulk-delete", resourceController.BulkDestroy)
	registerCoreResourceRoutes("users", userResourceController)
	registerCoreResourceRoutes("roles", roleResourceController)
	registerCoreResourceRoutes("permissions", permissionResourceController)
}

func registerCoreResourceRoutes(prefix string, controller *controllers.CoreResourceController) {
	base := "/api/resources/" + prefix
	facades.Route().Get(base, controller.Index)
	facades.Route().Get(base+"/:id", controller.Show)
	facades.Route().Post(base, controller.Store)
	facades.Route().Put(base+"/:id", controller.Update)
	facades.Route().Delete(base+"/:id", controller.Destroy)
	facades.Route().Post(base+"/bulk-delete", controller.BulkDestroy)
}
