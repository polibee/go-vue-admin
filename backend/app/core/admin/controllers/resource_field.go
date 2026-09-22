package controllers

import (
	"encoding/json"

	"github.com/goravel/framework/contracts/http"

	"goravel/app/core/resource"
	rbacservices "goravel/app/services/rbac"
)

func resourceFieldPolicies(manifest resource.Manifest) map[string]rbacservices.FieldPolicy {
	return rbacservices.NewFieldPermissionService().ManifestPolicies(manifest)
}

func resourceFieldPoliciesFor(ctx http.Context, manifest resource.Manifest, action string) (map[string]rbacservices.FieldPolicy, error) {
	return rbacservices.NewFieldPermissionService().EffectivePolicies(ctx, manifest, action)
}

func projectResourceRecord(record map[string]any, manifest resource.Manifest, export bool) map[string]any {
	service := rbacservices.NewFieldPermissionService()
	policies := service.ManifestPolicies(manifest)
	return projectResourceRecordWithPolicies(record, manifest, policies, export)
}

func projectResourceRecordWithPolicies(record map[string]any, manifest resource.Manifest, policies map[string]rbacservices.FieldPolicy, export bool) map[string]any {
	service := rbacservices.NewFieldPermissionService()
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
	return projectResourceRowsWithPolicies(rows, manifest, resourceFieldPolicies(manifest), export)
}

func projectResourceRowsWithPolicies(rows []map[string]any, manifest resource.Manifest, policies map[string]rbacservices.FieldPolicy, export bool) []map[string]any {
	projected := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		projected = append(projected, projectResourceRecordWithPolicies(row, manifest, policies, export))
	}
	return projected
}

func projectResourceValue(value any, manifest resource.Manifest, export bool) map[string]any {
	return projectResourceValueWithPolicies(value, manifest, resourceFieldPolicies(manifest), export)
}

func projectResourceValueWithPolicies(value any, manifest resource.Manifest, policies map[string]rbacservices.FieldPolicy, export bool) map[string]any {
	if record, ok := value.(map[string]any); ok {
		return projectResourceRecordWithPolicies(record, manifest, policies, export)
	}
	encoded, err := json.Marshal(value)
	if err != nil {
		return map[string]any{}
	}
	var record map[string]any
	if err := json.Unmarshal(encoded, &record); err != nil {
		return map[string]any{}
	}
	return projectResourceRecordWithPolicies(record, manifest, policies, export)
}

func resourceFieldAllowedForQuery(manifest resource.Manifest, name string, operation string) bool {
	service := rbacservices.NewFieldPermissionService()
	return service.ValidateQueryField(manifest, service.ManifestPolicies(manifest), name, operation) == nil
}
