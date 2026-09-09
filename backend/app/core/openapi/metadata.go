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
	{Method: "GET", Path: "/csrf", Description: "Issues the session-bound CSRF token."},
	{Method: "POST", Path: "/login", Description: "Creates an HttpOnly cookie session."},
	{Method: "POST", Path: "/logout", Description: "Revokes the current cookie session."},
	{Method: "GET", Path: "/me", Description: "Returns the current authenticated user."},
	{Method: "GET", Path: "/api/menu", Description: "Returns navigation items allowed by the current session permissions."},
	{Method: "GET", Path: "/api/resources/demo", Description: "Lists demo resources with pagination, filters, search, and sorting."},
	{Method: "GET", Path: "/api/resources/demo/:id", Description: "Returns one demo resource."},
	{Method: "POST", Path: "/api/resources/demo", Description: "Creates one demo resource."},
	{Method: "PUT", Path: "/api/resources/demo/:id", Description: "Updates one demo resource."},
	{Method: "DELETE", Path: "/api/resources/demo/:id", Description: "Deletes one demo resource."},
	{Method: "POST", Path: "/api/resources/demo/bulk-delete", Description: "Deletes multiple demo resources."},
}
