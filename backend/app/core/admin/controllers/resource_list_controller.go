package controllers

import (
	"strconv"
	"strings"

	"github.com/goravel/framework/contracts/database/orm"
	"github.com/goravel/framework/contracts/http"

	"goravel/app/core/resource"
	"goravel/app/facades"
	"goravel/app/models"
	"goravel/app/modules/admin/registry"
)

type resourceListQuery struct {
	Page    int
	PerPage int
	Search  string
	Status  string
	Sort    string
	Dir     string
}

func (r *ResourceController) List(ctx http.Context) http.Response {
	manifest, manifestErr := registry.AdminRegistry().Find(ctx.Request().Route("resource"))
	if manifestErr != nil || manifest.Table == "" {
		return ctx.Response().Status(404).Json(http.Json{"code": "RESOURCE_NOT_FOUND"})
	}
	query := resourceListQuery{
		Page:    positiveInt(ctx.Request().Query("page", "1"), 1),
		PerPage: positiveInt(ctx.Request().Query("per_page", "20"), 20),
		Search:  strings.TrimSpace(ctx.Request().Query("search")),
		Status:  normalizeUserStatusFilter(ctx.Request().Query("status")),
		Sort:    ctx.Request().Query("sort", "id"),
		Dir:     strings.ToLower(ctx.Request().Query("dir", "desc")),
	}
	if query.PerPage > 100 {
		query.PerPage = 100
	}
	if query.Dir != "asc" {
		query.Dir = "desc"
	}

	var rows any
	var q orm.Query
	var total int64
	switch ctx.Request().Route("resource") {
	case "users":
		rows = &[]models.User{}
		q = facades.Orm().Query()
		q, manifestErr = applyResourceScope(ctx, q, manifest, "view")
		if manifestErr != nil {
			return resourceScopeError(ctx, manifestErr)
		}
		q = applyResourceSearch(q, query.Search, "name", "email")
		if query.Status != "" {
			q = q.Where("status = ?", query.Status)
		}
		query.Sort = allowedSort(query.Sort, map[string]bool{"id": true, "name": true, "email": true, "status": true}, "id")
	case "roles":
		rows = &[]models.Role{}
		q = facades.Orm().Query()
		q, manifestErr = applyResourceScope(ctx, q, manifest, "view")
		if manifestErr != nil {
			return resourceScopeError(ctx, manifestErr)
		}
		q = applyResourceSearch(q, query.Search, "name", "display_name")
		query.Sort = allowedSort(query.Sort, map[string]bool{"id": true, "name": true, "display_name": true}, "id")
	case "permissions":
		rows = &[]models.Permission{}
		q = facades.Orm().Query()
		q, manifestErr = applyResourceScope(ctx, q, manifest, "view")
		if manifestErr != nil {
			return resourceScopeError(ctx, manifestErr)
		}
		q = applyResourceSearch(q, query.Search, "name", "display_name")
		query.Sort = allowedSort(query.Sort, map[string]bool{"id": true, "name": true, "display_name": true}, "id")
	default:
		rows = &[]map[string]any{}
		q = facades.Orm().Query().Table(manifest.Table)
		q, manifestErr = applyResourceScope(ctx, q, manifest, "view")
		if manifestErr != nil {
			return resourceScopeError(ctx, manifestErr)
		}
		searchColumns := make([]string, 0, len(manifest.Columns))
		allowedColumns := make(map[string]bool, len(manifest.Columns)+1)
		allowedColumns["id"] = true
		for _, column := range manifest.Columns {
			if resourceFieldAllowedForQuery(manifest, column.Name, "sort") {
				allowedColumns[column.Name] = true
			}
			searchColumns = append(searchColumns, column.Name)
		}
		q = applyResourceSearch(q, query.Search, fieldNames(resourceSearchFields(manifest))...)
		for _, field := range resourceFilterFields(manifest) {
			value := strings.TrimSpace(ctx.Request().Query(field.Name))
			if value != "" && value != "all" && resourceFilterValueAllowed(field, value) {
				q = q.Where(field.Name+" = ?", value)
			}
		}
		query.Sort = allowedSort(query.Sort, allowedColumns, "id")
		q = q.OrderBy(query.Sort, query.Dir)
		if err := q.Paginate(query.Page, query.PerPage, rows, &total); err != nil {
			return ctx.Response().Status(500).Json(http.Json{"code": "INTERNAL_ERROR"})
		}
		return resourceListResponse(ctx, rows, query, total)
	}

	q = q.OrderBy(query.Sort, query.Dir)
	if err := q.Paginate(query.Page, query.PerPage, rows, &total); err != nil {
		return ctx.Response().Status(500).Json(http.Json{"code": "INTERNAL_ERROR"})
	}

	if users, ok := rows.(*[]models.User); ok {
		items := make([]map[string]any, 0, len(*users))
		for _, user := range *users {
			items = append(items, projectResourceValue(user.Public(), manifest, false))
		}
		return resourceListResponse(ctx, items, query, total)
	}
	if roles, ok := rows.(*[]models.Role); ok {
		items := make([]map[string]any, 0, len(*roles))
		for _, role := range *roles {
			items = append(items, projectResourceValue(role, manifest, false))
		}
		return resourceListResponse(ctx, items, query, total)
	}
	if permissions, ok := rows.(*[]models.Permission); ok {
		items := make([]map[string]any, 0, len(*permissions))
		for _, permission := range *permissions {
			items = append(items, projectResourceValue(permission, manifest, false))
		}
		return resourceListResponse(ctx, items, query, total)
	}
	return resourceListResponse(ctx, rows, query, total)
}

