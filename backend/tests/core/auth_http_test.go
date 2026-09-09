package core_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	frameworkhttp "github.com/goravel/framework/contracts/testing/http"

	"goravel/tests"
)

func TestAuthEndpointsUseRotatingHttpOnlySessionCookie(t *testing.T) {
	testCase := new(tests.TestCase)

	csrfResponse, err := testCase.Http(t).Get("/csrf")
	require.NoError(t, err)
	csrfResponse.AssertOk()

	csrfToken := responseDataString(t, csrfResponse, "token")
	sessionCookie := csrfResponse.Cookie("goravel_session")
	require.NotNil(t, sessionCookie)
	require.True(t, sessionCookie.HttpOnly)
	require.Equal(t, http.SameSiteLaxMode, sessionCookie.SameSite)

	unauthenticated, err := testCase.Http(t).WithCookie(sessionCookie).Get("/me")
	require.NoError(t, err)
	unauthenticated.AssertUnauthorized()

	login, err := testCase.Http(t).
		WithCookie(sessionCookie).
		WithHeader("X-CSRF-TOKEN", csrfToken).
		Post("/login", strings.NewReader(`{"email":"admin@example.com","password":"test-only-password"}`))
	require.NoError(t, err)
	login.AssertOk()

	rotatedCookie := login.Cookie("goravel_session")
	require.NotNil(t, rotatedCookie)
	require.NotEqual(t, sessionCookie.Value, rotatedCookie.Value)
	require.True(t, rotatedCookie.HttpOnly)
	rotatedCSRFToken := login.Headers().Get("X-CSRF-TOKEN")
	require.NotEmpty(t, rotatedCSRFToken)

	me, err := testCase.Http(t).WithCookie(rotatedCookie).Get("/me")
	require.NoError(t, err)
	me.AssertOk()
	require.Equal(t, "admin@example.com", responseDataString(t, me, "email"))

	logout, err := testCase.Http(t).
		WithCookie(rotatedCookie).
		WithHeader("X-CSRF-TOKEN", rotatedCSRFToken).
		Post("/logout", nil)
	require.NoError(t, err)
	logout.AssertOk()

	afterLogout, err := testCase.Http(t).WithCookie(rotatedCookie).Get("/me")
	require.NoError(t, err)
	afterLogout.AssertUnauthorized()
}

func TestAuthMutationsRequireCsrfToken(t *testing.T) {
	testCase := new(tests.TestCase)

	csrfResponse, err := testCase.Http(t).Get("/csrf")
	require.NoError(t, err)
	sessionCookie := csrfResponse.Cookie("goravel_session")
	require.NotNil(t, sessionCookie)

	login, err := testCase.Http(t).
		WithCookie(sessionCookie).
		Post("/login", strings.NewReader(`{"email":"admin@example.com","password":"test-only-password"}`))
	require.NoError(t, err)
	login.AssertStatus(419)
}

func responseDataString(t *testing.T, response frameworkhttp.Response, key string) string {
	t.Helper()
	payload, err := response.Json()
	require.NoError(t, err)
	data, ok := payload["data"].(map[string]any)
	require.True(t, ok)
	value, ok := data[key].(string)
	require.True(t, ok)
	return value
}
