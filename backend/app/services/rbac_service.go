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
		Where("users.is_active = ? AND permissions.name IN (?, ?, ?)", true, "admin.users.view", "admin.roles.manage", "admin.permissions.manage").
		Get(&activeAdmins); err != nil {
		return false, err
	}
	if len(activeAdmins) != 1 {
		return false, nil
	}
	return activeAdmins[0].UserID == userID, nil
}
