package routes

import (
	authcontrollers "goravel/app/core/auth/controllers"
	"goravel/app/facades"
)

func registerAuthRoutes() {
	authController := authcontrollers.NewAuthController()
	facades.Route().Post("/api/v1/auth/login", authController.Login)
	facades.Route().Get("/api/v1/auth/me", authController.Me)
	facades.Route().Post("/api/v1/auth/refresh", authController.Refresh)
	facades.Route().Post("/api/v1/auth/logout", authController.Logout)
	facades.Route().Post("/api/v1/auth/logout-all", authController.LogoutAll)
}
