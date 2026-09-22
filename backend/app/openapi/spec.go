package openapi

func Spec() map[string]any {
	return map[string]any{
		"openapi":  "3.0.3",
		"info":     map[string]any{"title": "Go Vue Admin API", "version": "0.1.0", "description": "The checked-in API contract for the admin backend."},
		"servers":  []map[string]any{{"url": "/api/v1"}},
		"security": []map[string]any{{"bearerAuth": []any{}}},
		"components": map[string]any{
			"securitySchemes": map[string]any{"bearerAuth": map[string]any{"type": "http", "scheme": "bearer", "bearerFormat": "JWT"}},
			"schemas": map[string]any{
				"UserStatus":              map[string]any{"type": "string", "enum": []string{"active", "disabled", "locked"}},
				"DataScope":               map[string]any{"type": "string", "enum": []string{"all", "own"}},
				"FieldPermissionOverride": map[string]any{"type": "object", "required": []string{"readable", "writable"}, "properties": map[string]any{"readable": map[string]any{"type": "boolean"}, "writable": map[string]any{"type": "boolean"}}},
				"Error":                   map[string]any{"type": "object", "required": []string{"code"}, "properties": map[string]any{"code": map[string]any{"type": "string"}}},
				"ResourceOption":          resourceOptionSchema(),
				"ResourceField":           resourceFieldSchema(),
				"ResourceFilter":          resourceFilterSchema(),
				"ResourceColumn":          resourceColumnSchema(),
				"ResourceAction":          resourceActionSchema(),
				"ResourceRelation":        resourceRelationSchema(),
				"ResourceFormGroup":       resourceFormGroupSchema(),
				"ResourceDetailSection":   resourceDetailSectionSchema(),
				"ResourceFieldDependency": resourceFieldDependencySchema(),
				"RelationOption":          relationOptionSchema(),
				"ResourceManifest":        resourceManifestSchema(),
				"ActionRequest":           actionRequestSchema(),
				"ActionResponse":          actionResponseSchema(),
				"GlobalSearchResult":      globalSearchResultSchema(),
				"GlobalSearchResponse":    globalSearchResponseSchema(),
			},
		},
		"paths": map[string]any{
			"/auth/login": map[string]any{"post": map[string]any{
				"security": []any{}, "operationId": "login", "requestBody": jsonBody("LoginRequest", map[string]any{"type": "object", "required": []string{"email", "password"}, "properties": map[string]any{"email": map[string]any{"type": "string", "format": "email"}, "password": map[string]any{"type": "string", "format": "password"}}}),
				"responses": map[string]any{"200": jsonResponse("LoginResponse"), "403": errorResponse(), "429": errorResponse(), "503": errorResponse()},
			}},
			"/auth/refresh": map[string]any{"post": map[string]any{
				"security": []any{}, "operationId": "refresh", "responses": map[string]any{"200": jsonResponse("RefreshResponse"), "401": errorResponse()},
			}},
			"/auth/logout-all": map[string]any{"post": operation("logoutAll")},
			"/auth/me":         map[string]any{"get": operation("currentUser")},
			"/admin/registry":  map[string]any{"get": operation("listResources")},
			"/admin/search":    map[string]any{"get": globalSearchOperation()},
			"/admin/{resource}": map[string]any{
				"get":  listOperation("listResourceRows", []map[string]any{pathParameter("resource"), queryParameter("page", "integer"), queryParameter("per_page", "integer"), queryParameter("search", "string"), queryParameter("sort", "string"), queryParameter("dir", "string"), queryParameter("trashed", "string")}),
				"post": resourceWriteOperation("createResource", "201"),
			},
			"/admin/{resource}/export":                       exportOperation(),
			"/admin/{resource}/actions/{action}":             map[string]any{"post": resourceActionOperation()},
			"/admin/{resource}/relations/{relation}/options": map[string]any{"get": relationOptionsOperation()},
			"/admin/{resource}/{id}/relations/{relation}":    map[string]any{"get": relationRecordsOperation()},
			"/admin/{resource}/{id}": map[string]any{
				"get":    resourceItemOperation("showResource"),
				"put":    resourceWriteOperation("updateResource", "200", true),
				"delete": map[string]any{"operationId": "deleteResource", "parameters": []map[string]any{pathParameter("resource"), pathParameter("id")}, "responses": map[string]any{"204": map[string]any{"description": "Resource deleted"}, "401": errorResponse(), "403": errorResponse(), "404": errorResponse()}},
			},
			"/admin/overview":       map[string]any{"get": operation("adminOverview")},
			"/admin/audit-logs":     map[string]any{"get": listOperation("auditLogs", []map[string]any{queryParameter("page", "integer"), queryParameter("per_page", "integer"), queryParameter("action", "string"), queryParameter("user_id", "integer")})},
			"/admin/audit-logs/cleanup": map[string]any{"post": map[string]any{"operationId": "cleanupAuditLogs", "requestBody": jsonBody("AuditCleanupRequest", map[string]any{"type": "object", "required": []string{"retention_days"}, "properties": map[string]any{"retention_days": map[string]any{"type": "integer", "minimum": 1, "maximum": 3650}}}), "responses": map[string]any{"200": jsonResponse("AuditCleanupResponse"), "401": errorResponse(), "403": errorResponse(), "422": errorResponse()}}},
			"/admin/settings":       map[string]any{"get": operation("systemSettings")},
			"/admin/settings/{key}": map[string]any{"put": map[string]any{"operationId": "updateSystemSetting", "parameters": []map[string]any{{"name": "key", "in": "path", "required": true, "schema": map[string]any{"type": "string"}}}, "requestBody": jsonBody("SystemSettingRequest", map[string]any{"type": "object", "required": []string{"value"}, "properties": map[string]any{"value": map[string]any{"type": "string"}, "value_type": map[string]any{"type": "string", "enum": []string{"string", "boolean", "integer", "json"}}, "group": map[string]any{"type": "string"}, "description": map[string]any{"type": "string"}}}), "responses": map[string]any{"200": jsonResponse("SystemSettingResponse"), "422": errorResponse()}}},
			"/admin/roles/{id}/permissions": map[string]any{"put": map[string]any{
				"operationId": "replaceRolePermissions",
				"parameters":  []map[string]any{pathParameter("id")},
				"requestBody": jsonBody("RolePermissionsRequest", map[string]any{"type": "object", "required": []string{"permission_ids"}, "properties": map[string]any{
					"permission_ids": map[string]any{"type": "array", "items": map[string]any{"type": "integer", "format": "int64"}},
					"scopes":         map[string]any{"type": "object", "additionalProperties": map[string]any{"$ref": "#/components/schemas/DataScope"}},
					"fields":         map[string]any{"type": "object", "additionalProperties": map[string]any{"type": "object", "additionalProperties": map[string]any{"$ref": "#/components/schemas/FieldPermissionOverride"}}},
				}}),
				"responses": map[string]any{"204": map[string]any{"description": "Permissions replaced"}, "401": errorResponse(), "403": errorResponse(), "404": errorResponse(), "422": errorResponse()},
			}},
		},
	}
}

