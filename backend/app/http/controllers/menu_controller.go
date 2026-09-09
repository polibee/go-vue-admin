package controllers

import (
	"github.com/goravel/framework/contracts/http"

	"goravel/app/core/menu"
	apierrors "goravel/app/core/shared/errors"
	"goravel/app/core/shared/response"
)

type MenuController struct {
	auth     *AuthController
	registry *menu.Registry
}

func NewMenuController(authController *AuthController) *MenuController {
	return &MenuController{
		auth: authController,
		registry: menu.NewRegistry(
			menu.Item{ID: "dashboard", Label: "仪表盘", Route: "/admin/dashboard", Permission: "dashboard.view"},
			menu.Item{ID: "settings", Label: "设置", Route: "/admin/settings", Permission: "settings.view"},
		),
	}
}

func (c *MenuController) Index(ctx http.Context) http.Response {
	currentUser, err := c.auth.CurrentUser(ctx)
	if err != nil {
		return ctx.Response().Status(401).Json(apierrors.New("UNAUTHENTICATED", "请先登录", nil))
	}

	items := c.registry.ForPermissions(currentUser.Permissions)
	return ctx.Response().Success().Json(response.Success(items, response.Meta{RequestID: menuRequestID(ctx)}))
}

func menuRequestID(ctx http.Context) string {
	requestID := ctx.Request().Header("X-Request-ID")
	if requestID == "" {
		requestID = "menu-request"
	}
	ctx.Response().Header("X-Request-ID", requestID)
	return requestID
}
