package generator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

type APIAuditOptions struct{ RootDir string }

type APIAuditFinding struct {
	Method   string `json:"method"`
	Path     string `json:"path"`
	Layer    string `json:"layer"`
	Message  string `json:"message"`
	Severity string `json:"severity"`
}

// AuditAPI statically checks the route, contract, generated-client, and
// frontend-consumer boundaries without executing the application.
func AuditAPI(options APIAuditOptions) ([]APIAuditFinding, error) {
	root := options.RootDir
	if root == "" {
		root = "."
	}
	routesPath := filepath.Join(root, "backend", "routes", "web.go")
	routes, err := os.ReadFile(routesPath)
	if err != nil {
		return nil, fmt.Errorf("read routes: %w", err)
	}
	client, err := os.ReadFile(filepath.Join(root, "admin", "src", "generated", "api", "client.ts"))
	if err != nil {
		return nil, fmt.Errorf("read generated client: %w", err)
	}
	documentBytes, err := os.ReadFile(filepath.Join(root, "contracts", "openapi", "openapi.json"))
	if err != nil {
		return nil, fmt.Errorf("read OpenAPI document: %w", err)
	}
	var document map[string]any
	if err := json.Unmarshal(documentBytes, &document); err != nil {
		return nil, fmt.Errorf("decode OpenAPI document: %w", err)
	}
	findings := make([]APIAuditFinding, 0)
	paths, _ := document["paths"].(map[string]any)
	for _, match := range regexp.MustCompile(`facades\.Route\(\)\.(Get|Post|Put|Patch|Delete)\("([^"]+)"`).FindAllSubmatch(routes, -1) {
		method, path := strings.ToUpper(string(match[1])), normalizeAuditPath(string(match[2]))
		if !strings.HasPrefix(path, "/api/") || strings.HasPrefix(path, "/api/docs/") {
			continue
		}
		if _, ok := paths[path]; !ok {
			findings = append(findings, auditFinding(method, path, "openapi", "route is not documented in the OpenAPI contract", "error"))
		}
	}
	for path, item := range paths {
		if !auditClientPath(path) {
			continue
		}
		operations, _ := item.(map[string]any)
		for method, operation := range operations {
			op, _ := operation.(map[string]any)
			operationID, _ := op["operationId"].(string)
			if operationID != "" && !hasGeneratedOperation(client, operationID) {
				findings = append(findings, auditFinding(strings.ToUpper(method), path, "generated-client", "OpenAPI operation is missing from the generated admin client", "error"))
			}
		}
	}
	if generatedAPIDrift(root, documentBytes) {
		findings = append(findings, auditFinding("", "", "generated-artifacts", "OpenAPI, schema, or generated client output is stale; run admin-gen api", "error"))
	}
	findings = append(findings, auditDirectCalls(root)...)
	sort.Slice(findings, func(i, j int) bool {
		if findings[i].Path == findings[j].Path {
			return findings[i].Method < findings[j].Method
		}
		return findings[i].Path < findings[j].Path
	})
	return findings, nil
}

func hasGeneratedOperation(client []byte, operationID string) bool {
	return bytes.Contains(client, []byte(operationID+"(")) || bytes.Contains(client, []byte(operationID+":"))
}

func auditFinding(method, path, layer, message, severity string) APIAuditFinding {
	return APIAuditFinding{Method: method, Path: path, Layer: layer, Message: message, Severity: severity}
}

func normalizeAuditPath(path string) string {
	path = strings.ReplaceAll(path, ":id", "{id}")
	path = strings.ReplaceAll(path, ":namespace", "{namespace}")
	path = strings.ReplaceAll(path, ":key", "{key}")
	return path
}

func auditClientPath(path string) bool {
	// Resource operationIds intentionally match generated client method names.
	// Settings/media/audit/extensions use a curated facade with different
	// method names, so their route coverage is checked by frontend consumers.
	return strings.HasPrefix(path, "/api/resources/")
}

func generatedAPIDrift(root string, expected []byte) bool {
	actual, err := os.ReadFile(filepath.Join(root, "contracts", "openapi", "openapi.json"))
	if err != nil {
		return true
	}
	var want, got any
	if json.Unmarshal(expected, &want) != nil || json.Unmarshal(actual, &got) != nil {
		return true
	}
	normalizedWant, _ := json.Marshal(want)
	normalizedGot, _ := json.Marshal(got)
	if !bytes.Equal(normalizedWant, normalizedGot) {
		return true
	}
	temp, err := os.MkdirTemp("", "admin-api-audit-")
	if err != nil {
		return true
	}
	defer os.RemoveAll(temp)
	if err := GenerateAPI(APIOptions{RootDir: temp}); err != nil {
		return true
	}
	for _, name := range []string{"admin/src/generated/api/models.ts", "admin/src/generated/api/client.ts", "admin/src/generated/api/index.ts"} {
		wantFile, e1 := os.ReadFile(filepath.Join(temp, name))
		gotFile, e2 := os.ReadFile(filepath.Join(root, name))
		if e1 != nil || e2 != nil || !bytes.Equal(wantFile, gotFile) {
			return true
		}
	}
	return false
}

func auditDirectCalls(root string) []APIAuditFinding {
	var findings []APIAuditFinding
	paths, _ := filepath.Glob(filepath.Join(root, "admin", "src", "**", "*"))
	_ = paths // filepath.Glob does not recurse; walk below handles nested modules.
	_ = filepath.Walk(filepath.Join(root, "admin", "src"), func(path string, info os.FileInfo, err error) error {
		if err != nil || info == nil || info.IsDir() || strings.Contains(path, string(filepath.Separator)+"generated"+string(filepath.Separator)) {
			return nil
		}
		if !strings.HasSuffix(path, ".ts") && !strings.HasSuffix(path, ".vue") {
			return nil
		}
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return nil
		}
		for _, match := range regexp.MustCompile(`['"](/api/(?:resources|settings|media|audit|extensions)[^'"]*)['"]`).FindAllSubmatch(data, -1) {
			findings = append(findings, auditFinding("", string(match[1]), "frontend", "frontend contains a direct generated-API path; use the generated client", "warning"))
		}
		return nil
	})
	return findings
}
