package controllers

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/goravel/framework/contracts/http"

	"goravel/app/core/resource"
	"goravel/app/facades"
	"goravel/app/modules/admin/registry"
	rbacservices "goravel/app/services/rbac"
	userservices "goravel/app/services/users"
)

func (r *ResourceController) Create(ctx http.Context) http.Response {
	manifest, ok := generatedManifest(ctx)
	if !ok {
		return ctx.Response().Status(404).Json(http.Json{"code": "RESOURCE_NOT_FOUND"})
	}
	if err := validateResourceActionScope(ctx, manifest, "create"); err != nil {
		return resourceScopeError(ctx, err)
	}
	if manifest.Name == "users" {
		var payload userPayload
		if err := ctx.Request().Bind(&payload); err != nil {
			return rbacError(ctx, 422, "VALIDATION_ERROR")
		}
		if err := validateSpecializedPayload(ctx, manifest, "create", map[string]any{"name": payload.Name, "email": payload.Email, "password": payload.Password, "locale": payload.Locale, "status": payload.Status}); err != nil {
			return fieldWriteError(ctx, err)
		}
		user, err := userservices.NewUserService().Create(payload.Name, payload.Email, payload.Password, payload.Locale, payload.Status)
		if err != nil {
			return userServiceError(ctx, err)
		}
		recordManagementAudit(ctx, "user.create", map[string]any{"target_user_id": user.ID})
		return ctx.Response().Status(201).Json(http.Json{"data": user.Public()})
	}
	if manifest.Name == "roles" {
		var payload rolePayload
		if err := ctx.Request().Bind(&payload); err != nil {
			return rbacError(ctx, 422, "VALIDATION_ERROR")
		}
		if err := validateSpecializedPayload(ctx, manifest, "create", map[string]any{"name": payload.Name, "display_name": payload.DisplayName}); err != nil {
			return fieldWriteError(ctx, err)
		}
		role, err := rbacservices.NewRoleService().Create(payload.Name, payload.DisplayName)
		if err != nil {
			return roleServiceError(ctx, err)
		}
		recordManagementAudit(ctx, "role.create", map[string]any{"target_role_id": role.ID})
		return ctx.Response().Status(201).Json(http.Json{"data": role})
	}
	values, err := bindGeneratedValues(ctx, manifest, "create")
	if err != nil {
		return ctx.Response().Status(422).Json(http.Json{"code": "VALIDATION_ERROR"})
	}
	if err := enforceResourceCreateOwner(ctx, manifest, values); err != nil {
		return resourceScopeError(ctx, err)
	}
	if err := validateResourceRelations(ctx, manifest, values); err != nil {
		return relationWriteError(ctx, err)
	}
	createdID, err := insertGeneratedResource(manifest, values)
	if err != nil {
		return ctx.Response().Status(500).Json(http.Json{"code": "INTERNAL_ERROR"})
	}
	values["id"] = createdID
	return ctx.Response().Status(201).Json(http.Json{"data": values})
}

