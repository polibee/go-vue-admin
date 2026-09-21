package controllers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/goravel/framework/contracts/http"

	"goravel/app/core/resource"
	"goravel/app/facades"
	"goravel/app/modules/admin/registry"
)

func (r *ResourceController) Create(ctx http.Context) http.Response {
	manifest, ok := generatedManifest(ctx)
	if !ok {
		return ctx.Response().Status(404).Json(http.Json{"code": "RESOURCE_NOT_FOUND"})
	}
	values, err := bindGeneratedValues(ctx, manifest)
	if err != nil {
		return ctx.Response().Status(422).Json(http.Json{"code": "VALIDATION_ERROR"})
	}
	if err := facades.Orm().Query().Table(manifest.Table).Create(&values); err != nil {
		return ctx.Response().Status(500).Json(http.Json{"code": "INTERNAL_ERROR"})
	}
	return ctx.Response().Status(201).Json(http.Json{"data": values})
}

func (r *ResourceController) Update(ctx http.Context) http.Response {
	manifest, ok := generatedManifest(ctx)
	if !ok {
		return ctx.Response().Status(404).Json(http.Json{"code": "RESOURCE_NOT_FOUND"})
	}
	id, err := strconv.ParseInt(ctx.Request().Route("id"), 10, 64)
	if err != nil || id < 1 {
		return ctx.Response().Status(404).Json(http.Json{"code": "RESOURCE_NOT_FOUND"})
	}
	values, err := bindGeneratedValues(ctx, manifest)
	if err != nil {
		return ctx.Response().Status(422).Json(http.Json{"code": "VALIDATION_ERROR"})
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
	id, err := strconv.ParseInt(ctx.Request().Route("id"), 10, 64)
	if err != nil || id < 1 {
		return ctx.Response().Status(404).Json(http.Json{"code": "RESOURCE_NOT_FOUND"})
	}
	if _, err := facades.Orm().Query().Table(manifest.Table).Where("id = ?", id).Delete(); err != nil {
		return ctx.Response().Status(500).Json(http.Json{"code": "INTERNAL_ERROR"})
	}
	return ctx.Response().NoContent(204)
}

func generatedManifest(ctx http.Context) (resource.Manifest, bool) {
	manifest, err := registry.AdminRegistry().Find(ctx.Request().Route("resource"))
	return manifest, err == nil && manifest.Table != ""
}

func bindGeneratedValues(ctx http.Context, manifest resource.Manifest) (map[string]any, error) {
	var payload map[string]any
	if err := ctx.Request().Bind(&payload); err != nil {
		return nil, err
	}
	allowed := make(map[string]bool, len(manifest.Fields))
	for _, field := range manifest.Fields {
		allowed[field.Name] = true
	}
	return validateGeneratedValues(payload, manifest)
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
