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
				"UserStatus": map[string]any{"type": "string", "enum": []string{"active", "disabled", "locked"}},
				"Error":      map[string]any{"type": "object", "required": []string{"code"}, "properties": map[string]any{"code": map[string]any{"type": "string"}}},
			},
		},
		"paths": map[string]any{
			"/auth/login": map[string]any{"post": map[string]any{
				"security": []any{}, "operationId": "login", "requestBody": jsonBody("LoginRequest", map[string]any{"type": "object", "required": []string{"email", "password"}, "properties": map[string]any{"email": map[string]any{"type": "string", "format": "email"}, "password": map[string]any{"type": "string", "format": "password"}}}),
				"responses": map[string]any{"200": jsonResponse("LoginResponse"), "403": errorResponse()},
			}},
			"/auth/refresh": map[string]any{"post": map[string]any{
				"security": []any{}, "operationId": "refresh", "responses": map[string]any{"200": jsonResponse("RefreshResponse"), "401": errorResponse()},
			}},
			"/auth/logout-all":            map[string]any{"post": operation("logoutAll")},
			"/auth/me":                    map[string]any{"get": operation("currentUser")},
			"/admin/resources":            map[string]any{"get": operation("listResources")},
			"/admin/resources/{resource}": map[string]any{"get": operation("listResourceRows")},
			"/admin/overview":             map[string]any{"get": operation("adminOverview")},
			"/admin/audit-logs":           map[string]any{"get": operation("auditLogs")},
			"/admin/users/status": map[string]any{"put": map[string]any{
				"operationId": "bulkSetUserStatus", "requestBody": jsonBody("BulkUserStatusRequest", map[string]any{"type": "object", "required": []string{"user_ids", "status"}, "properties": map[string]any{"user_ids": map[string]any{"type": "array", "items": map[string]any{"type": "integer", "format": "int64"}}, "status": map[string]any{"$ref": "#/components/schemas/UserStatus"}}}),
				"responses": map[string]any{"204": map[string]any{"description": "Status updated"}, "409": errorResponse(), "422": errorResponse()},
			}},
		},
	}
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
