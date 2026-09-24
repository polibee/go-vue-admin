package userservices

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/goravel/framework/contracts/database/orm"

	"goravel/app/facades"
	"goravel/app/models"
	"goravel/app/rbac"
	notificationservices "goravel/app/services/notifications"
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
	previousRoles, err := s.Roles(userID)
	if err != nil {
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

	if err := facades.Orm().Transaction(func(tx orm.Query) error {
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
	}); err != nil {
		return err
	}
	if roleSetChanged(previousRoles, roles) {
		roleNames := make([]string, 0, len(roles))
		for _, role := range roles {
			roleNames = append(roleNames, role.DisplayName)
		}
		sort.Strings(roleNames)
		notificationservices.NewNotificationService().PublishBestEffort(uint(user.ID), notificationservices.NotificationInput{
			Type:  notificationservices.TypeUserRolesChanged,
			Title: "角色分配已变更",
			Body:  fmt.Sprintf("你的角色已更新为：%s。", strings.Join(roleNames, "、")),
			URL:   fmt.Sprintf("/admin/users/%d/edit", user.ID),
		})
	}
	return nil
}

func roleSetChanged(previous []models.Role, current []models.Role) bool {
	previousIDs := make(map[uint]struct{}, len(previous))
	currentIDs := make(map[uint]struct{}, len(current))
	for _, role := range previous {
		previousIDs[role.ID] = struct{}{}
	}
	for _, role := range current {
		currentIDs[role.ID] = struct{}{}
	}
	if len(previousIDs) != len(currentIDs) {
		return true
	}
	for id := range previousIDs {
		if _, ok := currentIDs[id]; !ok {
			return true
		}
	}
	return false
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
