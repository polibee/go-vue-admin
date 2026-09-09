package routes

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/support"

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
	for _, prefix := range []string{"", "/api"} {
		facades.Route().Get(prefix+"/csrf", authController.CSRF)
		facades.Route().Get(prefix+"/auth/bootstrap", authController.Bootstrap)
		facades.Route().Post(prefix+"/login", authController.Login)
		facades.Route().Post(prefix+"/logout", authController.Logout)
		facades.Route().Get(prefix+"/me", authController.Me)
	}
	facades.Route().Get("/menu", menuController.Index)
	facades.Route().Get("/api/menu", menuController.Index)
	facades.Route().Get("/api/resources/demo", resourceController.Index)
	facades.Route().Get("/api/resources/demo/:id", resourceController.Show)
	facades.Route().Post("/api/resources/demo", resourceController.Store)
	facades.Route().Put("/api/resources/demo/:id", resourceController.Update)
	facades.Route().Delete("/api/resources/demo/:id", resourceController.Destroy)
	facades.Route().Post("/api/resources/demo/bulk-delete", resourceController.BulkDestroy)
}
