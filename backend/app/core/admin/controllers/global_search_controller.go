package controllers

import (
	"fmt"
	"strings"

	"github.com/goravel/framework/contracts/database/orm"
	"github.com/goravel/framework/contracts/http"

	"goravel/app/core/resource"
	"goravel/app/facades"
	"goravel/app/models"
	"goravel/app/modules/admin/registry"
	rbacservices "goravel/app/services/rbac"
)

const (
	globalSearchMaxQuery       = 100
	globalSearchPerResource    = 5
	globalSearchMaxResultCount = 20
)

type GlobalSearchController struct{}

type globalSearchResult struct {
	Resource string `json:"resource"`
	Label    string `json:"label"`
	ID       any    `json:"id"`
	Title    string `json:"title"`
	Subtitle string `json:"subtitle,omitempty"`
	Route    string `json:"route"`
}

func NewGlobalSearchController() *GlobalSearchController { return &GlobalSearchController{} }

func (c *GlobalSearchController) Index(ctx http.Context) http.Response {
	query, ok := normalizeGlobalSearchQuery(ctx.Request().Query("q"))
	if !ok {
		return ctx.Response().Success().Json(http.Json{"data": []globalSearchResult{}})
	}

	results := make([]globalSearchResult, 0, globalSearchMaxResultCount)
	for _, manifest := range registry.AdminRegistry().All() {
		if len(results) >= globalSearchMaxResultCount || !globalSearchResourceAllowed(ctx, manifest.Permissions) {
			continue
		}
		items, err := searchManifest(manifest, query)
		if err != nil {
			return ctx.Response().Status(500).Json(http.Json{"code": "INTERNAL_ERROR"})
		}
		remaining := globalSearchMaxResultCount - len(results)
		if len(items) > remaining {
			items = items[:remaining]
		}
		results = append(results, items...)
	}

	return ctx.Response().Success().Json(http.Json{"data": results})
}

func normalizeGlobalSearchQuery(value string) (string, bool) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", false
	}
	if len(value) > globalSearchMaxQuery {
		value = value[:globalSearchMaxQuery]
	}
	return value, true
}

func globalSearchResourceAllowed(ctx http.Context, permissions []string) bool {
	if len(permissions) == 0 {
		return false
	}
	service := rbacservices.NewRBACService()
	for _, permission := range permissions {
		if strings.TrimSpace(permission) == "" {
			continue
		}
		allowed, err := service.UserHasPermission(ctx, permission)
		if err == nil && allowed {
			return true
		}
	}
	return false
}

func searchManifest(manifest resource.Manifest, query string) ([]globalSearchResult, error) {
	searchFields := globalSearchFields(manifest)
	if len(searchFields) == 0 || manifest.Table == "" {
		return []globalSearchResult{}, nil
	}

	like := "%" + query + "%"
	switch manifest.Name {
	case "users":
		var rows []models.User
		q := applyGlobalSearch(facades.Orm().Query(), like, "name", "email")
		if err := q.OrderBy("id", "desc").Paginate(1, globalSearchPerResource, &rows, new(int64)); err != nil {
			return nil, err
		}
		results := make([]globalSearchResult, 0, len(rows))
		for _, row := range rows {
			results = append(results, makeGlobalSearchResult(manifest, row.ID, row.Name, row.Email))
		}
		return results, nil
	case "roles", "permissions":
		var rows []models.Role
		if manifest.Name == "permissions" {
			var permissionRows []models.Permission
			q := applyGlobalSearch(facades.Orm().Query(), like, "name", "display_name")
			if err := q.OrderBy("id", "desc").Paginate(1, globalSearchPerResource, &permissionRows, new(int64)); err != nil {
				return nil, err
			}
			results := make([]globalSearchResult, 0, len(permissionRows))
			for _, row := range permissionRows {
				results = append(results, makeGlobalSearchResult(manifest, row.ID, row.DisplayName, row.Name))
			}
			return results, nil
		}
		q := applyGlobalSearch(facades.Orm().Query(), like, "name", "display_name")
		if err := q.OrderBy("id", "desc").Paginate(1, globalSearchPerResource, &rows, new(int64)); err != nil {
			return nil, err
		}
		results := make([]globalSearchResult, 0, len(rows))
		for _, row := range rows {
			results = append(results, makeGlobalSearchResult(manifest, row.ID, row.DisplayName, row.Name))
		}
		return results, nil
	default:
		columns := []string{"id"}
		for _, field := range searchFields {
			columns = append(columns, field.Name)
		}
		rows := []map[string]any{}
		q := facades.Orm().Query().Table(manifest.Table).Select(columns...)
		q = applyGlobalSearch(q, like, fieldNames(searchFields)...)
		if err := q.OrderBy("id", "desc").Paginate(1, globalSearchPerResource, &rows, new(int64)); err != nil {
			return nil, err
		}
		results := make([]globalSearchResult, 0, len(rows))
		for _, row := range rows {
			if id, exists := row["id"]; exists {
				title, subtitle := globalSearchText(row, searchFields)
				results = append(results, makeGlobalSearchResult(manifest, id, title, subtitle))
			}
		}
		return results, nil
	}
}

func globalSearchFields(manifest resource.Manifest) []resource.Field {
	fields := make([]resource.Field, 0)
	for _, field := range manifest.Fields {
		if (field.Type == "text" || field.Type == "email") && !isSensitiveSearchName(field.Name) {
			fields = append(fields, field)
		}
	}
	return fields
}

func isSensitiveSearchName(name string) bool {
	name = strings.ToLower(name)
	for _, part := range []string{"password", "secret", "token", "credential", "private_key", "access_key"} {
		if strings.Contains(name, part) {
			return true
		}
	}
	return false
}

func applyGlobalSearch(query orm.Query, like string, columns ...string) orm.Query {
	if len(columns) == 0 {
		return query
	}
	query = query.Where(columns[0]+" LIKE ?", like)
	for _, column := range columns[1:] {
		query = query.OrWhere(column+" LIKE ?", like)
	}
	return query
}

func makeGlobalSearchResult(manifest resource.Manifest, id any, title, subtitle string) globalSearchResult {
	return globalSearchResult{
		Resource: manifest.Name,
		Label:    manifest.Label,
		ID:       id,
		Title:    title,
		Subtitle: subtitle,
		Route:    frontendSearchRoute(manifest.Route, id),
	}
}

func frontendSearchRoute(route string, id any) string {
	route = strings.TrimSuffix(strings.TrimSpace(route), "/")
	route = strings.TrimPrefix(route, "/admin")
	if route == "" {
		route = "/"
	}
	return strings.TrimSuffix(route, "/") + "/" + urlPathPart(id)
}

func urlPathPart(value any) string { return strings.TrimSpace(fmt.Sprint(value)) }

func globalSearchText(row map[string]any, fields []resource.Field) (string, string) {
	values := make([]string, 0, 2)
	for _, field := range fields {
		if value, ok := row[field.Name]; ok {
			text := strings.TrimSpace(fmt.Sprint(value))
			if text != "" {
				values = append(values, text)
			}
		}
		if len(values) == 2 {
			break
		}
	}
	if len(values) == 0 {
		return "", ""
	}
	if len(values) == 1 {
		return values[0], ""
	}
	return values[0], values[1]
}