func resourceOptionSchema() map[string]any {
	return map[string]any{"type": "object", "required": []string{"value", "label"}, "properties": map[string]any{"value": map[string]any{"type": "string"}, "label": map[string]any{"type": "string"}}}
}

func resourceFieldSchema() map[string]any {
	return map[string]any{"type": "object", "required": []string{"name", "label", "type", "visible", "readable", "writable", "sensitive"}, "properties": map[string]any{"name": map[string]any{"type": "string"}, "label": map[string]any{"type": "string"}, "type": map[string]any{"type": "string"}, "required": map[string]any{"type": "boolean"}, "visible": map[string]any{"type": "boolean"}, "readable": map[string]any{"type": "boolean"}, "writable": map[string]any{"type": "boolean"}, "sensitive": map[string]any{"type": "boolean"}, "options": map[string]any{"type": "array", "items": map[string]any{"$ref": "#/components/schemas/ResourceOption"}}}}
}

func resourceColumnSchema() map[string]any {
	return map[string]any{"type": "object", "required": []string{"name", "label", "sortable"}, "properties": map[string]any{"name": map[string]any{"type": "string"}, "label": map[string]any{"type": "string"}, "sortable": map[string]any{"type": "boolean"}}}
}

func resourceFilterSchema() map[string]any {
	return map[string]any{"type": "object", "required": []string{"name", "label", "type"}, "properties": map[string]any{"name": map[string]any{"type": "string"}, "label": map[string]any{"type": "string"}, "type": map[string]any{"type": "string", "enum": []string{"select", "multi-select", "boolean", "text", "date-range", "relation"}}, "options": map[string]any{"type": "array", "items": map[string]any{"$ref": "#/components/schemas/ResourceOption"}}, "relation": map[string]any{"type": "string"}}}
}