func insertGeneratedResource(manifest resource.Manifest, values map[string]any) (int64, error) {
	columns := make([]string, 0, len(values))
	for column := range values {
		columns = append(columns, column)
	}
	sort.Strings(columns)
	placeholders := make([]string, len(columns))
	args := make([]any, len(columns))
	for index, column := range columns {
		placeholders[index] = "?"
		args[index] = values[column]
	}
	// Manifest table and field names are generated from validated identifiers; values
	// remain bound parameters. PostgreSQL RETURNING keeps the generic API response
	// deterministic without a race-prone follow-up lookup.
	statement := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) RETURNING id", manifest.Table, strings.Join(columns, ", "), strings.Join(placeholders, ", "))
	var id int64
	if err := facades.Orm().Query().Raw(statement, args...).Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *ResourceController) Update(ctx http.Context) http.Response {
	manifest, ok := generatedManifest(ctx)
	if !ok {
		return ctx.Response().Status(404).Json(http.Json{"code": "RESOURCE_NOT_FOUND"})
	}
	id := resourceID(ctx)
	if id < 1 {
		return ctx.Response().Status(404).Json(http.Json{"code": "RESOURCE_NOT_FOUND"})
	}
	allowed, scopeErr := resourceCanAccess(ctx, manifest, "update", id)
	if scopeErr != nil {
		return resourceScopeError(ctx, scopeErr)
	}
	if !allowed {
		return ctx.Response().Status(404).Json(http.Json{"code": "RESOURCE_NOT_FOUND"})
	}
	if manifest.Name == "users" {
		var payload userPayload
		if err := ctx.Request().Bind(&payload); err != nil {
			return rbacError(ctx, 422, "VALIDATION_ERROR")
		}
		if err := validateSpecializedPayload(ctx, manifest, "update", map[string]any{"name": payload.Name, "email": payload.Email, "password": payload.Password, "locale": payload.Locale, "status": payload.Status}); err != nil {
			return fieldWriteError(ctx, err)
		}
		user, err := userservices.NewUserService().Update(id, payload.Name, payload.Email, payload.Password, payload.Locale, payload.Status)
		if err != nil {
			return userServiceError(ctx, err)
		}
		recordManagementAudit(ctx, "user.update", map[string]any{"target_user_id": user.ID})
		return ctx.Response().Success().Json(http.Json{"data": user.Public()})
	}
	if manifest.Name == "roles" {
		var payload rolePayload
		if err := ctx.Request().Bind(&payload); err != nil {
			return rbacError(ctx, 422, "VALIDATION_ERROR")
		}
		if err := validateSpecializedPayload(ctx, manifest, "update", map[string]any{"name": payload.Name, "display_name": payload.DisplayName}); err != nil {
			return fieldWriteError(ctx, err)
		}
		role, err := rbacservices.NewRoleService().Update(id, payload.Name, payload.DisplayName)
		if err != nil {
			return roleServiceError(ctx, err)
		}
		recordManagementAudit(ctx, "role.update", map[string]any{"target_role_id": role.ID})
		return ctx.Response().Success().Json(http.Json{"data": role})
	}
	values, err := bindGeneratedValues(ctx, manifest, "update")
	if err != nil {
		return ctx.Response().Status(422).Json(http.Json{"code": "VALIDATION_ERROR"})
	}
	if err := validateResourceRelations(ctx, manifest, values); err != nil {
		return relationWriteError(ctx, err)
	}
	if _, err := facades.Orm().Query().Table(manifest.Table).Where("id = ?", id).Update(values); err != nil {
		return ctx.Response().Status(500).Json(http.Json{"code": "INTERNAL_ERROR"})
	}
	values["id"] = id
	return ctx.Response().Success().Json(http.Json{"data": values})
}

func (r *ResourceController) Delete(ctx http.Context) http.Response {
	manifest, ok := generatedManifest(ctx)
	if !ok {
		return ctx.Response().Status(404).Json(http.Json{"code": "RESOURCE_NOT_FOUND"})
	}
	id := resourceID(ctx)
	if id < 1 {
		return ctx.Response().Status(404).Json(http.Json{"code": "RESOURCE_NOT_FOUND"})
	}
	allowed, scopeErr := resourceCanAccess(ctx, manifest, "delete", id)
	if scopeErr != nil {
		return resourceScopeError(ctx, scopeErr)
	}
	if !allowed {
		return ctx.Response().Status(404).Json(http.Json{"code": "RESOURCE_NOT_FOUND"})
	}
	switch resourceDeleteStrategy(manifest) {
	case "users-service":
		if err := userservices.NewUserService().Delete(id); err != nil {
			return userServiceError(ctx, err)
		}
		recordManagementAudit(ctx, "user.delete", map[string]any{"target_user_id": id})
		return ctx.Response().NoContent(204)
	case "roles-service":
		if err := rbacservices.NewRoleService().Delete(id); err != nil {
			return roleServiceError(ctx, err)
		}
		recordManagementAudit(ctx, "role.delete", map[string]any{"target_role_id": id})
		return ctx.Response().NoContent(204)
	}
	if manifest.SoftDelete {
		if _, err := facades.Orm().Query().Table(manifest.Table).Where("id = ?", id).Update(map[string]any{"deleted_at": time.Now().UTC()}); err != nil {
			return ctx.Response().Status(500).Json(http.Json{"code": "INTERNAL_ERROR"})
		}
		return ctx.Response().NoContent(204)
	}
	if _, err := facades.Orm().Query().Table(manifest.Table).Where("id = ?", id).Delete(); err != nil {
		return ctx.Response().Status(500).Json(http.Json{"code": "INTERNAL_ERROR"})
	}
	return ctx.Response().NoContent(204)
}

func resourceDeleteStrategy(manifest resource.Manifest) string {
	switch manifest.Name {
	case "users":
		return "users-service"
	case "roles":
		return "roles-service"
	default:
		return "generic-table"
	}
}

func generatedManifest(ctx http.Context) (resource.Manifest, bool) {
	manifest, err := registry.AdminRegistry().Find(resourceName(ctx))
	return manifest, err == nil && manifest.Table != ""
}

