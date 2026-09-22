package auditservices

import (
	"encoding/json"
	"strings"
)

type HTTPAuditInput struct {
	Method       string
	Path         string
	Query        map[string]string
	Route        map[string]string
	RequestBody  any
	Status       int
	ContentType  string
	ResponseBody []byte
}

func ShouldAuditHTTP(method, path string) bool {
	if strings.EqualFold(method, "OPTIONS") || !strings.HasPrefix(path, "/api/v1/") {
		return false
	}
	return path != "/api/v1/admin/audit-logs"
}

func BuildHTTPAuditMetadata(input HTTPAuditInput) map[string]any {
	responseBody := any(nil)
	if len(input.ResponseBody) > 0 {
		if strings.Contains(strings.ToLower(input.ContentType), "json") {
			if err := json.Unmarshal(input.ResponseBody, &responseBody); err != nil {
				responseBody = string(input.ResponseBody)
			}
		} else {
			responseBody = string(input.ResponseBody)
		}
	}
	return map[string]any{
		"request": map[string]any{
			"method": input.Method,
			"path":   input.Path,
			"query":  input.Query,
			"route":  input.Route,
			"body":   input.RequestBody,
		},
		"response": map[string]any{
			"status":       input.Status,
			"content_type": input.ContentType,
			"body":         responseBody,
		},
	}
}