func resourceActionSchema() map[string]any {
	return map[string]any{"type": "object", "required": []string{"name", "label", "permission", "batch"}, "properties": map[string]any{"name": map[string]any{"type": "string"}, "label": map[string]any{"type": "string"}, "kind": map[string]any{"type": "string"}, "permission": map[string]any{"type": "string"}, "batch": map[string]any{"type": "boolean"}, "payload": map[string]any{"type": "string"}}}
}

func resourceManifestSchema() map[string]any {
	return map[string]any{"type": "object", "required": []string{"name", "label", "route", "permissions", "navigation", "fields", "columns"}, "properties": map[string]any{"name": map[string]any{"type": "string"}, "label": map[string]any{"type": "string"}, "route": map[string]any{"type": "string"}, "permissions": map[string]any{"type": "array", "items": map[string]any{"type": "string"}}, "navigation": map[string]any{"type": "object", "required": []string{"group", "order"}, "properties": map[string]any{"group": map[string]any{"type": "string"}, "order": map[string]any{"type": "integer"}, "hidden": map[string]any{"type": "boolean"}}}, "data_scope": map[string]any{"$ref": "#/components/schemas/DataScope"}, "owner_field": map[string]any{"type": "string"}, "soft_delete": map[string]any{"type": "boolean"}, "fields": map[string]any{"type": "array", "items": map[string]any{"$ref": "#/components/schemas/ResourceField"}}, "columns": map[string]any{"type": "array", "items": map[string]any{"$ref": "#/components/schemas/ResourceColumn"}}, "actions": map[string]any{"type": "array", "items": map[string]any{"$ref": "#/components/schemas/ResourceAction"}}, "filters": map[string]any{"type": "array", "items": map[string]any{"$ref": "#/components/schemas/ResourceFilter"}}, "relations": map[string]any{"type": "array", "items": map[string]any{"$ref": "#/components/schemas/ResourceRelation"}}, "form_groups": map[string]any{"type": "array", "items": map[string]any{"$ref": "#/components/schemas/ResourceFormGroup"}}, "details": map[string]any{"type": "array", "items": map[string]any{"$ref": "#/components/schemas/ResourceDetailSection"}}, "dependencies": map[string]any{"type": "array", "items": map[string]any{"$ref": "#/components/schemas/ResourceFieldDependency"}}}}
}

func resourceRelationSchema() map[string]any {
	return map[string]any{"type": "object", "required": []string{"name", "kind", "resource", "field", "foreign_field", "label_field", "selectable", "multiple"}, "properties": map[string]any{"name": map[string]any{"type": "string"}, "kind": map[string]any{"type": "string", "enum": []string{"belongsTo", "hasMany"}}, "resource": map[string]any{"type": "string"}, "field": map[string]any{"type": "string"}, "foreign_field": map[string]any{"type": "string"}, "label_field": map[string]any{"type": "string"}, "selectable": map[string]any{"type": "boolean"}, "multiple": map[string]any{"type": "boolean"}, "permission": map[string]any{"type": "string"}, "filter_fields": map[string]any{"type": "array", "items": map[string]any{"type": "string"}}}}
}

func resourceFormGroupSchema() map[string]any {
	return map[string]any{"type": "object", "required": []string{"name", "label", "fields"}, "properties": map[string]any{"name": map[string]any{"type": "string"}, "label": map[string]any{"type": "string"}, "columns": map[string]any{"type": "integer", "minimum": 1, "maximum": 4}, "fields": map[string]any{"type": "array", "items": map[string]any{"type": "string"}}}}
}

func resourceDetailSectionSchema() map[string]any {
	return map[string]any{"type": "object", "required": []string{"name", "label", "fields"}, "properties": map[string]any{"name": map[string]any{"type": "string"}, "label": map[string]any{"type": "string"}, "fields": map[string]any{"type": "array", "items": map[string]any{"type": "string"}}}}
}

