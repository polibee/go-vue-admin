package rbacservices

import (
	"errors"
	"fmt"
	"strings"

	"github.com/goravel/framework/contracts/database/orm"

	"goravel/app/core/resource"
	"goravel/app/facades"
	"goravel/app/models"
	"goravel/app/modules/admin/registry"
	"goravel/app/rbac"
	notificationservices "goravel/app/services/notifications"
)

var (
	ErrRoleNotFound       = errors.New("role not found")
	ErrPermissionNotFound = errors.New("permission not found")
	ErrSystemRole         = errors.New("system role cannot be changed")
	ErrDuplicateRole      = errors.New("role name already exists")
	ErrInvalidDataScope   = errors.New("invalid data scope")
)

const SystemRoleName = "super-admin"

type RoleService struct{}

type PermissionAssignment struct {
	models.Permission
	Scope  resource.DataScope       `json:"scope"`
	Fields map[string]FieldOverride `json:"fields,omitempty"`
}

func NewRoleService() *RoleService {
	return &RoleService{}
}

func validateRoleInput(name, displayName string) error {
	return rbac.ValidateRoleInput(name, displayName)
}

func (s *RoleService) Create(name, displayName string) (*models.Role, error) {
	if err := validateRoleInput(name, displayName); err != nil {
		return nil, err
	}
	name = strings.TrimSpace(name)
	displayName = strings.TrimSpace(displayName)
	exists, err := facades.Orm().Query().Table("roles").Where("name = ?", name).Exists()
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrDuplicateRole
	}

	role := &models.Role{Name: name, DisplayName: displayName}
	if err := facades.Orm().Query().Create(role); err != nil {
		return nil, err
	}
	return role, nil
}

func (s *RoleService) Find(id int64) (*models.Role, error) {
	role := &models.Role{}
	if err := facades.Orm().Query().Where("id = ?", id).First(role); err != nil {
		return nil, ErrRoleNotFound
	}
	return role, nil
}

func (s *RoleService) Permissions(roleID int64) ([]models.Permission, error) {
	if _, err := s.Find(roleID); err != nil {
		return nil, err
	}
	var permissions []models.Permission
	if err := facades.Orm().Query().
		Table("permissions").
		Join("JOIN permission_role ON permission_role.permission_id = permissions.id").
		Where("permission_role.role_id = ?", roleID).
		OrderBy("permissions.id").
		Get(&permissions); err != nil {
		return nil, err
	}
	return permissions, nil
}

func (s *RoleService) PermissionAssignments(roleID int64) ([]PermissionAssignment, error) {
	if _, err := s.Find(roleID); err != nil {
		return nil, err
	}
	var assignments []PermissionAssignment
	if err := facades.Orm().Query().
		Table("permissions").
		Select("permissions.id, permissions.name, permissions.display_name, permission_role.scope").
		Join("JOIN permission_role ON permission_role.permission_id = permissions.id").
		Where("permission_role.role_id = ?", roleID).
		OrderBy("permissions.id").
		Get(&assignments); err != nil {
		return nil, err
	}
	for index := range assignments {
		if assignments[index].Scope == "" {
			assignments[index].Scope = resource.DataScopeAll
		}
		fields, fieldErr := s.permissionFieldOverrides(roleID, int64(assignments[index].ID))
		if fieldErr != nil {
			return nil, fieldErr
		}
		assignments[index].Fields = fields
	}
	return assignments, nil
}

func (s *RoleService) Update(id int64, name, displayName string) (*models.Role, error) {
	if err := validateRoleInput(name, displayName); err != nil {
		return nil, err
	}
	role, err := s.Find(id)
	if err != nil {
		return nil, err
	}
	if role.Name == SystemRoleName && strings.TrimSpace(name) != SystemRoleName {
		return nil, ErrSystemRole
	}
	name = strings.TrimSpace(name)
	duplicate, err := facades.Orm().Query().Table("roles").Where("name = ? AND id <> ?", name, id).Exists()
	if err != nil {
		return nil, err
	}
	if duplicate {
		return nil, ErrDuplicateRole
	}
	role.Name = name
	role.DisplayName = strings.TrimSpace(displayName)
	if err := facades.Orm().Query().Save(role); err != nil {
		return nil, err
	}
	return role, nil
}

func (s *RoleService) Delete(id int64) error {
	role, err := s.Find(id)
	if err != nil {
		return err
	}
	if role.Name == SystemRoleName {
		return ErrSystemRole
	}
	return facades.Orm().Transaction(func(tx orm.Query) error {
		if _, err := tx.Table("permission_role_field").Where("role_id = ?", id).Delete(); err != nil {
			return err
		}
		if _, err := tx.Table("permission_role").Where("role_id = ?", id).Delete(); err != nil {
			return err
		}
		if _, err := tx.Table("role_user").Where("role_id = ?", id).Delete(); err != nil {
			return err
		}
		_, err := tx.Table("roles").Where("id = ?", id).Delete()
		return err
	})
}

