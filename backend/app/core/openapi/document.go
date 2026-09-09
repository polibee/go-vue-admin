package openapi

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// BuildDocument builds the versioned OpenAPI contract from CoreEndpoints.
// CoreEndpoints is intentionally kept next to the route boundary so a route
// cannot be added to the platform without making the contract decision visible.
func BuildDocument() map[string]any {
	paths := make(map[string]any)
	for _, endpoint := range CoreEndpoints {
		path := strings.ReplaceAll(endpoint.Path, ":id", "{id}")
		item, ok := paths[path].(map[string]any)
		if !ok {
			item = make(map[string]any)
			paths[path] = item
		}
		item[strings.ToLower(endpoint.Method)] = operationFor(endpoint, path)
	}

	components := map[string]any{"schemas": make(map[string]any)}
	for name, schema := range SchemaDocuments() {
		components["schemas"].(map[string]any)[schemaName(name)] = schema
	}
	return map[string]any{
		"openapi": "3.1.0",
		"info": map[string]any{
			"title":       "Go Vue Admin API",
			"version":     "0.1.0",
			"description": "Generated platform contract for the Go Vue Admin modular monolith.",
		},
		"servers":    []any{map[string]any{"url": "/"}},
		"tags":       []any{map[string]any{"name": "Resources", "description": "Generic resource CRUD endpoints."}},
		"paths":      paths,
		"components": components,
	}
}

// SchemaDocuments returns the schema files committed under contracts/schemas.
// The map keys are stable filenames without the .json suffix.
func SchemaDocuments() map[string]any {
	return map[string]any{
		"demo-resource": resourceSchema(map[string]any{
			"id": map[string]any{"type": "string"}, "name": map[string]any{"type": "string"},
			"status": map[string]any{"type": "string", "enum": []any{"draft", "active"}},
			"owner":  map[string]any{"type": "string"},
		}, []string{"id", "name", "status", "owner"}),
		"user-resource": resourceSchema(map[string]any{
			"id": map[string]any{"type": "string"}, "email": map[string]any{"type": "string", "format": "email"},
			"name": map[string]any{"type": "string"}, "active": map[string]any{"type": "boolean"},
			"role_ids": map[string]any{"type": "array", "items": map[string]any{"type": "string"}},
		}, []string{"id", "email", "name", "active", "role_ids"}),
		"role-resource": resourceSchema(map[string]any{
			"id": map[string]any{"type": "string"}, "name": map[string]any{"type": "string"},
			"description": map[string]any{"type": "string"},
			"permissions": map[string]any{"type": "array", "items": map[string]any{"type": "string"}},
		}, []string{"id", "name", "description", "permissions"}),
		"permission-resource": resourceSchema(map[string]any{
			"id": map[string]any{"type": "string"}, "name": map[string]any{"type": "string"},
			"description": map[string]any{"type": "string"},
		}, []string{"id", "name", "description"}),
		"resource-envelope": map[string]any{
			"type": "object", "required": []any{"data", "meta"},
			"properties": map[string]any{
				"data": map[string]any{},
				"meta": map[string]any{"$ref": "#/components/schemas/ResourceMeta"},
			},
		},
		"resource-list-envelope": map[string]any{
			"type": "object", "required": []any{"data", "meta"},
			"properties": map[string]any{
				"data": map[string]any{"type": "array", "items": map[string]any{}},
				"meta": map[string]any{"$ref": "#/components/schemas/ResourceMeta"},
			},
		},
		"resource-mutation": resourceSchema(map[string]any{"deleted": map[string]any{"type": "boolean"}}, []string{"deleted"}),
		"bulk-delete-request": resourceSchema(map[string]any{
			"ids": map[string]any{"type": "array", "items": map[string]any{"type": "string"}},
		}, []string{"ids"}),
		"resource-meta": map[string]any{
			"type": "object", "properties": map[string]any{
				"request_id": map[string]any{"type": "string"},
				"pagination": map[string]any{"$ref": "#/components/schemas/PaginationMeta"},
			},
		},
		"pagination-meta": map[string]any{
			"type": "object", "required": []any{"page", "per_page", "total", "total_pages"},
			"properties": map[string]any{
				"page":        map[string]any{"type": "integer", "minimum": 1},
				"per_page":    map[string]any{"type": "integer", "minimum": 1},
				"total":       map[string]any{"type": "integer", "minimum": 0},
				"total_pages": map[string]any{"type": "integer", "minimum": 0},
			},
		},
	}
}