func resourceFieldDependencySchema() map[string]any {
	return map[string]any{"type": "object", "required": []string{"field", "on", "value"}, "properties": map[string]any{"field": map[string]any{"type": "string"}, "on": map[string]any{"type": "string"}, "value": map[string]any{"type": "string"}}}
}

func relationOptionSchema() map[string]any {
	return map[string]any{"type": "object", "required": []string{"value", "label"}, "properties": map[string]any{"value": map[string]any{"type": "string"}, "label": map[string]any{"type": "string"}}}
}

func actionRequestSchema() map[string]any {
	return map[string]any{"type": "object", "properties": map[string]any{
		"ids":     map[string]any{"type": "array", "minItems": 1, "maxItems": 1000, "items": map[string]any{"type": "integer", "format": "int64"}},
		"selection": map[string]any{"type": "object", "required": []string{"mode"}, "properties": map[string]any{"mode": map[string]any{"type": "string", "enum": []string{"ids", "query"}}, "ids": map[string]any{"type": "array", "items": map[string]any{"type": "integer", "format": "int64"}}, "query": map[string]any{"type": "object", "additionalProperties": map[string]any{"type": "string"}}, "exclude_ids": map[string]any{"type": "array", "items": map[string]any{"type": "integer", "format": "int64"}}}},
		"payload": map[string]any{"type": "object", "additionalProperties": true},
	}}
}

func actionResponseSchema() map[string]any {
	return map[string]any{"type": "object", "required": []string{"action", "requested", "succeeded", "failed", "skipped", "failures", "skips"}, "properties": map[string]any{
		"action": map[string]any{"type": "string"}, "requested": map[string]any{"type": "integer"}, "succeeded": map[string]any{"type": "integer"}, "failed": map[string]any{"type": "integer"}, "skipped": map[string]any{"type": "integer"},
		"failures": map[string]any{"type": "array", "items": map[string]any{"type": "object", "required": []string{"id", "code"}, "properties": map[string]any{"id": map[string]any{"type": "integer", "format": "int64"}, "code": map[string]any{"type": "string"}}}},
		"skips":    map[string]any{"type": "array", "items": map[string]any{"type": "object", "required": []string{"id", "code"}, "properties": map[string]any{"id": map[string]any{"type": "integer", "format": "int64"}, "code": map[string]any{"type": "string"}}}},
	}}
}

func globalSearchResultSchema() map[string]any {
	return map[string]any{"type": "object", "required": []string{"resource", "label", "id", "title", "route"}, "properties": map[string]any{
		"resource": map[string]any{"type": "string"}, "label": map[string]any{"type": "string"}, "id": map[string]any{"oneOf": []map[string]any{{"type": "string"}, {"type": "integer"}}}, "title": map[string]any{"type": "string"}, "subtitle": map[string]any{"type": "string"}, "route": map[string]any{"type": "string"},
	}}
}

func globalSearchResponseSchema() map[string]any {
	return map[string]any{"type": "object", "required": []string{"data"}, "properties": map[string]any{
		"data": map[string]any{"type": "array", "items": map[string]any{"$ref": "#/components/schemas/GlobalSearchResult"}},
	}}
}

func jsonBody(name string, schema map[string]any) map[string]any {
	return map[string]any{"required": true, "content": map[string]any{"application/json": map[string]any{"schema": schema, "x-schema-name": name}}}
}

func jsonResponse(name string) map[string]any {
	return map[string]any{"description": "Successful response", "content": map[string]any{"application/json": map[string]any{"schema": map[string]any{"type": "object"}, "x-schema-name": name}}}
}

func errorResponse() map[string]any {
	return map[string]any{"description": "Error response", "content": map[string]any{"application/json": map[string]any{"schema": map[string]any{"$ref": "#/components/schemas/Error"}}}}
}

func operation(operationID string) map[string]any {
	return map[string]any{"operationId": operationID, "responses": map[string]any{"200": jsonResponse(operationID + "Response"), "401": errorResponse(), "403": errorResponse()}}
}

func listOperation(operationID string, parameters []map[string]any) map[string]any {
	result := operation(operationID)
	result["parameters"] = parameters
	return result
}

