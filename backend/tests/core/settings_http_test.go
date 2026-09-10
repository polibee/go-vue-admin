package core_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"goravel/tests"
)

func TestSettingsEndpointListsAndMutatesTypedSettings(t *testing.T) {
	testCase := new(tests.TestCase)
	cookie, csrf := loginSettingsTestAdmin(t, testCase)

	list, err := testCase.Http(t).WithCookie(cookie).Get("/api/settings?namespace=general")
	require.NoError(t, err)
	list.AssertOk()
	payload, err := list.Json()
	require.NoError(t, err)
	items, ok := payload["data"].([]any)
	require.True(t, ok)
	require.Len(t, items, 2)

	created, err := testCase.Http(t).
		WithCookie(cookie).
		WithHeader("X-CSRF-TOKEN", csrf).
		Post("/api/settings", strings.NewReader(`{"namespace":"general","key":"maintenance","value":true,"value_type":"boolean","description":"维护模式"}`))
	require.NoError(t, err)
	created.AssertStatus(201)

	updated, err := testCase.Http(t).
		WithCookie(cookie).
		WithHeader("X-CSRF-TOKEN", csrf).
		Put("/api/settings/general/maintenance", strings.NewReader(`{"value":false,"value_type":"boolean"}`))
	require.NoError(t, err)
	updated.AssertOk()

	deleted, err := testCase.Http(t).
		WithCookie(cookie).
		WithHeader("X-CSRF-TOKEN", csrf).
		Delete("/api/settings/general/maintenance", nil)
	require.NoError(t, err)
	deleted.AssertOk()
}

func loginSettingsTestAdmin(t *testing.T, testCase *tests.TestCase) (cookie *http.Cookie, csrf string) {
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
