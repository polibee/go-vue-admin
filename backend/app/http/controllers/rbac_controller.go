package controllers

import (
	"errors"
	"strconv"

	"github.com/goravel/framework/contracts/http"

	"goravel/app/facades"
	"goravel/app/models"
	"goravel/app/rbac"
	"goravel/app/services"
)

type RBACController struct{}

func NewRBACController() *RBACController { return &RBACController{} }

type rolePayload struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
}

type rolePermissionsPayload struct {
	PermissionIDs []int64 `json:"permission_ids"`
}

type userRolesPayload struct {
	RoleIDs []int64 `json:"role_ids"`
}

func (r *RBACController) Users(ctx http.Context) http.Response {
	if err := parseAuthToken(ctx); err != nil {
		return unauthorized(ctx)
	}
	var users []models.User
	if err := facades.Orm().Query().Select("id", "name", "email", "is_active", "locale", "created_at", "updated_at").OrderByDesc("id").Get(&users); err != nil {
		return ctx.Response().Status(500).Json(http.Json{"code": "RBAC_USERS_ERROR", "message": "could not load users"})
	}
	items := make([]map[string]any, 0, len(users))
	for _, user := range users {
		items = append(items, user.Public())
	}
	return ctx.Response().Success().Json(http.Json{"data": items})
}

func (r *RBACController) Roles(ctx http.Context) http.Response {
	if err := parseAuthToken(ctx); err != nil {
		return unauthorized(ctx)
	}
	var roles []models.Role
	if err := facades.Orm().Query().OrderBy("id").Get(&roles); err != nil {
		return ctx.Response().Status(500).Json(http.Json{"code": "RBAC_ROLES_ERROR", "message": "could not load roles"})
	}
	return ctx.Response().Success().Json(http.Json{"data": roles})
}

func (r *RBACController) Permissions(ctx http.Context) http.Response {
	if err := parseAuthToken(ctx); err != nil {
		return unauthorized(ctx)
	}
	var permissions []models.Permission
	if err := facades.Orm().Query().OrderBy("id").Get(&permissions); err != nil {
		return ctx.Response().Status(500).Json(http.Json{"code": "RBAC_PERMISSIONS_ERROR", "message": "could not load permissions"})
	}
	return ctx.Response().Success().Json(http.Json{"data": permissions})
}

func (r *RBACController) CreateRole(ctx http.Context) http.Response {
	var payload rolePayload
	if err := ctx.Request().Bind(&payload); err != nil {
		return rbacError(ctx, 422, "VALIDATION_ERROR")
	}
	role, err := services.NewRoleService().Create(payload.Name, payload.DisplayName)
	if err != nil {
		return roleServiceError(ctx, err)
	}
	return ctx.Response().Status(201).Json(http.Json{"data": role})
}

func (r *RBACController) ShowRole(ctx http.Context) http.Response {
	role, err := services.NewRoleService().Find(ctx.Request().RouteInt64("id"))
	if err != nil {
		return roleServiceError(ctx, err)
	}
	permissions, err := services.NewRoleService().Permissions(int64(role.ID))
	if err != nil {
		return rbacError(ctx, 500, "INTERNAL_ERROR")
	}
	return ctx.Response().Success().Json(http.Json{"data": http.Json{
		"id": role.ID, "name": role.Name, "display_name": role.DisplayName, "permissions": permissions,
	}})
}

func (r *RBACController) UpdateRole(ctx http.Context) http.Response {
	var payload rolePayload
	if err := ctx.Request().Bind(&payload); err != nil {
		return rbacError(ctx, 422, "VALIDATION_ERROR")
	}
	role, err := services.NewRoleService().Update(ctx.Request().RouteInt64("id"), payload.Name, payload.DisplayName)
	if err != nil {
		return roleServiceError(ctx, err)
	}
	return ctx.Response().Success().Json(http.Json{"data": role})
}

func (r *RBACController) DeleteRole(ctx http.Context) http.Response {
	if err := services.NewRoleService().Delete(ctx.Request().RouteInt64("id")); err != nil {
		return roleServiceError(ctx, err)
	}
	return ctx.Response().NoContent(204)
}

func (r *RBACController) ReplaceRolePermissions(ctx http.Context) http.Response {
	var payload rolePermissionsPayload
	if err := ctx.Request().Bind(&payload); err != nil {
		return rbacError(ctx, 422, "VALIDATION_ERROR")
	}
	if err := services.NewRoleService().ReplacePermissions(ctx.Request().RouteInt64("id"), payload.PermissionIDs); err != nil {
		return roleServiceError(ctx, err)
	}
	return ctx.Response().NoContent(204)
}

func (r *RBACController) UserRoles(ctx http.Context) http.Response {
	roles, err := services.NewUserRoleService().Roles(ctx.Request().RouteInt64("id"))
	if err != nil {
		return userRoleServiceError(ctx, err)
	}
	return ctx.Response().Success().Json(http.Json{"data": roles})
}

func (r *RBACController) ReplaceUserRoles(ctx http.Context) http.Response {
	var payload userRolesPayload
	if err := ctx.Request().Bind(&payload); err != nil {
		return rbacError(ctx, 422, "VALIDATION_ERROR")
	}
	identity, err := facades.Auth(ctx).ID()
	if err != nil {
		return unauthorized(ctx)
	}
	operatorID, err := strconv.ParseInt(identity, 10, 64)
	if err != nil {
		return unauthorized(ctx)
	}
	if err := services.NewUserRoleService().ReplaceRoles(operatorID, ctx.Request().RouteInt64("id"), payload.RoleIDs); err != nil {
		return userRoleServiceError(ctx, err)
	}
	return ctx.Response().NoContent(204)
}

func roleServiceError(ctx http.Context, err error) http.Response {
	switch {
	case errors.Is(err, services.ErrRoleNotFound):
		return rbacError(ctx, 404, "RBAC_ROLE_NOT_FOUND")
	case errors.Is(err, services.ErrPermissionNotFound):
		return rbacError(ctx, 404, "RBAC_PERMISSION_NOT_FOUND")
	case errors.Is(err, services.ErrSystemRole):
		return rbacError(ctx, 409, "RBAC_SYSTEM_ROLE")
	case errors.Is(err, rbac.ErrInvalidRoleInput), errors.Is(err, services.ErrDuplicateRole):
		return rbacError(ctx, 422, "VALIDATION_ERROR")
	default:
		return rbacError(ctx, 500, "INTERNAL_ERROR")
	}
}

func userRoleServiceError(ctx http.Context, err error) http.Response {
	switch {
	case errors.Is(err, services.ErrUserNotFound):
		return rbacError(ctx, 404, "RBAC_USER_NOT_FOUND")
	case errors.Is(err, services.ErrRoleNotFound):
		return rbacError(ctx, 404, "RBAC_ROLE_NOT_FOUND")
	case errors.Is(err, services.ErrLastAdmin):
		return rbacError(ctx, 409, "RBAC_LAST_ADMIN")
	case errors.Is(err, rbac.ErrInvalidRoleIDs):
		return rbacError(ctx, 422, "VALIDATION_ERROR")
	default:
		return rbacError(ctx, 500, "INTERNAL_ERROR")
	}
}

func rbacError(ctx http.Context, status int, code string) http.Response {
	return ctx.Response().Status(status).Json(http.Json{"code": code})
}

func parseAuthToken(ctx http.Context) error {
	token := ctx.Request().Header("Authorization")
	if token == "" {
		return errMissingToken
	}
	_, err := facades.Auth(ctx).Parse(token)
	return err
}
