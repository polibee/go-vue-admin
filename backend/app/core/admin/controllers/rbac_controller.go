package controllers

import (
	"errors"
	"strconv"

	"github.com/goravel/framework/contracts/http"

	"goravel/app/facades"
	"goravel/app/models"
	"goravel/app/rbac"
	auditservices "goravel/app/services/audit"
	rbacservices "goravel/app/services/rbac"
	userservices "goravel/app/services/users"
)

var errMissingToken = errors.New("missing authorization token")

func unauthorized(ctx http.Context) http.Response {
	return ctx.Response().Status(401).Json(http.Json{
		"code":    "AUTH_UNAUTHORIZED",
		"message": "authentication required",
	})
}

func recordAudit(userID uint, action string, metadata map[string]any) {
	if err := auditservices.NewAuditService().Record(userID, action, metadata); err != nil {
		facades.Log().Errorf("audit record failed action=%s user_id=%d error=%v", action, userID, err)
	}
}

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

type userPayload struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Locale   string `json:"locale"`
	Status   string `json:"status"`
}

type userBulkStatusPayload struct {
	UserIDs []int64 `json:"user_ids"`
	Status  string  `json:"status"`
}

func (r *RBACController) Users(ctx http.Context) http.Response {
	if err := parseAuthToken(ctx); err != nil {
		return unauthorized(ctx)
	}
	var users []models.User
	if err := facades.Orm().Query().Select("id", "name", "email", "status", "locale", "created_at", "updated_at").OrderByDesc("id").Get(&users); err != nil {
		return ctx.Response().Status(500).Json(http.Json{"code": "RBAC_USERS_ERROR", "message": "could not load users"})
	}
	items := make([]map[string]any, 0, len(users))
	for _, user := range users {
		items = append(items, user.Public())
	}
	return ctx.Response().Success().Json(http.Json{"data": items})
}

func (r *RBACController) CreateUser(ctx http.Context) http.Response {
	var payload userPayload
	if err := ctx.Request().Bind(&payload); err != nil {
		return rbacError(ctx, 422, "VALIDATION_ERROR")
	}
	user, err := userservices.NewUserService().Create(payload.Name, payload.Email, payload.Password, payload.Locale, payload.Status)
	if err != nil {
		return userServiceError(ctx, err)
	}
	recordManagementAudit(ctx, "user.create", map[string]any{"target_user_id": user.ID})
	return ctx.Response().Status(201).Json(http.Json{"data": user.Public()})
}

func (r *RBACController) UpdateUser(ctx http.Context) http.Response {
	var payload userPayload
	if err := ctx.Request().Bind(&payload); err != nil {
		return rbacError(ctx, 422, "VALIDATION_ERROR")
	}
	user, err := userservices.NewUserService().Update(ctx.Request().RouteInt64("id"), payload.Name, payload.Email, payload.Password, payload.Locale, payload.Status)
	if err != nil {
		return userServiceError(ctx, err)
	}
	recordManagementAudit(ctx, "user.update", map[string]any{"target_user_id": user.ID})
	return ctx.Response().Success().Json(http.Json{"data": user.Public()})
}

func (r *RBACController) DeleteUser(ctx http.Context) http.Response {
	targetID := ctx.Request().RouteInt64("id")
	if err := userservices.NewUserService().Delete(targetID); err != nil {
		return userServiceError(ctx, err)
	}
	recordManagementAudit(ctx, "user.delete", map[string]any{"target_user_id": targetID})
	return ctx.Response().NoContent(204)
}

func (r *RBACController) BulkSetUserStatus(ctx http.Context) http.Response {
	var payload userBulkStatusPayload
	if err := ctx.Request().Bind(&payload); err != nil {
		return rbacError(ctx, 422, "VALIDATION_ERROR")
	}
	if err := userservices.NewUserService().BulkSetStatus(payload.UserIDs, payload.Status); err != nil {
		return userServiceError(ctx, err)
	}
	recordManagementAudit(ctx, "user.status.bulk", map[string]any{"target_user_ids": payload.UserIDs, "status": payload.Status})
	return ctx.Response().NoContent(204)
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
	role, err := rbacservices.NewRoleService().Create(payload.Name, payload.DisplayName)
	if err != nil {
		return roleServiceError(ctx, err)
	}
	recordManagementAudit(ctx, "role.create", map[string]any{"target_role_id": role.ID})
	return ctx.Response().Status(201).Json(http.Json{"data": role})
}