func normalizeUserStatusFilter(status string) string {
	switch strings.TrimSpace(status) {
	case "active", "disabled", "locked":
		return strings.TrimSpace(status)
	default:
		return ""
	}
}

func applyResourceSearch(query orm.Query, search string, columns ...string) orm.Query {
	if search == "" || len(columns) == 0 {
		return query
	}
	like := "%" + search + "%"
	query = query.Where(columns[0]+" LIKE ?", like)
	for _, column := range columns[1:] {
		query = query.OrWhere(column+" LIKE ?", like)
	}
	return query
}

func resourceSearchFields(manifest resource.Manifest) []resource.Field {
	fields := make([]resource.Field, 0)
	for _, field := range manifest.Fields {
		if (field.Type == "text" || field.Type == "email") && resourceFieldAllowedForQuery(manifest, field.Name, "search") {
			fields = append(fields, field)
		}
	}
	return fields
}

func resourceFilterFields(manifest resource.Manifest) []resource.Field {
	fields := make([]resource.Field, 0)
	for _, field := range manifest.Fields {
		if (field.Type == "select" || field.Type == "boolean") && field.Visible && field.Readable {
			fields = append(fields, field)
		}
	}
	return fields
}

func fieldNames(fields []resource.Field) []string {
	names := make([]string, 0, len(fields))
	for _, field := range fields {
		names = append(names, field.Name)
	}
	return names
}

func resourceFilterValueAllowed(field resource.Field, value string) bool {
	if field.Type == "boolean" {
		return value == "true" || value == "false" || value == "1" || value == "0"
	}
	if len(field.Options) == 0 {
		return true
	}
	for _, option := range field.Options {
		if option.Value == value {
			return true
		}
	}
	return false
}

func allowedSort(value string, allowed map[string]bool, fallback string) string {
	if allowed[value] {
		return value
	}
	return fallback
}

func positiveInt(value string, fallback int) int {
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed < 1 {
		return fallback
	}
	return parsed
}

func resourceListResponse(ctx http.Context, data any, query resourceListQuery, total int64) http.Response {
	lastPage := int((total + int64(query.PerPage) - 1) / int64(query.PerPage))
	if lastPage == 0 {
		lastPage = 1
	}
	return ctx.Response().Success().Json(http.Json{
		"data": data,
		"meta": http.Json{
			"page": query.Page, "per_page": query.PerPage, "total": total, "last_page": lastPage,
		},
	})
}
