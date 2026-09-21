package services

import (
	"errors"
	"strconv"

	"github.com/goravel/framework/contracts/http"

	"goravel/app/facades"
)

var ErrInvalidIdentity = errors.New("invalid authenticated user identity")

type RBACService struct{}

func NewRBACService() *RBACService {
	return &RBACService{}
}

func (s *RBACService) UserHasPermission(ctx http.Context, permission string) (bool, error) {
	if ctx.Request().Header("Authorization") == "" {
		return false, ErrInvalidIdentity
	}

	identity, err := facades.Auth(ctx).ID()
	if err != nil {
		return false, err
	}
	userID, err := strconv.ParseInt(identity, 10, 64)
	if err != nil {
		return false, ErrInvalidIdentity
	}

	return facades.Orm().Query().
		Table("role_user").
		Join("JOIN roles ON roles.id = role_user.role_id").
		Join("JOIN permission_role ON permission_role.role_id = roles.id").
		Join("JOIN permissions ON permissions.id = permission_role.permission_id").
		Where("role_user.user_id = ? AND permissions.name = ?", userID, permission).
		Exists()
}

func (s *RBACService) PermissionsForUser(userID uint) ([]string, error) {
	var rows []struct {
		Name string `db:"name"`
	}
	if err := facades.Orm().Query().
		Table("permissions").
		Select("permissions.name").
		Join("JOIN permission_role ON permission_role.permission_id = permissions.id").
		Join("JOIN role_user ON role_user.role_id = permission_role.role_id").
		Where("role_user.user_id = ?", userID).
		OrderBy("permissions.name").
		Get(&rows); err != nil {
		return nil, err
	}
	permissions := make([]string, 0, len(rows))
	seen := make(map[string]struct{}, len(rows))
	for _, row := range rows {
		if _, ok := seen[row.Name]; ok {
			continue
		}
		seen[row.Name] = struct{}{}
		permissions = append(permissions, row.Name)
	}
	return permissions, nil
}

func (s *RBACService) IsLastActiveAdmin(userID int64) (bool, error) {
	var activeAdmins []struct {
		UserID int64 `db:"user_id"`
	}
	if err := facades.Orm().Query().
		Table("role_user").
		Select("DISTINCT role_user.user_id").
		Join("JOIN users ON users.id = role_user.user_id").
		Join("JOIN permission_role ON permission_role.role_id = role_user.role_id").
		Join("JOIN permissions ON permissions.id = permission_role.permission_id").
		Where("users.status = ? AND permissions.name IN (?, ?, ?)", userStatusActive, "admin.users.view", "admin.roles.manage", "admin.permissions.manage").
		Get(&activeAdmins); err != nil {
		return false, err
	}
	if len(activeAdmins) != 1 {
		return false, nil
	}
	return activeAdmins[0].UserID == userID, nil
}