func (r *RBACController) ShowRole(ctx http.Context) http.Response {
	role, err := rbacservices.NewRoleService().Find(ctx.Request().RouteInt64("id"))
	if err != nil {
		return roleServiceError(ctx, err)
	}
	permissions, err := rbacservices.NewRoleService().Permissions(int64(role.ID))
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
	role, err := rbacservices.NewRoleService().Update(ctx.Request().RouteInt64("id"), payload.Name, payload.DisplayName)
	if err != nil {
		return roleServiceError(ctx, err)
	}
	recordManagementAudit(ctx, "role.update", map[string]any{"target_role_id": role.ID})
	return ctx.Response().Success().Json(http.Json{"data": role})
}

func (r *RBACController) DeleteRole(ctx http.Context) http.Response {
	targetID := ctx.Request().RouteInt64("id")
	if err := rbacservices.NewRoleService().Delete(targetID); err != nil {
		return roleServiceError(ctx, err)
	}
	recordManagementAudit(ctx, "role.delete", map[string]any{"target_role_id": targetID})
	return ctx.Response().NoContent(204)
}

func (r *RBACController) ReplaceRolePermissions(ctx http.Context) http.Response {
	var payload rolePermissionsPayload
	if err := ctx.Request().Bind(&payload); err != nil {
		return rbacError(ctx, 422, "VALIDATION_ERROR")
	}
	targetID := ctx.Request().RouteInt64("id")
	if err := rbacservices.NewRoleService().ReplacePermissions(targetID, payload.PermissionIDs); err != nil {
		return roleServiceError(ctx, err)
	}
	recordManagementAudit(ctx, "role.permissions.replace", map[string]any{"target_role_id": targetID, "permission_ids": payload.PermissionIDs})
	return ctx.Response().NoContent(204)
}

func (r *RBACController) UserRoles(ctx http.Context) http.Response {
	roles, err := userservices.NewUserRoleService().Roles(ctx.Request().RouteInt64("id"))
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
	targetID := ctx.Request().RouteInt64("id")
	if err := userservices.NewUserRoleService().ReplaceRoles(operatorID, targetID, payload.RoleIDs); err != nil {
		return userRoleServiceError(ctx, err)
	}
	recordManagementAudit(ctx, "user.roles.replace", map[string]any{"target_user_id": targetID, "role_ids": payload.RoleIDs})
	return ctx.Response().NoContent(204)
}

func recordManagementAudit(ctx http.Context, action string, metadata map[string]any) {
	identity, err := facades.Auth(ctx).ID()
	if err != nil {
		return
	}
	operatorID, err := strconv.ParseUint(identity, 10, 64)
	if err != nil {
		return
	}
	recordAudit(uint(operatorID), action, metadata)
}

func roleServiceError(ctx http.Context, err error) http.Response {
	switch {
	case errors.Is(err, rbacservices.ErrRoleNotFound):
		return rbacError(ctx, 404, "RBAC_ROLE_NOT_FOUND")
	case errors.Is(err, rbacservices.ErrPermissionNotFound):
		return rbacError(ctx, 404, "RBAC_PERMISSION_NOT_FOUND")
	case errors.Is(err, rbacservices.ErrSystemRole):
		return rbacError(ctx, 409, "RBAC_SYSTEM_ROLE")
	case errors.Is(err, rbac.ErrInvalidRoleInput), errors.Is(err, rbacservices.ErrDuplicateRole):
		return rbacError(ctx, 422, "VALIDATION_ERROR")
	default:
		return rbacError(ctx, 500, "INTERNAL_ERROR")
	}
}

func userRoleServiceError(ctx http.Context, err error) http.Response {
	switch {
	case errors.Is(err, userservices.ErrUserNotFound):
		return rbacError(ctx, 404, "RBAC_USER_NOT_FOUND")
	case errors.Is(err, rbacservices.ErrRoleNotFound):
		return rbacError(ctx, 404, "RBAC_ROLE_NOT_FOUND")
	case errors.Is(err, userservices.ErrLastAdmin):
		return rbacError(ctx, 409, "RBAC_LAST_ADMIN")
	case errors.Is(err, rbac.ErrInvalidRoleIDs):
		return rbacError(ctx, 422, "VALIDATION_ERROR")
	default:
		return rbacError(ctx, 500, "INTERNAL_ERROR")
	}
}

func userServiceError(ctx http.Context, err error) http.Response {
	switch {
	case errors.Is(err, userservices.ErrUserNotFound):
		return rbacError(ctx, 404, "RBAC_USER_NOT_FOUND")
	case errors.Is(err, userservices.ErrLastAdmin):
		return rbacError(ctx, 409, "RBAC_LAST_ADMIN")
	case errors.Is(err, userservices.ErrUserExists), errors.Is(err, userservices.ErrInvalidUser):
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
