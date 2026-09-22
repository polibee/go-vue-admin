package controllers

import (
	"fmt"
	"strconv"
	"strings"

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

type relationOptionsQuery struct {
	Search   string
	Selected string
	Page     int
	PerPage  int
}

func normalizeRelationOptionsQuery(search, selected, pageValue, perPageValue string) relationOptionsQuery {
	return relationOptionsQuery{
		Search:   strings.TrimSpace(search),
		Selected: strings.TrimSpace(selected),
		Page:     positiveInt(pageValue, 1),
		PerPage:  minPositiveInt(perPageValue, 100, 20),
	}
}

func minPositiveInt(value string, maximum, fallback int) int {
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed < 1 {
		return fallback
	}
	if parsed > maximum {
		return maximum
	}
	return parsed
}

func (r *ResourceController) RelationOptions(ctx http.Context) http.Response {
	source, err := registry.AdminRegistry().Find(resourceName(ctx))
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
	optionsQuery := normalizeRelationOptionsQuery(ctx.Request().Query("search"), ctx.Request().Query("selected"), ctx.Request().Query("page", "1"), ctx.Request().Query("per_page", "20"))
	if optionsQuery.Search != "" {
		query = query.Where(labelField+" LIKE ?", "%"+optionsQuery.Search+"%")
	}
	var rows []map[string]any
	var total int64
	if err := paginateRelationOptions(query, relation.ForeignField, labelField, optionsQuery.Page, optionsQuery.PerPage, &rows, &total); err != nil {
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
	if optionsQuery.Selected != "" && !relationOptionExists(options, optionsQuery.Selected) {
		selectedQuery, selectedScopeErr := applyResourceScope(ctx, facades.Orm().Query().Table(target.Table), target, "view")
		if selectedScopeErr != nil {
			return resourceScopeError(ctx, selectedScopeErr)
		}
		selectedQuery = selectedQuery.Where(relation.ForeignField+" = ?", optionsQuery.Selected)
		var selectedRows []map[string]any
		var selectedTotal int64
		if selectedErr := paginateRelationOptions(selectedQuery, relation.ForeignField, labelField, 1, 1, &selectedRows, &selectedTotal); selectedErr != nil {
			return ctx.Response().Status(500).Json(http.Json{"code": "INTERNAL_ERROR"})
		}
		for _, row := range selectedRows {
			if value, valueOK := row[relation.ForeignField]; valueOK && row[labelField] != nil {
				options = append(options, relationOption{Value: fmt.Sprint(value), Label: fmt.Sprint(row[labelField])})
			}
		}
	}
	return ctx.Response().Success().Json(http.Json{"data": options, "meta": relationOptionsMeta(optionsQuery, total)})
}

func (r *ResourceController) RelationRecords(ctx http.Context) http.Response {
	source, err := registry.AdminRegistry().Find(resourceName(ctx))
	if err != nil || source.Table == "" {
		return ctx.Response().Status(404).Json(http.Json{"code": "RESOURCE_NOT_FOUND"})
	}
	relation, ok := findRelation(source, ctx.Request().Route("relation"))
	if !ok || relation.Kind != "hasMany" {
		return ctx.Response().Status(404).Json(http.Json{"code": "RELATION_NOT_FOUND"})
	}
	id := resourceID(ctx)
	if id < 1 || !hasResourceViewPermission(ctx, source) {
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
	var total int64
	if err := paginateRelationOptions(query, "id", relation.LabelField, 1, 100, &rows, &total); err != nil {
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

func paginateRelationOptions(query orm.Query, valueField, labelField string, page, perPage int, rows *[]map[string]any, total *int64) error {
	// Paginate keeps option payloads bounded without exposing arbitrary limits to the client.
	return query.Select(valueField, labelField).OrderBy(labelField, "asc").Paginate(page, perPage, rows, total)
}

func relationOptionsMeta(query relationOptionsQuery, total int64) http.Json {
	lastPage := int((total + int64(query.PerPage) - 1) / int64(query.PerPage))
	if lastPage == 0 {
		lastPage = 1
	}
	return http.Json{"page": query.Page, "per_page": query.PerPage, "total": total, "last_page": lastPage}
}

func relationOptionExists(options []relationOption, value string) bool {
	for _, option := range options {
		if option.Value == value {
			return true
		}
	}
	return false
}
