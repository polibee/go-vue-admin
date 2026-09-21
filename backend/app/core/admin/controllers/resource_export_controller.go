package controllers

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strings"

	"github.com/goravel/framework/contracts/database/orm"
	"github.com/goravel/framework/contracts/http"

	"goravel/app/core/resource"
	"goravel/app/facades"
	"goravel/app/models"
	"goravel/app/modules/admin/registry"
)

func (r *ResourceController) Export(ctx http.Context) http.Response {
	resourceName := ctx.Request().Route("resource")
	manifest, err := registry.AdminRegistry().Find(resourceName)
	if err != nil || manifest.Table == "" {
		return ctx.Response().Status(404).Json(http.Json{"code": "RESOURCE_NOT_FOUND"})
	}
	query := resourceListQuery{
		Search: strings.TrimSpace(ctx.Request().Query("search")),
		Status: normalizeUserStatusFilter(ctx.Request().Query("status")),
		Sort:   ctx.Request().Query("sort", "id"),
		Dir:    strings.ToLower(ctx.Request().Query("dir", "desc")),
	}
	if query.Dir != "asc" {
		query.Dir = "desc"
	}
	columns := exportColumns(manifest)
	var rows []map[string]any
	var q orm.Query
	switch resourceName {
	case "users":
		var records []models.User
		q = facades.Orm().Query()
		q = applyResourceSearch(q, query.Search, "name", "email")
		if query.Status != "" {
			q = q.Where("status = ?", query.Status)
		}
		query.Sort = allowedSort(query.Sort, map[string]bool{"id": true, "name": true, "email": true, "status": true}, "id")
		q = q.OrderBy(query.Sort, query.Dir)
		if err := q.Get(&records); err != nil {
			return ctx.Response().Status(500).Json(http.Json{"code": "INTERNAL_ERROR"})
		}
		for _, record := range records {
			rows = append(rows, record.Public())
		}
	case "roles":
		var records []models.Role
		q = facades.Orm().Query()
		q = applyResourceSearch(q, query.Search, "name", "display_name")
		query.Sort = allowedSort(query.Sort, map[string]bool{"id": true, "name": true, "display_name": true}, "id")
		q = q.OrderBy(query.Sort, query.Dir)
		if err := q.Get(&records); err != nil {
			return ctx.Response().Status(500).Json(http.Json{"code": "INTERNAL_ERROR"})
		}
		for _, record := range records {
			rows = append(rows, map[string]any{"id": record.ID, "name": record.Name, "display_name": record.DisplayName})
		}
	case "permissions":
		var records []models.Permission
		q = facades.Orm().Query()
		q = applyResourceSearch(q, query.Search, "name", "display_name")
		query.Sort = allowedSort(query.Sort, map[string]bool{"id": true, "name": true, "display_name": true}, "id")
		q = q.OrderBy(query.Sort, query.Dir)
		if err := q.Get(&records); err != nil {
			return ctx.Response().Status(500).Json(http.Json{"code": "INTERNAL_ERROR"})
		}
		for _, record := range records {
			rows = append(rows, map[string]any{"id": record.ID, "name": record.Name, "display_name": record.DisplayName})
		}
	default:
		q = facades.Orm().Query().Table(manifest.Table)
		q = applyResourceSearch(q, query.Search, fieldNames(resourceSearchFields(manifest))...)
		for _, field := range resourceFilterFields(manifest) {
			value := strings.TrimSpace(ctx.Request().Query(field.Name))
			if value != "" && value != "all" && resourceFilterValueAllowed(field, value) {
				q = q.Where(field.Name+" = ?", value)
			}
		}
		allowed := map[string]bool{"id": true}
		for _, column := range columns[1:] {
			allowed[column.Name] = true
		}
		query.Sort = allowedSort(query.Sort, allowed, "id")
		q = q.OrderBy(query.Sort, query.Dir)
		if err := q.Get(&rows); err != nil {
			return ctx.Response().Status(500).Json(http.Json{"code": "INTERNAL_ERROR"})
		}
	}

	body, err := renderResourceCSV(columns, rows)
	if err != nil {
		return ctx.Response().Status(500).Json(http.Json{"code": "INTERNAL_ERROR"})
	}
	filename := resourceName + "-export.csv"
	return ctx.Response().Header("Content-Type", "text/csv; charset=utf-8").Header("Content-Disposition", "attachment; filename=\""+filename+"\"").String(200, body)
}

func exportColumns(manifest resource.Manifest) []resource.Column {
	columns := []resource.Column{{Name: "id", Label: "ID"}}
	for _, column := range manifest.Columns {
		name := strings.ToLower(column.Name)
		if strings.Contains(name, "password") || strings.Contains(name, "secret") || strings.Contains(name, "token") || strings.Contains(name, "credential") {
			continue
		}
		columns = append(columns, column)
	}
	return columns
}

func csvCell(value any) string {
	if value == nil {
		return ""
	}
	return fmt.Sprint(value)
}

func renderResourceCSV(columns []resource.Column, rows []map[string]any) (string, error) {
	var buffer bytes.Buffer
	writer := csv.NewWriter(&buffer)
	header := make([]string, 0, len(columns))
	for _, column := range columns {
		header = append(header, column.Label)
	}
	if err := writer.Write(header); err != nil {
		return "", err
	}
	for _, row := range rows {
		values := make([]string, 0, len(columns))
		for _, column := range columns {
			values = append(values, csvCell(row[column.Name]))
		}
		if err := writer.Write(values); err != nil {
			return "", err
		}
	}
	writer.Flush()
	return buffer.String(), writer.Error()
}
