package openapi

import (
	"strings"
	"testing"
)

func TestDocsHTMLPointsToLocalContract(t *testing.T) {
	html := DocsHTML()
	if !strings.Contains(html, `data-url="/api/openapi/admin.json"`) {
		t.Fatal("Scalar page must point to the local OpenAPI contract")
	}
}
