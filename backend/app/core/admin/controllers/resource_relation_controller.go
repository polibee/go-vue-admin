package controllers

import (
	"fmt"
	"strconv"

	"github.com/goravel/framework/contracts/database/orm"
	"github.com/goravel/framework/contracts/http"

	"goravel/app/core/resource"
	"goravel/app/facades"
	"goravel/app/modules/admin/registry"
	rbacservices "goravel/app/services/rbac"
)

type relationOption struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

func (r *ResourceController) RelationOptions(ctx http.Context) http.Response {
	source, err := registry.AdminRegistry().Find(ctx.Request().Route("resource"))
	if err != nil || source.Table == "" {
		return ctx.Response().Status(404).Json(http.Json{"code": "RESOURCE_NOT_FOUND"})
	}
	relation, ok := findRelation(source, ctx.Request().Route("relation"))
	if !ok || relation.Kind != "belongsTo" || !relation.Selectable {
		return ctx.Response().Status(404).Json(http.Json{"code": "RELATION_NOT_FOUND"})
	}
	if !hasResourceViewPermission(ctx, source) {
		return ctx.Response().Status(403).Json(http.Json{"code": "RBAC_FORBIDDEN"})
	}
	target, err := registry.AdminRegistry().Find(relation.Resource)
	if err != nil || target.Table == "" || !targetFieldDeclared(target, relation.ForeignField) || !targetFieldDeclared(target, relation.LabelField) {
		return ctx.Response().Status(404).Json(http.Json{"code": "RELATION_NOT_FOUND"})
	}
	labelField, readable := relationLabelField(target, relation.LabelField)
	if !readable {
		return ctx.Response().Status(403).Json(http.Json{"code": "RELATION_FORBIDDEN"})
	}
	permission := relation.Permission
	if permission == "" && len(target.Permissions) > 0 {
		permission = target.Permissions[0]
	}
	allowed, permissionErr := rbacservices.NewRBACService().UserHasPermission(ctx, permission)
	if permissionErr != nil {
		return ctx.Response().Status(401).Json(http.Json{"code": "AUTH_UNAUTHORIZED"})
	}
	if !allowed {
		return ctx.Response().Status(403).Json(http.Json{"code": "RELATION_FORBIDDEN"})
	}
	query, scopeErr := applyResourceScope(ctx, facades.Orm().Query().Table(target.Table), target, "view")
	if scopeErr != nil {
		return resourceScopeError(ctx, scopeErr)
	}
	var rows []map[string]any
	if err := paginateRelationOptions(query, relation.ForeignField, labelField, &rows); err != nil {
		return ctx.Response().Status(500).Json(http.Json{"code": "INTERNAL_ERROR"})
	}
	options := make([]relationOption, 0, len(rows))
	for _, row := range rows {
		value, valueOK := row[relation.ForeignField]
		label, labelOK := row[labelField]
		if !valueOK || !labelOK || value == nil || label == nil {
			continue
		}
		options = append(options, relationOption{Value: fmt.Sprint(value), Label: fmt.Sprint(label)})
	}
	return ctx.Response().Success().Json(http.Json{"data": options})
}

func (r *ResourceController) RelationRecords(ctx http.Context) http.Response {
	source, err := registry.AdminRegistry().Find(ctx.Request().Route("resource"))
	if err != nil || source.Table == "" {
		return ctx.Response().Status(404).Json(http.Json{"code": "RESOURCE_NOT_FOUND"})
	}
	relation, ok := findRelation(source, ctx.Request().Route("relation"))
	if !ok || relation.Kind != "hasMany" {
		return ctx.Response().Status(404).Json(http.Json{"code": "RELATION_NOT_FOUND"})
	}
	id, parseErr := strconv.ParseInt(ctx.Request().Route("id"), 10, 64)
	if parseErr != nil || id < 1 || !hasResourceViewPermission(ctx, source) {
		return ctx.Response().Status(404).Json(http.Json{"code": "RESOURCE_NOT_FOUND"})
	}
	if allowed, scopeErr := resourceCanAccess(ctx, source, "view", id); scopeErr != nil {
		return resourceScopeError(ctx, scopeErr)
	} else if !allowed {
		return ctx.Response().Status(404).Json(http.Json{"code": "RESOURCE_NOT_FOUND"})
	}
	target, err := registry.AdminRegistry().Find(relation.Resource)
	if err != nil || target.Table == "" || !targetFieldDeclared(target, relation.ForeignField) || !targetFieldReadable(target, relation.LabelField) {
		return ctx.Response().Status(404).Json(http.Json{"code": "RELATION_NOT_FOUND"})
	}
	permission := relation.Permission
	if permission == "" && len(target.Permissions) > 0 {
		permission = target.Permissions[0]
	}
	allowed, permissionErr := rbacservices.NewRBACService().UserHasPermission(ctx, permission)
	if permissionErr != nil {
		return ctx.Response().Status(401).Json(http.Json{"code": "AUTH_UNAUTHORIZED"})
	}
	if !allowed {
		return ctx.Response().Status(403).Json(http.Json{"code": "RELATION_FORBIDDEN"})
	}
	query, scopeErr := applyResourceScope(ctx, facades.Orm().Query().Table(target.Table), target, "view")
	if scopeErr != nil {
		return resourceScopeError(ctx, scopeErr)
	}
	query = query.Where(relation.ForeignField+" = ?", id)
	var rows []map[string]any
	if err := paginateRelationOptions(query, "id", relation.LabelField, &rows); err != nil {
		return ctx.Response().Status(500).Json(http.Json{"code": "INTERNAL_ERROR"})
	}
	options := make([]relationOption, 0, len(rows))
	for _, row := range rows {
		if row["id"] == nil || row[relation.LabelField] == nil {
			continue
		}
		options = append(options, relationOption{Value: fmt.Sprint(row["id"]), Label: fmt.Sprint(row[relation.LabelField])})
	}
	return ctx.Response().Success().Json(http.Json{"data": options})
}

func findRelation(manifest resource.Manifest, name string) (resource.Relation, bool) {
	for _, relation := range manifest.Relations {
		if relation.Name == name {
			return relation, true
		}
	}
	return resource.Relation{}, false
}

func hasResourceViewPermission(ctx http.Context, manifest resource.Manifest) bool {
	if len(manifest.Permissions) == 0 {
		return false
	}
	allowed, err := rbacservices.NewRBACService().UserHasPermission(ctx, manifest.Permissions[0])
	return err == nil && allowed
}

func targetFieldDeclared(manifest resource.Manifest, name string) bool {
	if name == "id" {
		return true
	}
	for _, field := range manifest.Fields {
		if field.Name == name {
			return true
		}
	}
	return false
}

func targetFieldReadable(manifest resource.Manifest, name string) bool {
	if name == "id" {
		return true
	}
	for _, field := range manifest.Fields {
		if field.Name == name {
			return field.Visible && field.Readable && !field.Sensitive
		}
	}
	return false
}

func relationLabelField(manifest resource.Manifest, name string) (string, bool) {
	return name, targetFieldReadable(manifest, name)
}

func paginateRelationOptions(query orm.Query, valueField, labelField string, rows *[]map[string]any) error {
	// Paginate keeps option payloads bounded without exposing arbitrary limits to the client.
	return query.Select(valueField, labelField).OrderBy(labelField, "asc").Paginate(1, 100, rows, new(int64))
}
