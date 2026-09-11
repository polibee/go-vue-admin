package controllers

import (
	"errors"
	"github.com/goravel/framework/contracts/http"
	example "go-vue-admin-module/example"
	"goravel/app/core/module"
	"goravel/app/core/permission"
	"goravel/app/core/plugin"
	apierrors "goravel/app/core/shared/errors"
	"goravel/app/core/shared/response"
)

// ExtensionController owns the compiled catalog for this application process.
// Plugin state survives browser reloads, but is reset by a backend restart.
type ExtensionController struct {
	auth      *AuthController
	modules   *module.Registry
	app       *module.ApplicationContainer
	plugins   *plugin.Registry
	initError error
}

type extensionRecord struct {
	ID      string                  `json:"id"`
	Name    string                  `json:"name"`
	Kind    string                  `json:"kind"`
	State   string                  `json:"state"`
	Version string                  `json:"version,omitempty"`
	Config  *extensionConfigSummary `json:"config,omitempty"`
	Message string                  `json:"message,omitempty"`
}

type extensionConfigSummary struct {
	Route        string   `json:"route"`
	Label        string   `json:"label"`
	Permission   string   `json:"permission"`
	Schema       string   `json:"schema,omitempty"`
	SecretFields []string `json:"secret_fields,omitempty"`
	Status       string   `json:"status"`
}

type builtinExamplePlugin struct{}

func (builtinExamplePlugin) Manifest() plugin.PluginManifest {
	return plugin.PluginManifest{ID: "example-plugin", Name: "示例插件", Version: "1.0.0", Runtime: plugin.RuntimeBuiltin, UICompatibility: plugin.UICompatibilityShadcnVue, Config: &plugin.ConfigDeclaration{Route: "/admin/plugins/example-plugin/config", Label: "示例插件配置", Permission: "plugins.configure", Schema: "example-plugin-config", SecretFields: []string{"apiKey"}}}
}
func (builtinExamplePlugin) Enable() error  { return nil }
func (builtinExamplePlugin) Disable() error { return nil }

func NewExtensionController(auth *AuthController) *ExtensionController {
	c := &ExtensionController{auth: auth, modules: module.NewRegistry(), app: module.NewApplication(), plugins: plugin.NewRegistry()}
	if err := c.modules.Register(example.Module{}); err != nil {
		c.initError = err
		return c
	}
	if err := c.modules.RegisterAll(c.app); err != nil {
		c.initError = err
		return c
	}
	if err := c.modules.BootAll(c.app); err != nil {
		c.initError = err
		return c
	}
	c.initError = c.plugins.Register(builtinExamplePlugin{})
	return c
}

func (c *ExtensionController) authorize(ctx http.Context, required string) http.Response {
	user, err := c.auth.CurrentUser(ctx)
	if err != nil {
		return ctx.Response().Status(401).Json(apierrors.New("UNAUTHENTICATED", "请先登录", nil))
	}
	if permission.NewAuthorizer().Require(user.Permissions, required) != nil {
		return ctx.Response().Status(403).Json(apierrors.New("FORBIDDEN", "没有执行该操作的权限", nil))
	}
	if c.initError != nil {
		return ctx.Response().Status(503).Json(apierrors.New("EXTENSION_UNAVAILABLE", "扩展初始化失败", nil))
	}
	return nil
}

func (c *ExtensionController) Index(ctx http.Context) http.Response {
	if denied := c.authorize(ctx, "dashboard.view"); denied != nil {
		return denied
	}
	items := c.catalog()
	return ctx.Response().Success().Json(response.Success(items, response.Meta{}))
}

func (c *ExtensionController) Modules(ctx http.Context) http.Response {
	if denied := c.authorize(ctx, "dashboard.view"); denied != nil {
		return denied
	}
	items := []extensionRecord{}
	for _, id := range c.modules.Names() {
		items = append(items, extensionRecord{ID: id, Name: "示例模块", Kind: "module", State: "enabled"})
	}
	return ctx.Response().Success().Json(response.Success(items, response.Meta{}))
}

func (c *ExtensionController) Plugins(ctx http.Context) http.Response {
	if denied := c.authorize(ctx, "dashboard.view"); denied != nil {
		return denied
	}
	items := []extensionRecord{}
	for _, item := range c.plugins.List() {
		items = append(items, c.pluginRecord(item))
	}
	return ctx.Response().Success().Json(response.Success(items, response.Meta{}))
}

func (c *ExtensionController) ModuleShow(ctx http.Context) http.Response {
	return c.showCatalogItem(ctx, "module")
}

