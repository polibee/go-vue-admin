package controllers

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/goravel/framework/contracts/http"

	"goravel/app/core/resource"
	"goravel/app/facades"
	"goravel/app/modules/admin/registry"
	rbacservices "goravel/app/services/rbac"
)

const maxImportBytes int64 = 5 * 1024 * 1024

type resourceImportRowError struct {
	Row      int      `json:"row"`
	Messages []string `json:"messages"`
}

type resourceImportResult struct {
	Rows      []map[string]any         `json:"rows"`
	RowErrors []resourceImportRowError `json:"row_errors"`
	Errors    []string                 `json:"errors"`
}

func (r *ResourceController) ImportPreview(ctx http.Context) http.Response {
	manifest, ok := importManifest(ctx)
	if !ok {
		return ctx.Response().Status(404).Json(http.Json{"code": "RESOURCE_NOT_FOUND"})
	}
	if hasSensitiveImportField(manifest) {
		return ctx.Response().Status(422).Json(http.Json{"code": "IMPORT_SENSITIVE_FIELD"})
	}
	result, err := parseImportRequest(ctx, manifest)
	if err != nil {
		return ctx.Response().Status(422).Json(http.Json{"code": err.Error()})
	}
	return ctx.Response().Success().Json(http.Json{"data": result})
}

func (r *ResourceController) Import(ctx http.Context) http.Response {
	manifest, ok := importManifest(ctx)
	if !ok {
		return ctx.Response().Status(404).Json(http.Json{"code": "RESOURCE_NOT_FOUND"})
	}
	if hasSensitiveImportField(manifest) {
		return ctx.Response().Status(422).Json(http.Json{"code": "IMPORT_SENSITIVE_FIELD"})
	}
	result, err := parseImportRequest(ctx, manifest)
	if err != nil {
		return ctx.Response().Status(422).Json(http.Json{"code": err.Error()})
	}
	if len(result.Errors) > 0 || len(result.RowErrors) > 0 {
		return ctx.Response().Status(422).Json(http.Json{"data": result, "code": "IMPORT_VALIDATION_ERROR"})
	}
	if err := persistImportRows(manifest, result.Rows); err != nil {
		return ctx.Response().Status(422).Json(http.Json{"code": err.Error()})
	}
	return ctx.Response().Success().Json(http.Json{"data": http.Json{"imported": len(result.Rows)}})
}

func importManifest(ctx http.Context) (resource.Manifest, bool) {
	manifest, err := registry.AdminRegistry().Find(ctx.Request().Route("resource"))
	return manifest, err == nil && manifest.Table != ""
}

func parseImportRequest(ctx http.Context, manifest resource.Manifest) (resourceImportResult, error) {
	file, err := ctx.Request().File("file")
	if err != nil {
		return resourceImportResult{}, fmt.Errorf("IMPORT_FILE_REQUIRED")
	}
	if size, sizeErr := file.Size(); sizeErr != nil || size > maxImportBytes {
		return resourceImportResult{}, fmt.Errorf("IMPORT_FILE_TOO_LARGE")
	}
	input, err := os.Open(file.File())
	if err != nil {
		return resourceImportResult{}, fmt.Errorf("IMPORT_FILE_INVALID")
	}
	defer input.Close()
	return parseResourceCSV(io.LimitReader(input, maxImportBytes+1), manifest), nil
}

func parseResourceCSV(input io.Reader, manifest resource.Manifest) resourceImportResult {
	reader := csv.NewReader(input)
	reader.FieldsPerRecord = -1
	records, err := reader.ReadAll()
	result := resourceImportResult{}
	if err != nil {
		result.Errors = append(result.Errors, "invalid CSV")
		return result
	}
	if len(records) == 0 || len(records[0]) == 0 {
		result.Errors = append(result.Errors, "CSV header is required")
		return result
	}
	fields := make(map[string]resource.Field, len(manifest.Fields))
	for _, field := range manifest.Fields {
		fields[field.Name] = field
	}
	columns := make([]string, len(records[0]))
	for index, value := range records[0] {
		columns[index] = strings.TrimSpace(value)
		field, ok := fields[columns[index]]
		if !ok || isSensitiveImportField(columns[index]) {
			result.Errors = append(result.Errors, "unsupported column: "+columns[index])
			continue
		}
		if field.Name == "id" {
			result.Errors = append(result.Errors, "unsupported column: "+columns[index])
		}
	}
	if len(result.Errors) > 0 {
		return result
	}
	for index, record := range records[1:] {
		rowNumber := index + 2
		if len(record) != len(columns) {
			result.RowErrors = append(result.RowErrors, resourceImportRowError{Row: rowNumber, Messages: []string{"column count does not match header"}})
			continue
		}
		values := make(map[string]any, len(columns))
		messages := make([]string, 0)
		for columnIndex, name := range columns {
			field := fields[name]
			value, valueErr := importValue(field, strings.TrimSpace(record[columnIndex]))
			if valueErr != nil {
				messages = append(messages, valueErr.Error())
				continue
			}
			values[name] = value
		}
		for _, field := range manifest.Fields {
			value := values[field.Name]
			if field.Required && (value == nil || strings.TrimSpace(fmt.Sprint(value)) == "") {
				messages = append(messages, "required field: "+field.Name)
			}
		}
		if len(messages) > 0 {
			result.RowErrors = append(result.RowErrors, resourceImportRowError{Row: rowNumber, Messages: messages})
			continue
		}
		result.Rows = append(result.Rows, values)
	}
	return result
}

func importValue(field resource.Field, value string) (any, error) {
	if value == "" {
		return "", nil
	}
	switch field.Type {
	case "boolean":
		parsed, err := strconv.ParseBool(value)
		if err != nil {
			return nil, fmt.Errorf("invalid boolean field: %s", field.Name)
		}
		return parsed, nil
	case "integer":
		parsed, err := strconv.Atoi(value)
		if err != nil {
			return nil, fmt.Errorf("invalid integer field: %s", field.Name)
		}
		return parsed, nil
	case "select":
		for _, option := range field.Options {
			if option.Value == value {
				return value, nil
			}
		}
		return nil, fmt.Errorf("invalid option field: %s", field.Name)
	default:
		return value, nil
	}
}

func hasSensitiveImportField(manifest resource.Manifest) bool {
	for _, field := range manifest.Fields {
		if isSensitiveImportField(field.Name) {
			return true
		}
	}
	return false
}

func isSensitiveImportField(name string) bool {
	name = strings.ToLower(name)
	return strings.Contains(name, "password") || strings.Contains(name, "secret") || strings.Contains(name, "token") || strings.Contains(name, "credential")
}

func persistImportRows(manifest resource.Manifest, rows []map[string]any) error {
	for _, row := range rows {
		if manifest.Name == "roles" {
			_, err := rbacservices.NewRoleService().Create(fmt.Sprint(row["name"]), fmt.Sprint(row["display_name"]))
			if err != nil {
				return fmt.Errorf("IMPORT_PERSIST_FAILED")
			}
			continue
		}
		if err := facades.Orm().Query().Table(manifest.Table).Create(&row); err != nil {
			return fmt.Errorf("IMPORT_PERSIST_FAILED")
		}
	}
	return nil
}
