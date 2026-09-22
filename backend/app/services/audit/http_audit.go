package auditservices

import (
	"encoding/json"
	"strings"
)

type HTTPAuditInput struct {
	Method            string
	Path              string
	Query             map[string]string
	Route             map[string]string
	RequestBody       any
	Status            int
	ContentType       string
	ResponseBody      []byte
	ResponseBodyBytes int
	ResponseTruncated bool
}

func ShouldAuditHTTP(method, path string) bool {
	if strings.EqualFold(method, "OPTIONS") || !strings.HasPrefix(path, "/api/v1/") {
		return false
	}
	return path != "/api/v1/admin/audit-logs"
}

func BuildHTTPAuditMetadata(input HTTPAuditInput) map[string]any {
	bodyBytes := input.ResponseBodyBytes
	if bodyBytes == 0 {
		bodyBytes = len(input.ResponseBody)
	}
	responseBody := any(nil)
	if input.ResponseTruncated {
		responseBody = map[string]any{"truncated": true, "original_bytes": bodyBytes}
	} else if len(input.ResponseBody) > 0 {
		if strings.Contains(strings.ToLower(input.ContentType), "json") {
			if err := json.Unmarshal(input.ResponseBody, &responseBody); err != nil {
				responseBody = map[string]any{"truncated": true, "invalid_json": true, "bytes": len(input.ResponseBody)}
			}
		}
	}
	metadata := map[string]any{
		"request": map[string]any{
			"method": input.Method,
			"path":   input.Path,
			"query":  BoundedValue(input.Query),
			"route":  BoundedValue(input.Route),
			"body":   BoundedValue(input.RequestBody),
		},
		"response": map[string]any{
			"status":       input.Status,
			"content_type": input.ContentType,
			"body":         BoundedValue(responseBody),
		},
	}
	if len(input.ResponseBody) > 0 && !strings.Contains(strings.ToLower(input.ContentType), "json") {
		response := metadata["response"].(map[string]any)
		delete(response, "body")
		response["body_bytes"] = bodyBytes
	}
	return metadata
}
