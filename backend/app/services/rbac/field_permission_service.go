package rbacservices

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/goravel/framework/contracts/http"

	"goravel/app/core/resource"
	"goravel/app/facades"
)

var (
	ErrFieldNotFound         = errors.New("resource field is not declared")
	ErrFieldQueryDenied      = errors.New("resource field query is not allowed")
	ErrFieldPermissionDenied = errors.New("resource field permission denied")
	ErrFieldPolicyExpansion  = errors.New("field policy cannot expand manifest access")
)

type FieldPolicy struct {
	Visible   bool
	Readable  bool
	Writable  bool
	Sensitive bool
}

type FieldOverride struct {
	Readable bool `json:"readable"`
	Writable bool `json:"writable"`
}

type FieldPermissionService struct{}

func NewFieldPermissionService() *FieldPermissionService {
	return &FieldPermissionService{}
}

func (s *FieldPermissionService) ManifestPolicies(manifest resource.Manifest) map[string]FieldPolicy {
	manifest = resource.NormalizeManifestFields(manifest)
	policies := make(map[string]FieldPolicy, len(manifest.Fields))
	for _, field := range manifest.Fields {
		policies[field.Name] = FieldPolicy{
			Visible: field.Visible, Readable: field.Readable,
			Writable: field.Writable, Sensitive: field.Sensitive,
		}
	}
	return policies
}

func (s *FieldPermissionService) EffectivePolicies(ctx http.Context, manifest resource.Manifest, action string) (map[string]FieldPolicy, error) {
	identity, err := facades.Auth(ctx).ID()
	if err != nil {
		return nil, ErrInvalidScopeUser
	}
	userID, err := strconv.ParseInt(identity, 10, 64)
	if err != nil || userID < 1 {
		return nil, ErrInvalidScopeUser
	}
	permission := manifestPermission(manifest, action)
	if permission == "" {
		return nil, ErrDataScopeNotAssigned
	}
	base := s.ManifestPolicies(manifest)
	var assignments []struct {
		RoleID int64 `db:"role_id"`
	}
	if err := facades.Orm().Query().Table("permission_role").Select("permission_role.role_id").Join("JOIN role_user ON role_user.role_id = permission_role.role_id").Join("JOIN permissions ON permissions.id = permission_role.permission_id").Where("role_user.user_id = ? AND permissions.name = ?", userID, permission).Get(&assignments); err != nil {
		return nil, err
	}
	if len(assignments) == 0 {
		return base, nil
	}
	rolePolicies := make([]map[string]FieldPolicy, 0, len(assignments))
	for _, assignment := range assignments {
		rolePolicy := make(map[string]FieldPolicy, len(base))
		for name, policy := range base {
			rolePolicy[name] = policy
		}
		var overrides []struct {
			FieldName string `db:"field_name"`
			Readable  bool   `db:"readable"`
			Writable  bool   `db:"writable"`
		}
		if err := facades.Orm().Query().Table("permission_role_field").Where("role_id = ? AND permission_id = (SELECT id FROM permissions WHERE name = ?)", assignment.RoleID, permission).Get(&overrides); err != nil {
			return nil, err
		}
		for _, override := range overrides {
			policy, ok := rolePolicy[override.FieldName]
			if !ok {
				continue
			}
			policy.Readable = override.Readable
			policy.Writable = override.Writable
			rolePolicy[override.FieldName] = policy
		}
		rolePolicies = append(rolePolicies, rolePolicy)
	}
	return mergeFieldPolicies(base, rolePolicies), nil
}

func mergeFieldPolicies(base map[string]FieldPolicy, roles []map[string]FieldPolicy) map[string]FieldPolicy {
	effective := make(map[string]FieldPolicy, len(base))
	for name, policy := range base {
		effective[name] = FieldPolicy{Sensitive: policy.Sensitive}
	}
	for _, role := range roles {
		for name, policy := range role {
			if _, exists := base[name]; !exists {
				continue
			}
			merged := effective[name]
			merged.Visible = base[name].Visible && (merged.Visible || policy.Visible)
			merged.Readable = base[name].Readable && (merged.Readable || policy.Readable)
			merged.Writable = base[name].Writable && (merged.Writable || policy.Writable)
			effective[name] = merged
		}
	}
	return effective
}

func (s *FieldPermissionService) ReadableFields(manifest resource.Manifest, policies map[string]FieldPolicy, export bool) []resource.Field {
	manifest = resource.NormalizeManifestFields(manifest)
	fields := make([]resource.Field, 0, len(manifest.Fields))
	for _, field := range manifest.Fields {
		policy := policies[field.Name]
		if !policy.Visible || !policy.Readable || (export && policy.Sensitive) {
			continue
		}
		fields = append(fields, field)
	}
	return fields
}

func (s *FieldPermissionService) ValidateQueryField(manifest resource.Manifest, policies map[string]FieldPolicy, name string, operation string) error {
	field, ok := fieldByName(manifest, name)
	if !ok {
		return fmt.Errorf("%w: %s", ErrFieldNotFound, name)
	}
	policy := policies[field.Name]
	if !policy.Visible || !policy.Readable || ((operation == "search" || operation == "sort" || operation == "export") && policy.Sensitive) {
		return fmt.Errorf("%w: %s", ErrFieldQueryDenied, name)
	}
	return nil
}

func (s *FieldPermissionService) ValidateWritablePayload(payload map[string]any, manifest resource.Manifest, policies map[string]FieldPolicy) (map[string]any, error) {
	values := make(map[string]any, len(payload))
	for name, value := range payload {
		field, ok := fieldByName(manifest, name)
		if !ok {
			continue
		}
		if !policies[field.Name].Writable {
			return nil, fmt.Errorf("%w: %s", ErrFieldPermissionDenied, name)
		}
		values[name] = value
	}
	return values, nil
}

func (s *FieldPermissionService) ValidateFieldOverrides(manifest resource.Manifest, overrides map[string]FieldOverride) error {
	policies := s.ManifestPolicies(manifest)
	for name, override := range overrides {
		policy, ok := policies[name]
		if !ok {
			return fmt.Errorf("%w: %s", ErrFieldNotFound, name)
		}
		if override.Readable && !policy.Readable {
			return fmt.Errorf("%w: %s", ErrFieldPolicyExpansion, name)
		}
		if override.Writable && !policy.Writable {
			return fmt.Errorf("%w: %s", ErrFieldPolicyExpansion, name)
		}
	}
	return nil
}

func (s *FieldPermissionService) ProjectRecord(record map[string]any, manifest resource.Manifest, policies map[string]FieldPolicy) map[string]any {
	projected := make(map[string]any, len(record))
	if id, ok := record["id"]; ok {
		projected["id"] = id
	}
	for _, field := range s.ReadableFields(manifest, policies, false) {
		if value, ok := record[field.Name]; ok {
			projected[field.Name] = value
		}
	}
	return projected
}

func fieldByName(manifest resource.Manifest, name string) (resource.Field, bool) {
	for _, field := range manifest.Fields {
		if field.Name == name {
			return field, true
		}
	}
	return resource.Field{}, false
}
