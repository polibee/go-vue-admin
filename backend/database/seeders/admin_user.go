package seeders

import (
	"goravel/app/facades"
	"goravel/app/models"
)

type AdminUser struct {
}

// Signature The name and signature of the seeder.
func (s *AdminUser) Signature() string {
	return "AdminUser"
}

// Run executes the seeder logic.
func (s *AdminUser) Run() error {
	const email = "admin@example.com"

	exists, err := facades.Orm().Query().Model(&models.User{}).Where("email = ?", email).Exists()
	if err != nil {
		return err
	}
	if !exists {
		password, err := facades.Hash().Make("Admin123!")
		if err != nil {
			return err
		}
		if err := facades.Orm().Query().Create(&models.User{
			Name:     "Administrator",
			Email:    email,
			Password: password,
			IsActive: true,
			Locale:   "zh-CN",
		}); err != nil {
			return err
		}
	}

	var roles []models.Role
	if err := facades.Orm().Query().Where("name = ?", "super-admin").Get(&roles); err != nil {
		return err
	}
	if len(roles) == 0 {
		if err := facades.Orm().Query().Create(&models.Role{Name: "super-admin", DisplayName: "Super Administrator"}); err != nil {
			return err
		}
		if err := facades.Orm().Query().Where("name = ?", "super-admin").Get(&roles); err != nil {
			return err
		}
	}

	var users []models.User
	if err := facades.Orm().Query().Where("email = ?", email).Get(&users); err != nil {
		return err
	}
	if len(users) == 0 || len(roles) == 0 {
		return nil
	}
	assigned, err := facades.Orm().Query().Table("role_user").Where("user_id = ? AND role_id = ?", users[0].ID, roles[0].ID).Exists()
	if err != nil {
		return err
	}
	if !assigned {
		if err := facades.Orm().Query().Table("role_user").Create(&map[string]any{"user_id": users[0].ID, "role_id": roles[0].ID}); err != nil {
			return err
		}
	}

	for _, permission := range []models.Permission{
		{Name: "admin.users.view", DisplayName: "View users"},
		{Name: "admin.users.manage", DisplayName: "Manage users"},
		{Name: "admin.roles.manage", DisplayName: "Manage roles"},
		{Name: "admin.permissions.manage", DisplayName: "Manage permissions"},
	} {
		var permissionExists []models.Permission
		if err := facades.Orm().Query().Where("name = ?", permission.Name).Get(&permissionExists); err != nil {
			return err
		}
		if len(permissionExists) == 0 {
			if err := facades.Orm().Query().Create(&permission); err != nil {
				return err
			}
			permissionExists = []models.Permission{permission}
		}
		assignedPermission, err := facades.Orm().Query().Table("permission_role").Where("permission_id = ? AND role_id = ?", permissionExists[0].ID, roles[0].ID).Exists()
		if err != nil {
			return err
		}
		if !assignedPermission {
			if err := facades.Orm().Query().Table("permission_role").Create(&map[string]any{"permission_id": permissionExists[0].ID, "role_id": roles[0].ID}); err != nil {
				return err
			}
		}
	}
	return nil
}
