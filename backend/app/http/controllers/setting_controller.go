package controllers

import (
	"errors"
	"strings"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"goravel/app/core/audit"
	"goravel/app/core/permission"
	"goravel/app/core/resource"
	"goravel/app/core/setting"
	apierrors "goravel/app/core/shared/errors"
	"goravel/app/core/shared/response"
)

const settingsUpdatePermission = "settings.update"

type SettingController struct {
	auth      *AuthController
	service   *setting.Service
	audit     *audit.Service
	initError error
}

func NewConfiguredSettingController(auth *AuthController) *SettingController {
	return NewConfiguredSettingControllerWithAudit(auth, nil)
}

func NewConfiguredSettingControllerWithAudit(auth *AuthController, auditService *audit.Service) *SettingController {
	mode, err := resource.ParseResourceProviderMode(facades.Config().GetString("resource.provider", "memory"))
	if err != nil {
		return &SettingController{auth: auth, audit: auditService, initError: err}
	}
	var repository setting.Repository
	switch mode {
	case resource.ResourceProviderMemory:
		repository = setting.NewMemoryRepository(
			setting.Setting{Namespace: "general", Key: "site_name", Value: "Go Vue Admin", ValueType: setting.ValueTypeString, Description: "后台显示名称"},
			setting.Setting{Namespace: "general", Key: "timezone", Value: "Asia/Shanghai", ValueType: setting.ValueTypeString, Description: "默认时区"},
			setting.Setting{Namespace: "auth", Key: "session_ttl", Value: 120, ValueType: setting.ValueTypeInteger, Description: "会话有效期（分钟）"},
		)
	case resource.ResourceProviderDatabase, resource.ResourceProviderMySQL:
		repository = setting.NewMySQLRepository()
	default:
		return &SettingController{auth: auth, audit: auditService, initError: errors.New("setting provider is not configured")}
	}
	return &SettingController{auth: auth, service: setting.NewService(repository), audit: auditService}
}

type settingRequest struct {
	Namespace   string            `json:"namespace"`
	Key         string            `json:"key"`
	Value       any               `json:"value"`
	ValueType   setting.ValueType `json:"value_type"`
	Description string            `json:"description"`
}

func (c *SettingController) Index(ctx http.Context) http.Response {
	if response := c.authorize(ctx, "view"); response != nil {
		return response
	}
	items, err := c.service.List(ctx, strings.TrimSpace(ctx.Request().Query("namespace")))
	if err != nil {
		return c.storageError(ctx, err)
	}
	return c.success(ctx, items)
}

func (c *SettingController) Store(ctx http.Context) http.Response {
	if response := c.authorize(ctx, "update"); response != nil {
		return response
	}
	var input settingRequest
	if err := ctx.Request().Bind(&input); err != nil {
		return c.error(ctx, 400, "INVALID_REQUEST", "设置请求格式无效", nil)
	}
	var before any
	if existing, getErr := c.service.Get(ctx, input.Namespace, input.Key); getErr == nil {
		before = existing
	}
	item, err := c.service.Upsert(ctx, setting.Setting{
		Namespace: input.Namespace, Key: input.Key, Value: input.Value,
		ValueType: input.ValueType, Description: input.Description,
	})
	if err != nil {
		return c.domainError(ctx, err)
	}
	if err := recordAudit(ctx, c.auth, c.audit, "settings.upsert", "setting", input.Namespace+"."+input.Key, before, item); err != nil {
		return c.storageError(ctx, err)
	}
	return ctx.Response().Status(201).Json(response.Success(item, response.Meta{RequestID: settingRequestID(ctx)}))
}

func (c *SettingController) Update(ctx http.Context) http.Response {
	if response := c.authorize(ctx, "update"); response != nil {
		return response
	}
	var input settingRequest
	if err := ctx.Request().Bind(&input); err != nil {
		return c.error(ctx, 400, "INVALID_REQUEST", "设置请求格式无效", nil)
	}
	namespace, key := ctx.Request().Route("namespace"), ctx.Request().Route("key")
	var before any
	if existing, getErr := c.service.Get(ctx, namespace, key); getErr == nil {
		before = existing
	}
	item, err := c.service.Upsert(ctx, setting.Setting{
		Namespace: namespace, Key: key, Value: input.Value,
		ValueType: input.ValueType, Description: input.Description,
	})
	if err != nil {
		return c.domainError(ctx, err)
	}
	if err := recordAudit(ctx, c.auth, c.audit, "settings.update", "setting", namespace+"."+key, before, item); err != nil {
		return c.storageError(ctx, err)
	}
	return c.success(ctx, item)
}

func (c *SettingController) Destroy(ctx http.Context) http.Response {
	if response := c.authorize(ctx, "update"); response != nil {
		return response
	}
	namespace, key := ctx.Request().Route("namespace"), ctx.Request().Route("key")
	before, _ := c.service.Get(ctx, namespace, key)
	if err := c.service.Delete(ctx, namespace, key); err != nil {
		return c.domainError(ctx, err)
	}
	if err := recordAudit(ctx, c.auth, c.audit, "settings.delete", "setting", namespace+"."+key, before, nil); err != nil {
		return c.storageError(ctx, err)
	}
	return c.success(ctx, map[string]bool{"deleted": true})
}

func (c *SettingController) authorize(ctx http.Context, action string) http.Response {
	if c.initError != nil || c.service == nil {
		return c.error(ctx, 503, "SETTING_PROVIDER_NOT_CONFIGURED", "设置存储尚未配置", nil)
	}
	if c.auth == nil {
		return c.error(ctx, 401, "UNAUTHENTICATED", "请先登录", nil)
	}
	user, err := c.auth.CurrentUser(ctx)
	if err != nil {
		return c.error(ctx, 401, "UNAUTHENTICATED", "请先登录", nil)
	}
	required := "settings.view"
	if action == "update" {
		required = settingsUpdatePermission
	}
	if permission.NewAuthorizer().Require(user.Permissions, required) != nil {
		return c.error(ctx, 403, "FORBIDDEN", "没有执行该操作的权限", map[string]string{"permission": required})
	}
	return nil
}

func (c *SettingController) domainError(ctx http.Context, err error) http.Response {
	switch {
	case errors.Is(err, setting.ErrInvalidSetting), errors.Is(err, setting.ErrInvalidValue), errors.Is(err, setting.ErrInvalidValueType):
		return c.error(ctx, 400, "INVALID_SETTING", "设置数据无效", nil)
	case errors.Is(err, setting.ErrNotFound):
		return c.error(ctx, 404, "SETTING_NOT_FOUND", "设置不存在", nil)
	default:
		return c.storageError(ctx, err)
	}
}

func (c *SettingController) storageError(ctx http.Context, err error) http.Response {
	return c.error(ctx, 503, "SETTING_STORAGE_UNAVAILABLE", "设置存储暂不可用", map[string]string{"reason": err.Error()})
}

func (c *SettingController) success(ctx http.Context, data any) http.Response {
	return ctx.Response().Success().Json(response.Success(data, response.Meta{RequestID: settingRequestID(ctx)}))
}

func (c *SettingController) error(ctx http.Context, status int, code, message string, details any) http.Response {
	return ctx.Response().Status(status).Json(apierrors.New(code, message, details))
}

func settingRequestID(ctx http.Context) string {
	requestID := ctx.Request().Header("X-Request-ID")
	if requestID == "" {
		requestID = "setting-request"
	}
	ctx.Response().Header("X-Request-ID", requestID)
	return requestID
}