func operationFor(endpoint CoreEndpoint, normalizedPath string) map[string]any {
	operationID := operationIDFor(endpoint, normalizedPath)
	operation := map[string]any{
		"operationId": operationID,
		"description": endpoint.Description,
		"responses": map[string]any{
			"200": map[string]any{"description": "Successful response."},
		},
	}
	if strings.HasPrefix(normalizedPath, "/api/resources/") {
		resourceName, item, bulk := resourcePathParts(normalizedPath)
		operation["tags"] = []any{"Resources"}
		if endpoint.Method == "GET" && !item {
			operation["parameters"] = resourceQueryParameters()
			operation["responses"] = map[string]any{"200": responseSchema(resourceName, true)}
		} else if endpoint.Method == "GET" && item {
			operation["parameters"] = []any{pathIDParameter()}
			operation["responses"] = map[string]any{"200": responseSchema(resourceName, false)}
		} else if bulk {
			operation["requestBody"] = jsonRequestBody(map[string]any{"$ref": "#/components/schemas/BulkDeleteRequest"})
			operation["responses"] = map[string]any{"200": responseSchema("ResourceMutation", false)}
		} else {
			if item {
				operation["parameters"] = []any{pathIDParameter()}
			}
			operation["requestBody"] = jsonRequestBody(map[string]any{"$ref": "#/components/schemas/" + schemaTitle(resourceName)})
			if endpoint.Method == "POST" {
				operation["responses"] = map[string]any{"201": responseSchema(resourceName, false)}
			} else if endpoint.Method == "DELETE" {
				operation["responses"] = map[string]any{"200": responseSchema("ResourceMutation", false)}
			} else {
				operation["responses"] = map[string]any{"200": responseSchema(resourceName, false)}
			}
		}
	}
	return operation
}

func resourcePathParts(path string) (string, bool, bool) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	name := parts[2]
	return name, strings.Contains(path, "/{id}"), strings.HasSuffix(path, "/bulk-delete")
}

func operationIDFor(endpoint CoreEndpoint, path string) string {
	if !strings.HasPrefix(path, "/api/resources/") {
		parts := strings.Split(strings.Trim(path, "/"), "/")
		name := ""
		for _, part := range parts {
			if part == "" {
				continue
			}
			name += strings.Title(part)
		}
		return strings.ToLower(endpoint.Method) + name
	}
	name, item, bulk := resourcePathParts(path)
	prefix := strings.TrimSuffix(schemaTitle(name), "Resource")
	isDemo := name == "demo"
	plural := strings.Title(name)
	switch {
	case bulk:
		if isDemo {
			return "bulkDeleteDemoResources"
		}
		return "bulkDelete" + plural
	case endpoint.Method == "GET" && !item:
		if isDemo {
			return "listDemoResources"
		}
		return "list" + plural
	case endpoint.Method == "GET":
		if isDemo {
			return "getDemoResource"
		}
		return "get" + prefix
	case endpoint.Method == "POST":
		if isDemo {
			return "createDemoResource"
		}
		return "create" + prefix
	case endpoint.Method == "PUT":
		if isDemo {
			return "updateDemoResource"
		}
		return "update" + prefix
	default:
		if isDemo {
			return "deleteDemoResource"
		}
		return "delete" + prefix
	}
}

func resourceSchema(properties map[string]any, required []string) map[string]any {
	requiredAny := make([]any, len(required))
	for i, value := range required {
		requiredAny[i] = value
	}
	return map[string]any{"type": "object", "required": requiredAny, "properties": properties, "additionalProperties": false}
}

func resourceQueryParameters() []any {
	return []any{
		map[string]any{"name": "page", "in": "query", "schema": map[string]any{"type": "integer", "minimum": 1}},
		map[string]any{"name": "per_page", "in": "query", "schema": map[string]any{"type": "integer", "minimum": 1}},
		map[string]any{"name": "search", "in": "query", "schema": map[string]any{"type": "string"}},
		map[string]any{"name": "sort", "in": "query", "schema": map[string]any{"type": "string"}},
		map[string]any{"name": "sort_dir", "in": "query", "schema": map[string]any{"type": "string", "enum": []any{"asc", "desc"}}},
		map[string]any{"name": "filter", "in": "query", "style": "deepObject", "explode": true, "schema": map[string]any{"type": "object", "additionalProperties": map[string]any{"type": "string"}}},
	}
}

func pathIDParameter() map[string]any {
	return map[string]any{"name": "id", "in": "path", "required": true, "schema": map[string]any{"type": "string"}}
}

func jsonRequestBody(schema map[string]any) map[string]any {
	return map[string]any{"required": true, "content": map[string]any{"application/json": map[string]any{"schema": schema}}}
}

func responseSchema(resourceName string, list bool) map[string]any {
	if list {
		return map[string]any{"description": "Successful response.", "content": map[string]any{"application/json": map[string]any{"schema": map[string]any{"$ref": "#/components/schemas/ResourceListEnvelope"}}}}
	}
	return map[string]any{"description": "Successful response.", "content": map[string]any{"application/json": map[string]any{"schema": map[string]any{"$ref": "#/components/schemas/ResourceEnvelope"}}}}
}

func schemaTitle(value string) string {
	if value == "demo" {
		return "DemoResource"
	}
	return strings.Title(strings.TrimSuffix(value, "s")) + "Resource"
}

func schemaName(value string) string {
	return strings.ReplaceAll(strings.Title(strings.ReplaceAll(value, "-", " ")), " ", "")
}

// WriteJSON writes an indented, newline-terminated JSON document.
func WriteJSON(path string, value any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}
