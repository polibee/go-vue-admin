package rbacservices

import (
	"errors"
	"strconv"

	"github.com/goravel/framework/contracts/database/orm"
	"github.com/goravel/framework/contracts/http"

	"goravel/app/core/resource"
	"goravel/app/facades"
)

var (
	ErrInvalidOwnerScope = errors.New("resource own scope requires an owner field")
	ErrInvalidScopeUser  = errors.New("invalid authenticated user identity")
)

type ResourceScopeService struct{}

func NewResourceScopeService() *ResourceScopeService {
	return &ResourceScopeService{}
}

func (s *ResourceScopeService) Apply(ctx http.Context, query orm.Query, manifest resource.Manifest, action string) (orm.Query, error) {
	scope, userID, err := s.Resolve(ctx, manifest, action)
	if err != nil {
		return query, err
	}
	field, value, err := OwnerPredicate(manifest, scope, userID)
	if err != nil {
		return query, err
	}
	if field == "" {
		return query, nil
	}
	return query.Where(field+" = ?", value), nil
}

func (s *ResourceScopeService) CanAccess(ctx http.Context, manifest resource.Manifest, action string, id int64) (bool, error) {
	scope, userID, err := s.Resolve(ctx, manifest, action)
	if err != nil {
		return false, err
	}
	field, value, err := OwnerPredicate(manifest, scope, userID)
	if err != nil {
		return false, err
	}
	if field == "" {
		return true, nil
	}
	return facades.Orm().Query().Table(manifest.Table).Where("id = ? AND "+field+" = ?", id, value).Exists()
}

func (s *ResourceScopeService) Resolve(ctx http.Context, manifest resource.Manifest, action string) (resource.DataScope, int64, error) {
	identity, err := facades.Auth(ctx).ID()
	if err != nil {
		return "", 0, ErrInvalidScopeUser
	}
	userID, err := strconv.ParseInt(identity, 10, 64)
	if err != nil || userID < 1 {
		return "", 0, ErrInvalidScopeUser
	}
	permission := manifestPermission(manifest, action)
	if permission == "" {
		return "", 0, ErrDataScopeNotAssigned
	}
	assigned, err := NewDataScopeService().EffectiveScope(userID, permission)
	if err != nil {
		return "", 0, err
	}
	declared := manifest.DataScope
	if declared == "" {
		declared = resource.DataScopeAll
	}
	if declared == resource.DataScopeOwn || assigned == resource.DataScopeOwn {
		return resource.DataScopeOwn, userID, nil
	}
	return resource.DataScopeAll, userID, nil
}

func OwnerPredicate(manifest resource.Manifest, scope resource.DataScope, userID int64) (string, any, error) {
	if scope != resource.DataScopeOwn {
		return "", nil, nil
	}
	if manifest.OwnerField == "" {
		return "", nil, ErrInvalidOwnerScope
	}
	return manifest.OwnerField, userID, nil
}

func manifestPermission(manifest resource.Manifest, action string) string {
	if action == "view" && len(manifest.Permissions) > 0 {
		return manifest.Permissions[0]
	}
	for _, candidate := range manifest.Actions {
		if candidate.Name == action {
			return candidate.Permission
		}
	}
	return ""
}
