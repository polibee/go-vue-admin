package controllers

import (
	"strconv"

	"github.com/goravel/framework/contracts/http"

	"goravel/app/facades"
	"goravel/app/models"
)

func (r *ResourceController) Show(ctx http.Context) http.Response {
	id, err := strconv.ParseInt(ctx.Request().Route("id"), 10, 64)
	if err != nil || id < 1 {
		return ctx.Response().Status(404).Json(http.Json{"code": "RESOURCE_NOT_FOUND"})
	}

	resourceName := ctx.Request().Route("resource")
	switch resourceName {
	case "users":
		var user models.User
		if err := facades.Orm().Query().Where("id", id).First(&user); err != nil {
			return ctx.Response().Status(404).Json(http.Json{"code": "RESOURCE_NOT_FOUND"})
		}
		return ctx.Response().Success().Json(http.Json{"data": user.Public()})
	case "roles":
		var role models.Role
		if err := facades.Orm().Query().Where("id", id).First(&role); err != nil {
			return ctx.Response().Status(404).Json(http.Json{"code": "RESOURCE_NOT_FOUND"})
		}
		return ctx.Response().Success().Json(http.Json{"data": role})
	case "permissions":
		var permission models.Permission
		if err := facades.Orm().Query().Where("id", id).First(&permission); err != nil {
			return ctx.Response().Status(404).Json(http.Json{"code": "RESOURCE_NOT_FOUND"})
		}
		return ctx.Response().Success().Json(http.Json{"data": permission})
	default:
		return ctx.Response().Status(404).Json(http.Json{"code": "RESOURCE_NOT_FOUND"})
	}
}
