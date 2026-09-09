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
}
