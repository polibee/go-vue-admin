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
				"UserStatus":       map[string]any{"type": "string", "enum": []string{"active", "disabled", "locked"}},
				"Error":            map[string]any{"type": "object", "required": []string{"code"}, "properties": map[string]any{"code": map[string]any{"type": "string"}}},
				"ResourceOption":   resourceOptionSchema(),
				"ResourceField":    resourceFieldSchema(),
				"ResourceColumn":   resourceColumnSchema(),
				"ResourceAction":   resourceActionSchema(),
				"ResourceManifest": resourceManifestSchema(),
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
			"/admin/{resource}": map[string]any{
				"get":  listOperation("listResourceRows", []map[string]any{pathParameter("resource"), queryParameter("page", "integer"), queryParameter("per_page", "integer"), queryParameter("search", "string"), queryParameter("sort", "string"), queryParameter("dir", "string")}),
				"post": resourceWriteOperation("createResource", "201"),
			},
			"/admin/{resource}/export": exportOperation(),
			"/admin/{resource}/{id}": map[string]any{
				"get":    resourceItemOperation("showResource"),
				"put":    resourceWriteOperation("updateResource", "200", true),
				"delete": map[string]any{"operationId": "deleteResource", "parameters": []map[string]any{pathParameter("resource"), pathParameter("id")}, "responses": map[string]any{"204": map[string]any{"description": "Resource deleted"}, "401": errorResponse(), "403": errorResponse(), "404": errorResponse()}},
			},
			"/admin/overview":       map[string]any{"get": operation("adminOverview")},
			"/admin/audit-logs":     map[string]any{"get": listOperation("auditLogs", []map[string]any{queryParameter("page", "integer"), queryParameter("per_page", "integer"), queryParameter("action", "string"), queryParameter("user_id", "integer")})},
			"/admin/settings":       map[string]any{"get": operation("systemSettings")},
			"/admin/settings/{key}": map[string]any{"put": map[string]any{"operationId": "updateSystemSetting", "parameters": []map[string]any{{"name": "key", "in": "path", "required": true, "schema": map[string]any{"type": "string"}}}, "requestBody": jsonBody("SystemSettingRequest", map[string]any{"type": "object", "required": []string{"value"}, "properties": map[string]any{"value": map[string]any{"type": "string"}, "value_type": map[string]any{"type": "string", "enum": []string{"string", "boolean", "integer", "json"}}, "group": map[string]any{"type": "string"}, "description": map[string]any{"type": "string"}}}), "responses": map[string]any{"200": jsonResponse("SystemSettingResponse"), "422": errorResponse()}}},
			"/admin/users/status": map[string]any{"put": map[string]any{
				"operationId": "bulkSetUserStatus", "requestBody": jsonBody("BulkUserStatusRequest", map[string]any{"type": "object", "required": []string{"user_ids", "status"}, "properties": map[string]any{"user_ids": map[string]any{"type": "array", "items": map[string]any{"type": "integer", "format": "int64"}}, "status": map[string]any{"$ref": "#/components/schemas/UserStatus"}}}),
				"responses": map[string]any{"204": map[string]any{"description": "Status updated"}, "409": errorResponse(), "422": errorResponse()},
			}},
		},
	}
}

func resourceOptionSchema() map[string]any {
	return map[string]any{"type": "object", "required": []string{"value", "label"}, "properties": map[string]any{"value": map[string]any{"type": "string"}, "label": map[string]any{"type": "string"}}}
}

func resourceFieldSchema() map[string]any {
	return map[string]any{"type": "object", "required": []string{"name", "label", "type"}, "properties": map[string]any{"name": map[string]any{"type": "string"}, "label": map[string]any{"type": "string"}, "type": map[string]any{"type": "string"}, "required": map[string]any{"type": "boolean"}, "options": map[string]any{"type": "array", "items": map[string]any{"$ref": "#/components/schemas/ResourceOption"}}}}
}

func resourceColumnSchema() map[string]any {
	return map[string]any{"type": "object", "required": []string{"name", "label", "sortable"}, "properties": map[string]any{"name": map[string]any{"type": "string"}, "label": map[string]any{"type": "string"}, "sortable": map[string]any{"type": "boolean"}}}
}

func resourceActionSchema() map[string]any {
	return map[string]any{"type": "object", "required": []string{"name", "label", "permission"}, "properties": map[string]any{"name": map[string]any{"type": "string"}, "label": map[string]any{"type": "string"}, "kind": map[string]any{"type": "string"}, "permission": map[string]any{"type": "string"}}}
}

func resourceManifestSchema() map[string]any {
	return map[string]any{"type": "object", "required": []string{"name", "label", "route", "fields", "columns"}, "properties": map[string]any{"name": map[string]any{"type": "string"}, "label": map[string]any{"type": "string"}, "route": map[string]any{"type": "string"}, "fields": map[string]any{"type": "array", "items": map[string]any{"$ref": "#/components/schemas/ResourceField"}}, "columns": map[string]any{"type": "array", "items": map[string]any{"$ref": "#/components/schemas/ResourceColumn"}}, "actions": map[string]any{"type": "array", "items": map[string]any{"$ref": "#/components/schemas/ResourceAction"}}}}
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
