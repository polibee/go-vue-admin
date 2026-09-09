package core_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"goravel/tests"
)

func TestResourceListRequiresAuthentication(t *testing.T) {
	testCase := new(tests.TestCase)

	response, err := testCase.Http(t).Get("/api/resources/demo")
	require.NoError(t, err)
	response.AssertUnauthorized()
}

func TestResourceListReturnsPaginatedEnvelopeForCurrentUser(t *testing.T) {
	testCase := new(tests.TestCase)
	cookie, _ := loginResourceTestAdmin(t, testCase)

	response, err := testCase.Http(t).WithCookie(cookie).Get("/api/resources/demo?page=1&per_page=10&sort=name")
	require.NoError(t, err)
	response.AssertOk()

	payload, err := response.Json()
	require.NoError(t, err)
	data, ok := payload["data"].([]any)
	require.True(t, ok)
	require.Len(t, data, 2)
	meta, ok := payload["meta"].(map[string]any)
	require.True(t, ok)
	pagination, ok := meta["pagination"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, float64(2), pagination["total"])
}

func TestResourceCreateRejectsInvalidStatus(t *testing.T) {
	testCase := new(tests.TestCase)
	cookie, csrf := loginResourceTestAdmin(t, testCase)

	response, err := testCase.Http(t).
		WithCookie(cookie).
		WithHeader("X-CSRF-TOKEN", csrf).
		Post("/api/resources/demo", strings.NewReader(`{"id":"invalid-status","name":"Invalid","status":"paused","owner":"Admin"}`))
	require.NoError(t, err)
	response.AssertStatus(400)
}

func TestResourceGetReturnsNotFound(t *testing.T) {
	testCase := new(tests.TestCase)
	cookie, _ := loginResourceTestAdmin(t, testCase)

	response, err := testCase.Http(t).WithCookie(cookie).Get("/api/resources/demo/missing")
	require.NoError(t, err)
	response.AssertNotFound()
}

func loginResourceTestAdmin(t *testing.T, testCase *tests.TestCase) (cookie *http.Cookie, csrf string) {
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