func bindGeneratedValues(ctx http.Context, manifest resource.Manifest, action string) (map[string]any, error) {
	var payload map[string]any
	if err := ctx.Request().Bind(&payload); err != nil {
		return nil, err
	}
	allowed := make(map[string]bool, len(manifest.Fields))
	for _, field := range manifest.Fields {
		allowed[field.Name] = true
	}
	values, err := validateGeneratedValues(payload, manifest)
	if err != nil {
		return nil, err
	}
	policies, err := resourceFieldPoliciesFor(ctx, manifest, action)
	if err != nil {
		return nil, err
	}
	return rbacservices.NewFieldPermissionService().ValidateWritablePayload(values, manifest, policies)
}

func validateSpecializedPayload(ctx http.Context, manifest resource.Manifest, action string, payload map[string]any) error {
	policies, err := resourceFieldPoliciesFor(ctx, manifest, action)
	if err != nil {
		return err
	}
	_, err = rbacservices.NewFieldPermissionService().ValidateWritablePayload(payload, manifest, policies)
	return err
}

func fieldWriteError(ctx http.Context, err error) http.Response {
	if errors.Is(err, rbacservices.ErrFieldPermissionDenied) || errors.Is(err, rbacservices.ErrFieldPolicyExpansion) || errors.Is(err, rbacservices.ErrFieldNotFound) {
		return ctx.Response().Status(422).Json(http.Json{"code": "FIELD_PERMISSION_DENIED"})
	}
	return resourceScopeError(ctx, err)
}

func validateGeneratedValues(payload map[string]any, manifest resource.Manifest) (map[string]any, error) {
	values := make(map[string]any)
	fields := make(map[string]resource.Field, len(manifest.Fields))
	for _, field := range manifest.Fields {
		fields[field.Name] = field
	}
	for key, value := range payload {
		field, ok := fields[key]
		if !ok {
			continue
		}
		if err := validateGeneratedField(field, value); err != nil {
			return nil, err
		}
		values[key] = value
	}
	for _, field := range manifest.Fields {
		if field.Required {
			value, exists := values[field.Name]
			if !exists || value == nil || (field.Type != "boolean" && strings.TrimSpace(fmt.Sprint(value)) == "") {
				return nil, fmt.Errorf("required field %q is missing", field.Name)
			}
		}
	}
	return values, nil
}

func validateGeneratedField(field resource.Field, value any) error {
	if value == nil {
		return nil
	}
	switch field.Type {
	case "boolean":
		if _, ok := value.(bool); !ok {
			return fmt.Errorf("field %q must be boolean", field.Name)
		}
	case "integer":
		switch value.(type) {
		case float64, float32, int, int32, int64:
		default:
			return fmt.Errorf("field %q must be integer", field.Name)
		}
	case "select":
		text, ok := value.(string)
		if !ok {
			return fmt.Errorf("field %q must be a select value", field.Name)
		}
		if len(field.Options) > 0 {
			for _, option := range field.Options {
				if option.Value == text {
					return nil
				}
			}
			return fmt.Errorf("field %q has an invalid option", field.Name)
		}
	}
	return nil
}

var (
	errRelationNotFound  = errors.New("relation target not found")
	errRelationForbidden = errors.New("relation target forbidden")
)

func validateResourceRelations(ctx http.Context, manifest resource.Manifest, values map[string]any) error {
	for _, relation := range manifest.Relations {
		if relation.Kind != "belongsTo" || !relation.Selectable {
			continue
		}
		value, exists := values[relation.Field]
		if !exists || value == nil || value == "" {
			continue
		}
		target, err := registry.AdminRegistry().Find(relation.Resource)
		if err != nil || target.Table == "" || !targetFieldDeclared(target, relation.ForeignField) {
			return errRelationNotFound
		}
		permission := relation.Permission
		if permission == "" && len(target.Permissions) > 0 {
			permission = target.Permissions[0]
		}
		allowed, permissionErr := rbacservices.NewRBACService().UserHasPermission(ctx, permission)
		if permissionErr != nil || !allowed {
			return errRelationForbidden
		}
		query, scopeErr := applyResourceScope(ctx, facades.Orm().Query().Table(target.Table), target, "view")
		if scopeErr != nil {
			return errRelationForbidden
		}
		exists, existsErr := query.Where(relation.ForeignField+" = ?", value).Exists()
		if existsErr != nil {
			return errRelationForbidden
		}
		if !exists {
			return errRelationNotFound
		}
	}
	return nil
}

func relationWriteError(ctx http.Context, err error) http.Response {
	if errors.Is(err, errRelationForbidden) {
		return ctx.Response().Status(422).Json(http.Json{"code": "RELATION_FORBIDDEN"})
	}
	if errors.Is(err, errRelationNotFound) {
		return ctx.Response().Status(422).Json(http.Json{"code": "RELATION_NOT_FOUND"})
	}
	return ctx.Response().Status(422).Json(http.Json{"code": "VALIDATION_ERROR"})
}
