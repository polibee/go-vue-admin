package rbacservices

import (
	"errors"
	"fmt"

	"goravel/app/core/resource"
)

var (
	ErrFieldNotFound         = errors.New("resource field is not declared")
	ErrFieldQueryDenied      = errors.New("resource field query is not allowed")
	ErrFieldPermissionDenied = errors.New("resource field permission denied")
)

type FieldPolicy struct {
	Visible   bool
	Readable  bool
	Writable  bool
	Sensitive bool
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
