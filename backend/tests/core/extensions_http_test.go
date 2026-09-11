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

func TestSeparateModuleAndPluginCatalogs(t *testing.T) {
	tc := new(tests.TestCase)
	for _, path := range []string{"/api/modules", "/api/plugins", "/api/modules/example", "/api/plugins/example-plugin"} {
		guest, err := tc.Http(t).Get(path)
		require.NoError(t, err)
		guest.AssertStatus(401)
	}
	cookie, csrf := loginSettingsTestAdmin(t, tc)

	modules, err := tc.Http(t).WithCookie(cookie).Get("/api/modules")
	require.NoError(t, err)
	modules.AssertOk()
	modulePayload, err := modules.Json()
	require.NoError(t, err)
	moduleItems := modulePayload["data"].([]any)
	require.NotEmpty(t, moduleItems)
	require.Equal(t, "module", moduleItems[0].(map[string]any)["kind"])

	plugins, err := tc.Http(t).WithCookie(cookie).Get("/api/plugins")
	require.NoError(t, err)
	plugins.AssertOk()
	pluginPayload, err := plugins.Json()
	require.NoError(t, err)
	pluginItems := pluginPayload["data"].([]any)
	require.NotEmpty(t, pluginItems)
	require.Equal(t, "plugin", pluginItems[0].(map[string]any)["kind"])

	moduleDetail, err := tc.Http(t).WithCookie(cookie).Get("/api/modules/example")
	require.NoError(t, err)
	moduleDetail.AssertOk()
	moduleData, err := moduleDetail.Json()
	require.NoError(t, err)
	require.Equal(t, "module", moduleData["data"].(map[string]any)["kind"])
	resetPlugin, err := tc.Http(t).WithCookie(cookie).WithHeader("X-CSRF-TOKEN", csrf).
		Put("/api/plugins/example-plugin/state", strings.NewReader(`{"state":"enabled"}`))
	require.NoError(t, err)
	resetPlugin.AssertOk()

	pluginDetail, err := tc.Http(t).WithCookie(cookie).Get("/api/plugins/example-plugin")
	require.NoError(t, err)
	pluginDetail.AssertOk()
	pluginData, err := pluginDetail.Json()
	require.NoError(t, err)
	require.Equal(t, "plugin", pluginData["data"].(map[string]any)["kind"])

	change, err := tc.Http(t).WithCookie(cookie).WithHeader("X-CSRF-TOKEN", csrf).
		Put("/api/plugins/example-plugin/state", strings.NewReader(`{"state":"disabled"}`))
	require.NoError(t, err)
	change.AssertOk()
	modulesAfterDisable, err := tc.Http(t).WithCookie(cookie).Get("/api/modules")
	require.NoError(t, err)
	modulesAfterDisable.AssertOk()
	pluginAfterDisable, err := tc.Http(t).WithCookie(cookie).Get("/api/plugins/example-plugin")
	require.NoError(t, err)
	pluginAfterDisable.AssertStatus(404)
}
