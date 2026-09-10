package core_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	frameworkhttp "github.com/goravel/framework/contracts/testing/http"

	"goravel/tests"
)

func TestMenuEndpointReturnsOnlyPermissionsGrantedToCurrentSession(t *testing.T) {
	testCase := new(tests.TestCase)

	unauthenticated, err := testCase.Http(t).Get("/api/menu")
	require.NoError(t, err)
	unauthenticated.AssertUnauthorized()

	csrfResponse, err := testCase.Http(t).Get("/csrf")
	require.NoError(t, err)
	sessionCookie := csrfResponse.Cookie("goravel_session")
	require.NotNil(t, sessionCookie)

	login, err := testCase.Http(t).
		WithCookie(sessionCookie).
		WithHeader("X-CSRF-TOKEN", responseHeaderToken(csrfResponse)).
		Post("/login", strings.NewReader(`{"email":"admin@example.com","password":"test-only-password"}`))
	require.NoError(t, err)
	login.AssertOk()

	rotatedCookie := login.Cookie("goravel_session")
	require.NotNil(t, rotatedCookie)

	menuResponse, err := testCase.Http(t).WithCookie(rotatedCookie).Get("/api/menu")
	require.NoError(t, err)
	menuResponse.AssertOk()

	payload, err := menuResponse.Json()
	require.NoError(t, err)
	items, ok := payload["data"].([]any)
	require.True(t, ok)
	require.Len(t, items, 7)
	require.Equal(t, "dashboard", items[0].(map[string]any)["id"])
	require.Equal(t, "settings", items[1].(map[string]any)["id"])
	require.Equal(t, "media", items[2].(map[string]any)["id"])
	require.Equal(t, "audit", items[3].(map[string]any)["id"])
	require.Equal(t, "users", items[4].(map[string]any)["id"])
	require.Equal(t, "roles", items[5].(map[string]any)["id"])
	require.Equal(t, "permissions", items[6].(map[string]any)["id"])
}

func responseHeaderToken(response frameworkhttp.Response) string {
	return response.Headers().Get("X-CSRF-TOKEN")
}
