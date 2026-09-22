package rbacservices

import (
	"errors"

	"goravel/app/core/resource"
	"goravel/app/facades"
)

var ErrDataScopeNotAssigned = errors.New("resource data scope is not assigned")

type DataScopeService struct{}

func NewDataScopeService() *DataScopeService {
	return &DataScopeService{}
}

func (s *DataScopeService) EffectiveScope(userID int64, permission string) (resource.DataScope, error) {
	var rows []struct {
		Scope string `db:"scope"`
	}
	if err := facades.Orm().Query().
		Table("permission_role").
		Select("permission_role.scope").
		Join("JOIN role_user ON role_user.role_id = permission_role.role_id").
		Join("JOIN permissions ON permissions.id = permission_role.permission_id").
		Where("role_user.user_id = ? AND permissions.name = ?", userID, permission).
		Get(&rows); err != nil {
		return "", err
	}

	scopes := make([]resource.DataScope, 0, len(rows))
	for _, row := range rows {
		scopes = append(scopes, resource.DataScope(row.Scope))
	}
	scope := mergeDataScopes(scopes)
	if scope == "" {
		return "", ErrDataScopeNotAssigned
	}
	return scope, nil
}

func mergeDataScopes(scopes []resource.DataScope) resource.DataScope {
	for _, scope := range scopes {
		if scope == resource.DataScopeAll {
			return resource.DataScopeAll
		}
	}
	for _, scope := range scopes {
		if scope == resource.DataScopeOwn {
			return resource.DataScopeOwn
		}
	}
	return ""
}
