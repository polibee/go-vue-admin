package controllers

import (
	"encoding/json"

	"goravel/app/core/resource"
	rbacservices "goravel/app/services/rbac"
)

func resourceFieldPolicies(manifest resource.Manifest) map[string]rbacservices.FieldPolicy {
	return rbacservices.NewFieldPermissionService().ManifestPolicies(manifest)
}

func projectResourceRecord(record map[string]any, manifest resource.Manifest, export bool) map[string]any {
	service := rbacservices.NewFieldPermissionService()
	policies := service.ManifestPolicies(manifest)
	if export {
		projected := make(map[string]any, len(record))
		if id, ok := record["id"]; ok {
			projected["id"] = id
		}
		for _, field := range service.ReadableFields(manifest, policies, true) {
			if value, ok := record[field.Name]; ok {
				projected[field.Name] = value
			}
		}
		return projected
	}
	return service.ProjectRecord(record, manifest, policies)
}

func projectResourceRows(rows []map[string]any, manifest resource.Manifest, export bool) []map[string]any {
	projected := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		projected = append(projected, projectResourceRecord(row, manifest, export))
	}
	return projected
}

func projectResourceValue(value any, manifest resource.Manifest, export bool) map[string]any {
	if record, ok := value.(map[string]any); ok {
		return projectResourceRecord(record, manifest, export)
	}
	encoded, err := json.Marshal(value)
	if err != nil {
		return map[string]any{}
	}
	var record map[string]any
	if err := json.Unmarshal(encoded, &record); err != nil {
		return map[string]any{}
	}
	return projectResourceRecord(record, manifest, export)
}

func resourceFieldAllowedForQuery(manifest resource.Manifest, name string, operation string) bool {
	service := rbacservices.NewFieldPermissionService()
	return service.ValidateQueryField(manifest, service.ManifestPolicies(manifest), name, operation) == nil
}
