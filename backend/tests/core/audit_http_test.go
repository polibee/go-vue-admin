package core_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"goravel/tests"
)

func TestAuditEndpointExposesSettingsMutationDiff(t *testing.T) {
	testCase := new(tests.TestCase)
	cookie, csrf := loginAuditTestAdmin(t, testCase)

	created, err := testCase.Http(t).
		WithCookie(cookie).
		WithHeader("X-CSRF-TOKEN", csrf).
		Post("/api/settings", strings.NewReader(`{"namespace":"audit-test","key":"enabled","value":true,"value_type":"boolean"}`))
	require.NoError(t, err)
	created.AssertStatus(201)

	list, err := testCase.Http(t).WithCookie(cookie).Get("/api/audit")
	require.NoError(t, err)
	list.AssertOk()
	payload, err := list.Json()
	require.NoError(t, err)
	items := payload["data"].([]any)
	require.NotEmpty(t, items)

	var auditID string
	for _, raw := range items {
		item := raw.(map[string]any)
		if item["resource_id"] == "audit-test.enabled" {
			auditID = item["id"].(string)
			require.Equal(t, "settings.upsert", item["action"])
			require.NotNil(t, item["after"])
		}
	}
	require.NotEmpty(t, auditID)

	detail, err := testCase.Http(t).WithCookie(cookie).Get("/api/audit/" + auditID)
	require.NoError(t, err)
	detail.AssertOk()
}

func loginAuditTestAdmin(t *testing.T, testCase *tests.TestCase) (cookie *http.Cookie, csrf string) {
	t.Helper()
	csrfResponse, err := testCase.Http(t).Get("/csrf")
	require.NoError(t, err)
	sessionCookie := csrfResponse.Cookie("goravel_session")
	require.NotNil(t, sessionCookie)

	login, err := testCase.Http(t).
		WithCookie(sessionCookie).
		WithHeader("X-CSRF-TOKEN", csrfResponse.Headers().Get("X-CSRF-TOKEN")).
		Post("/login", strings.NewReader(`{"email":"admin@example.com","password":"test-only-password"}`))
	require.NoError(t, err)
	login.AssertOk()
	rotated := login.Cookie("goravel_session")
	require.NotNil(t, rotated)
	return rotated, login.Headers().Get("X-CSRF-TOKEN")
}