func (c *ExtensionController) PluginShow(ctx http.Context) http.Response {
	if denied := c.authorize(ctx, "dashboard.view"); denied != nil {
		return denied
	}
	id := ctx.Request().Route("id")
	for _, item := range c.plugins.List() {
		if item.Manifest.ID == id {
			return ctx.Response().Success().Json(response.Success(c.pluginRecord(item), response.Meta{}))
		}
	}
	return ctx.Response().Status(404).Json(apierrors.New("EXTENSION_NOT_FOUND", "插件不存在", nil))
}

func (c *ExtensionController) catalog() []extensionRecord {
	items := []extensionRecord{}
	for _, id := range c.modules.Names() {
		items = append(items, extensionRecord{ID: id, Name: "示例模块", Kind: "module", State: "enabled"})
	}
	for _, item := range c.plugins.List() {
		items = append(items, c.pluginRecord(item))
	}
	return items
}

func (c *ExtensionController) pluginRecord(item plugin.PluginInfo) extensionRecord {
	var config *extensionConfigSummary
	if item.Manifest.Config != nil {
		declaration := item.Manifest.Config
		config = &extensionConfigSummary{Route: declaration.Route, Label: declaration.Label, Permission: declaration.Permission, Schema: declaration.Schema, SecretFields: append([]string(nil), declaration.SecretFields...), Status: "not_configured"}
	}
	return extensionRecord{ID: item.Manifest.ID, Name: item.Manifest.Name, Kind: "plugin", State: string(item.State), Version: item.Manifest.Version, Config: config}
}

func (c *ExtensionController) Update(ctx http.Context) http.Response {
	return c.updatePluginState(ctx)
}

func (c *ExtensionController) UpdatePlugin(ctx http.Context) http.Response {
	return c.updatePluginState(ctx)
}

func (c *ExtensionController) updatePluginState(ctx http.Context) http.Response {
	if denied := c.authorize(ctx, "plugins.manage"); denied != nil {
		return denied
	}
	var input struct {
		State string `json:"state"`
	}
	if err := ctx.Request().Bind(&input); err != nil || (input.State != "enabled" && input.State != "disabled") {
		return ctx.Response().Status(400).Json(apierrors.New("INVALID_STATE", "状态必须为 enabled 或 disabled", nil))
	}
	id := ctx.Request().Route("id")
	var err error
	if input.State == "enabled" {
		err = c.plugins.Enable(id)
	} else {
		err = c.plugins.Disable(id)
	}
	if err != nil {
		status := 409
		if errors.Is(err, plugin.ErrPluginNotFound) {
			status = 404
		}
		return ctx.Response().Status(status).Json(apierrors.New("PLUGIN_TRANSITION_FAILED", "插件状态更新失败", nil))
	}
	return c.Plugins(ctx)
}

func (c *ExtensionController) Show(ctx http.Context) http.Response {
	return c.showCatalogItem(ctx, "")
}

func (c *ExtensionController) showCatalogItem(ctx http.Context, expectedKind string) http.Response {
	if denied := c.authorize(ctx, "dashboard.view"); denied != nil {
		return denied
	}
	id := ctx.Request().Route("id")
	message := ""
	if id == "example" {
		if value, ok := c.app.Service("example.message"); ok {
			message, _ = value.(string)
		}
	}
	if id == "example-plugin" && c.plugins.State(id) == plugin.PluginEnabled {
		message = "示例插件正在运行"
	}
	if message == "" {
		return ctx.Response().Status(404).Json(apierrors.New("EXTENSION_NOT_ACTIVE", "扩展不存在或已停用", nil))
	}
	kind, name := "module", "示例模块"
	if id == "example-plugin" {
		kind, name = "plugin", "示例插件"
	}
	if expectedKind != "" && kind != expectedKind {
		return ctx.Response().Status(404).Json(apierrors.New("EXTENSION_NOT_FOUND", "模块或插件不存在", nil))
	}
	if kind == "plugin" {
		var item extensionRecord
		for _, pluginInfo := range c.plugins.List() {
			if pluginInfo.Manifest.ID == id {
				item = c.pluginRecord(pluginInfo)
				break
			}
		}
		item.Message = message
		return ctx.Response().Success().Json(response.Success(item, response.Meta{}))
	}
	return ctx.Response().Success().Json(response.Success(map[string]string{"id": id, "name": name, "kind": kind, "state": "enabled", "message": message}, response.Meta{}))
}