func globalSearchOperation() map[string]any {
	return map[string]any{
		"operationId": "globalSearch",
		"parameters":  []map[string]any{queryParameter("q", "string")},
		"responses": map[string]any{
			"200": map[string]any{"description": "Search results", "content": map[string]any{"application/json": map[string]any{"schema": map[string]any{"$ref": "#/components/schemas/GlobalSearchResponse"}}}},
			"401": errorResponse(), "403": errorResponse(),
		},
	}
}

func exportOperation() map[string]any {
	return map[string]any{
		"operationId": "exportResource",
		"parameters":  []map[string]any{pathParameter("resource"), queryParameter("search", "string"), queryParameter("status", "string"), queryParameter("sort", "string"), queryParameter("dir", "string")},
		"responses": map[string]any{
			"200": map[string]any{"description": "CSV export", "content": map[string]any{"text/csv": map[string]any{"schema": map[string]any{"type": "string", "format": "binary"}}}},
			"401": errorResponse(), "403": errorResponse(), "404": errorResponse(),
		},
	}
}

func resourceActionOperation() map[string]any {
	return map[string]any{
		"operationId": "executeResourceAction",
		"parameters":  []map[string]any{pathParameter("resource"), pathParameter("action")},
		"requestBody": jsonBody("ActionRequest", actionRequestSchema()),
		"responses": map[string]any{
			"200": map[string]any{"description": "Action executed", "content": map[string]any{"application/json": map[string]any{"schema": map[string]any{"type": "object", "properties": map[string]any{"data": map[string]any{"$ref": "#/components/schemas/ActionResponse"}}}}}},
			"401": errorResponse(), "403": errorResponse(), "404": errorResponse(), "422": errorResponse(),
		},
	}
}

func relationOptionsOperation() map[string]any {
	return map[string]any{"operationId": "resourceRelationOptions", "parameters": []map[string]any{pathParameter("resource"), pathParameter("relation")}, "responses": map[string]any{"200": map[string]any{"description": "Relation options", "content": map[string]any{"application/json": map[string]any{"schema": map[string]any{"type": "object", "properties": map[string]any{"data": map[string]any{"type": "array", "items": map[string]any{"$ref": "#/components/schemas/RelationOption"}}}}}}}, "401": errorResponse(), "403": errorResponse(), "404": errorResponse()}}
}

func relationRecordsOperation() map[string]any {
	return map[string]any{"operationId": "resourceRelationRecords", "parameters": []map[string]any{pathParameter("resource"), pathParameter("id"), pathParameter("relation")}, "responses": map[string]any{"200": map[string]any{"description": "Related records", "content": map[string]any{"application/json": map[string]any{"schema": map[string]any{"type": "object", "properties": map[string]any{"data": map[string]any{"type": "array", "items": map[string]any{"$ref": "#/components/schemas/RelationOption"}}}}}}}, "401": errorResponse(), "403": errorResponse(), "404": errorResponse()}}
}

func queryParameter(name string, valueType string) map[string]any {
	return map[string]any{"name": name, "in": "query", "required": false, "schema": map[string]any{"type": valueType}}
}

func pathParameter(name string) map[string]any {
	return map[string]any{"name": name, "in": "path", "required": true, "schema": map[string]any{"type": "string"}}
}

func resourceItemOperation(operationID string) map[string]any {
	return map[string]any{"operationId": operationID, "parameters": []map[string]any{pathParameter("resource"), pathParameter("id")}, "responses": map[string]any{"200": jsonResponse(operationID + "Response"), "401": errorResponse(), "403": errorResponse(), "404": errorResponse()}}
}

func resourceWriteOperation(operationID, success string, item ...bool) map[string]any {
	parameters := []map[string]any{pathParameter("resource")}
	if len(item) > 0 && item[0] {
		parameters = append(parameters, pathParameter("id"))
	}
	return map[string]any{"operationId": operationID, "parameters": parameters, "requestBody": jsonBody("ResourceRecord", map[string]any{"type": "object", "additionalProperties": true}), "responses": map[string]any{success: jsonResponse(operationID + "Response"), "401": errorResponse(), "403": errorResponse(), "404": errorResponse(), "422": errorResponse()}}
}