func (s *RoleService) ReplacePermissions(roleID int64, permissionIDs []int64, scopes map[int64]resource.DataScope, fields map[int64]map[string]FieldOverride) error {
	role, err := s.Find(roleID)
	if err != nil {
		return err
	}
	if role.Name == SystemRoleName {
		return ErrSystemRole
	}

	seen := make(map[int64]struct{}, len(permissionIDs))
	for _, permissionID := range permissionIDs {
		if permissionID <= 0 {
			return ErrPermissionNotFound
		}
		if _, ok := seen[permissionID]; ok {
			continue
		}
		scope := scopes[permissionID]
		if scope == "" {
			scope = resource.DataScopeAll
		}
		if err := validatePermissionScope(scope); err != nil {
			return err
		}
		scopes[permissionID] = scope
		seen[permissionID] = struct{}{}
		permission := &models.Permission{}
		if err := facades.Orm().Query().Where("id = ?", permissionID).First(permission); err != nil {
			return ErrPermissionNotFound
		}
		if overrides := fields[permissionID]; len(overrides) > 0 {
			manifest, ok := manifestForPermission(permission.Name)
			if !ok {
				return ErrPermissionNotFound
			}
			if err := NewFieldPermissionService().ValidateFieldOverrides(manifest, overrides); err != nil {
				return err
			}
		}
	}
	for permissionID := range fields {
		if _, ok := seen[permissionID]; !ok {
			return ErrPermissionNotFound
		}
	}

	if err := facades.Orm().Transaction(func(tx orm.Query) error {
		if _, err := tx.Table("permission_role_field").Where("role_id = ?", roleID).Delete(); err != nil {
			return err
		}
		if _, err := tx.Table("permission_role").Where("role_id = ?", roleID).Delete(); err != nil {
			return err
		}
		for permissionID := range seen {
			if err := tx.Table("permission_role").Create(&map[string]any{
				"permission_id": permissionID,
				"role_id":       roleID,
				"scope":         scopes[permissionID],
			}); err != nil {
				return err
			}
			for fieldName, override := range fields[permissionID] {
				if err := tx.Table("permission_role_field").Create(&map[string]any{
					"role_id": roleID, "permission_id": permissionID, "field_name": fieldName,
					"readable": override.Readable, "writable": override.Writable,
				}); err != nil {
					return err
				}
			}
		}
		return nil
	}); err != nil {
		return err
	}
	var members []struct {
		UserID uint `db:"user_id"`
	}
	if err := facades.Orm().Query().Table("role_user").Select("user_id").Where("role_id = ?", roleID).Get(&members); err != nil {
		return err
	}
	for _, member := range members {
		notificationservices.NewNotificationService().PublishBestEffort(member.UserID, notificationservices.NotificationInput{
			Type:  notificationservices.TypeRolePermissionsChanged,
			Title: "角色权限已变更",
			Body:  fmt.Sprintf("你所属的角色“%s”的权限已更新。", role.DisplayName),
			URL:   fmt.Sprintf("/admin/roles/%d/edit", role.ID),
		})
	}
	return nil
}

func (s *RoleService) permissionFieldOverrides(roleID, permissionID int64) (map[string]FieldOverride, error) {
	var rows []struct {
		FieldName string `db:"field_name"`
		Readable  bool   `db:"readable"`
		Writable  bool   `db:"writable"`
	}
	if err := facades.Orm().Query().Table("permission_role_field").Where("role_id = ? AND permission_id = ?", roleID, permissionID).Get(&rows); err != nil {
		return nil, err
	}
	fields := make(map[string]FieldOverride, len(rows))
	for _, row := range rows {
		fields[row.FieldName] = FieldOverride{Readable: row.Readable, Writable: row.Writable}
	}
	return fields, nil
}

func manifestForPermission(permission string) (resource.Manifest, bool) {
	parts := strings.Split(permission, ".")
	if len(parts) < 3 || parts[0] != "admin" {
		return resource.Manifest{}, false
	}
	manifest, err := registry.AdminRegistry().Find(parts[1])
	return manifest, err == nil
}

func validatePermissionScope(scope resource.DataScope) error {
	if scope != resource.DataScopeAll && scope != resource.DataScopeOwn {
		return ErrInvalidDataScope
	}
	return nil
}
