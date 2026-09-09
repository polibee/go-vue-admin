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
}
