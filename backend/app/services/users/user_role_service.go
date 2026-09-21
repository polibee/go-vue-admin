package userservices

import (
	"errors"

	"github.com/goravel/framework/contracts/database/orm"

	"goravel/app/facades"
	"goravel/app/models"
	"goravel/app/rbac"
	rbacservices "goravel/app/services/rbac"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrLastAdmin    = errors.New("last active administrator cannot be removed")
)

type UserRoleService struct{}

func NewUserRoleService() *UserRoleService {
	return &UserRoleService{}
}

func (s *UserRoleService) Roles(userID int64) ([]models.Role, error) {
	var user models.User
	if err := facades.Orm().Query().Where("id = ?", userID).First(&user); err != nil {
		return nil, ErrUserNotFound
	}
	var roles []models.Role
	if err := facades.Orm().Query().
		Table("roles").
		Join("JOIN role_user ON role_user.role_id = roles.id").
		Where("role_user.user_id = ?", userID).
		OrderBy("roles.id").
		Get(&roles); err != nil {
		return nil, err
	}
	return roles, nil
}

func (s *UserRoleService) ReplaceRoles(operatorID, userID int64, roleIDs []int64) error {
	var user models.User
	if err := facades.Orm().Query().Where("id = ?", userID).First(&user); err != nil {
		return ErrUserNotFound
	}
	if err := rbac.ValidateRoleIDs(roleIDs); err != nil {
		return err
	}

	roles := make([]models.Role, 0, len(roleIDs))
	for _, roleID := range roleIDs {
		var role models.Role
		if err := facades.Orm().Query().Where("id = ?", roleID).First(&role); err != nil {
			return rbacservices.ErrRoleNotFound
		}
		roles = append(roles, role)
	}
	adminRole := s.rolesGrantAdministration(roleIDs)
	if userID == operatorID && !adminRole {
		return ErrLastAdmin
	}
	if user.Status == userStatusActive && !adminRole {
		lastAdmin, err := rbacservices.NewRBACService().IsLastActiveAdmin(userID)
		if err != nil {
			return err
		}
		if lastAdmin {
			return ErrLastAdmin
		}
	}

	return facades.Orm().Transaction(func(tx orm.Query) error {
		if _, err := tx.Table("role_user").Where("user_id = ?", userID).Delete(); err != nil {
			return err
		}
		for _, role := range roles {
			if err := tx.Table("role_user").Create(&map[string]any{
				"user_id": userID,
				"role_id": role.ID,
			}); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *UserRoleService) rolesGrantAdministration(roleIDs []int64) bool {
	for _, roleID := range roleIDs {
		allowed, err := facades.Orm().Query().
			Table("permission_role").
			Join("JOIN permissions ON permissions.id = permission_role.permission_id").
			Where("permission_role.role_id = ? AND permissions.name IN (?, ?, ?)", roleID, "admin.users.view", "admin.roles.manage", "admin.permissions.manage").
			Exists()
		if err == nil && allowed {
			return true
		}
	}
	return false
}
