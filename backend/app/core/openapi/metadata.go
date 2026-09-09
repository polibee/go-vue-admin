package openapi

// CoreEndpoint describes the small set of platform endpoints that are
// stable enough to be consumed by generated clients.
type CoreEndpoint struct {
	Method      string
	Path        string
	Description string
}

var CoreEndpoints = []CoreEndpoint{
	{Method: "GET", Path: "/api/health", Description: "Returns service health and API contract version."},
	{Method: "GET", Path: "/api/csrf", Description: "Issues the session-bound CSRF token."},
	{Method: "GET", Path: "/api/auth/bootstrap", Description: "Returns development-only bootstrap credentials when enabled."},
	{Method: "POST", Path: "/api/login", Description: "Creates an HttpOnly cookie session."},
	{Method: "POST", Path: "/api/logout", Description: "Revokes the current cookie session."},
	{Method: "GET", Path: "/api/me", Description: "Returns the current authenticated user."},
	{Method: "GET", Path: "/api/menu", Description: "Returns navigation items allowed by the current session permissions."},
	{Method: "GET", Path: "/api/resources/demo", Description: "Lists demo resources with pagination, filters, search, and sorting."},
	{Method: "GET", Path: "/api/resources/demo/:id", Description: "Returns one demo resource."},
	{Method: "POST", Path: "/api/resources/demo", Description: "Creates one demo resource."},
	{Method: "PUT", Path: "/api/resources/demo/:id", Description: "Updates one demo resource."},
	{Method: "DELETE", Path: "/api/resources/demo/:id", Description: "Deletes one demo resource."},
	{Method: "POST", Path: "/api/resources/demo/bulk-delete", Description: "Deletes multiple demo resources."},
	{Method: "GET", Path: "/api/resources/users", Description: "Lists core users through the generic resource contract."},
	{Method: "GET", Path: "/api/resources/users/:id", Description: "Returns one core user."},
	{Method: "POST", Path: "/api/resources/users", Description: "Creates one core user."},
	{Method: "PUT", Path: "/api/resources/users/:id", Description: "Updates one core user."},
	{Method: "DELETE", Path: "/api/resources/users/:id", Description: "Deletes one core user."},
	{Method: "POST", Path: "/api/resources/users/bulk-delete", Description: "Deletes multiple core users."},
	{Method: "GET", Path: "/api/resources/roles", Description: "Lists core roles through the generic resource contract."},
	{Method: "GET", Path: "/api/resources/roles/:id", Description: "Returns one core role."},
	{Method: "POST", Path: "/api/resources/roles", Description: "Creates one core role."},
	{Method: "PUT", Path: "/api/resources/roles/:id", Description: "Updates one core role."},
	{Method: "DELETE", Path: "/api/resources/roles/:id", Description: "Deletes one core role."},
	{Method: "POST", Path: "/api/resources/roles/bulk-delete", Description: "Deletes multiple core roles."},
	{Method: "GET", Path: "/api/resources/permissions", Description: "Lists core permissions through the generic resource contract."},
	{Method: "GET", Path: "/api/resources/permissions/:id", Description: "Returns one core permission."},
	{Method: "POST", Path: "/api/resources/permissions", Description: "Creates one core permission."},
	{Method: "PUT", Path: "/api/resources/permissions/:id", Description: "Updates one core permission."},
	{Method: "DELETE", Path: "/api/resources/permissions/:id", Description: "Deletes one core permission."},
	{Method: "POST", Path: "/api/resources/permissions/bulk-delete", Description: "Deletes multiple core permissions."},
}
