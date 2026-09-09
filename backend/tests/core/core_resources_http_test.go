package core_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"goravel/tests"
)

func TestCoreResourcesRequireAuthentication(t *testing.T) {
	testCase := new(tests.TestCase)

	for _, path := range []string{"/api/resources/users", "/api/resources/roles", "/api/resources/permissions"} {
		response, err := testCase.Http(t).Get(path)
		require.NoError(t, err)
		response.AssertUnauthorized()
	}
}

func TestCoreResourcesUseGenericCrudRoutesForCurrentAdmin(t *testing.T) {
	testCase := new(tests.TestCase)
	cookie, _ := loginResourceTestAdmin(t, testCase)

	for _, item := range []struct {
		path string
		id   string
	}{
		{path: "/api/resources/users", id: "bootstrap-admin"},
		{path: "/api/resources/roles", id: "platform-admin"},
		{path: "/api/resources/permissions", id: "dashboard.view"},
	} {
		response, err := testCase.Http(t).WithCookie(cookie).Get(item.path)
		require.NoError(t, err)
		response.AssertOk()
		payload, err := response.Json()
		require.NoError(t, err)
		rows, ok := payload["data"].([]any)
		require.True(t, ok)
		require.NotEmpty(t, rows)
		require.Equal(t, item.id, rows[0].(map[string]any)["id"])
	}
}
