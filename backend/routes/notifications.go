package routes

import (
	admincontrollers "goravel/app/core/admin/controllers"
	"goravel/app/facades"
	adminmiddleware "goravel/app/http/middleware"
)

func registerNotificationRoutes() {
	notificationController := admincontrollers.NewNotificationController()
	facades.Route().Middleware(adminmiddleware.RequireAuthentication()).Get("/api/v1/notifications", notificationController.Index)
	facades.Route().Middleware(adminmiddleware.RequireAuthentication()).Get("/api/v1/notifications/unread-count", notificationController.UnreadCount)
	facades.Route().Middleware(adminmiddleware.RequireAuthentication()).Put("/api/v1/notifications/read-all", notificationController.MarkAllRead)
	facades.Route().Middleware(adminmiddleware.RequireAuthentication()).Put("/api/v1/notifications/{id}/read", notificationController.MarkRead)
}
