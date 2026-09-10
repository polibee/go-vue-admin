package core_test

import (
	"github.com/stretchr/testify/require"
	"goravel/tests"
	"strings"
	"testing"
)

func TestExtensionLifecycle(t *testing.T) {
	tc := new(tests.TestCase)
	guest, err := tc.Http(t).Get("/api/extensions")
	require.NoError(t, err)
	guest.AssertStatus(401)
	cookie, csrf := loginSettingsTestAdmin(t, tc)
	list, err := tc.Http(t).WithCookie(cookie).Get("/api/extensions")
	require.NoError(t, err)
	list.AssertOk()
	module, err := tc.Http(t).WithCookie(cookie).Get("/api/extensions/example")
	require.NoError(t, err)
	module.AssertOk()
	for _, state := range []string{"enabled", "enabled", "disabled", "disabled"} {
		change, err := tc.Http(t).WithCookie(cookie).WithHeader("X-CSRF-TOKEN", csrf).
			Put("/api/extensions/example-plugin", strings.NewReader(`{"state":"`+state+`"}`))
		require.NoError(t, err)
		change.AssertOk()
		probe, err := tc.Http(t).WithCookie(cookie).Get("/api/extensions/example-plugin")
		require.NoError(t, err)
		if state == "enabled" {
			probe.AssertOk()
		} else {
			probe.AssertStatus(404)
		}
	}
	invalid, err := tc.Http(t).WithCookie(cookie).WithHeader("X-CSRF-TOKEN", csrf).
		Put("/api/extensions/example-plugin", strings.NewReader(`{"state":"other"}`))
	require.NoError(t, err)
	invalid.AssertStatus(400)
}
