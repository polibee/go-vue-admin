package controllers

import (
	"errors"
	"strings"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"goravel/app/core/audit"
	"goravel/app/core/permission"
	"goravel/app/core/resource"
	apierrors "goravel/app/core/shared/errors"
	"goravel/app/core/shared/response"
)

type AuditController struct {
	auth      *AuthController
	service   *audit.Service
	initError error
}

func NewConfiguredAuditService() (*audit.Service, error) {
	mode, err := resource.ParseResourceProviderMode(facades.Config().GetString("resource.provider", "memory"))
	if err != nil {
		return nil, err
	}
	switch mode {
	case resource.ResourceProviderMemory:
		return audit.NewService(audit.NewMemoryRepository()), nil
	case resource.ResourceProviderMySQL:
		return audit.NewService(audit.NewMySQLRepository()), nil
	default:
		return nil, errors.New("audit provider is not configured")
	}
}

func NewAuditController(auth *AuthController, service *audit.Service) *AuditController {
	return &AuditController{auth: auth, service: service}
}

func NewConfiguredAuditController(auth *AuthController) *AuditController {
	service, err := NewConfiguredAuditService()
	return &AuditController{auth: auth, service: service, initError: err}
}

func (c *AuditController) Index(ctx http.Context) http.Response {
	if response := c.authorize(ctx); response != nil {
		return response
	}
	items, err := c.service.List(ctx)
	if err != nil {
		return c.storageError(ctx, err)
	}
	return c.success(ctx, items)
}

func (c *AuditController) Show(ctx http.Context) http.Response {
	if response := c.authorize(ctx); response != nil {
		return response
	}
	item, err := c.service.Get(ctx, strings.TrimSpace(ctx.Request().Route("id")))
	if err != nil {
		if errors.Is(err, audit.ErrNotFound) {
			return c.error(ctx, 404, "AUDIT_NOT_FOUND", "审计记录不存在", nil)
		}
		return c.storageError(ctx, err)
	}
	return c.success(ctx, item)
}

func (c *AuditController) authorize(ctx http.Context) http.Response {
	if c.initError != nil || c.service == nil {
		return c.error(ctx, 503, "AUDIT_PROVIDER_NOT_CONFIGURED", "审计存储尚未配置", nil)
	}
	if c.auth == nil {
		return c.error(ctx, 401, "UNAUTHENTICATED", "请先登录", nil)
	}
	user, err := c.auth.CurrentUser(ctx)
	if err != nil {
		return c.error(ctx, 401, "UNAUTHENTICATED", "请先登录", nil)
	}
	if permission.NewAuthorizer().Require(user.Permissions, "audit.view") != nil {
		return c.error(ctx, 403, "FORBIDDEN", "没有执行该操作的权限", map[string]string{"permission": "audit.view"})
	}
	return nil
}

func (c *AuditController) storageError(ctx http.Context, err error) http.Response {
	return c.error(ctx, 503, "AUDIT_STORAGE_UNAVAILABLE", "审计存储暂不可用", map[string]string{"reason": err.Error()})
}

func (c *AuditController) success(ctx http.Context, data any) http.Response {
	return ctx.Response().Success().Json(response.Success(data, response.Meta{RequestID: auditRequestID(ctx)}))
}

func (c *AuditController) error(ctx http.Context, status int, code, message string, details any) http.Response {
	return ctx.Response().Status(status).Json(apierrors.New(code, message, details))
}

func auditRequestID(ctx http.Context) string {
	requestID := ctx.Request().Header("X-Request-ID")
	if requestID == "" {
		requestID = "audit-request"
	}
	ctx.Response().Header("X-Request-ID", requestID)
	return requestID
}

func recordAudit(ctx http.Context, auth *AuthController, service *audit.Service, action, resourceType, resourceID string, before, after any) error {
	if service == nil {
		return nil
	}
	if auth == nil {
		return errors.New("audit actor is not configured")
	}
	user, err := auth.CurrentUser(ctx)
	if err != nil {
		return err
	}
	_, err = service.Record(ctx, audit.Entry{
		ActorID: user.ID, ActorEmail: user.Email, Action: action, ResourceType: resourceType, ResourceID: resourceID,
		Before: before, After: after, IP: ctx.Request().Ip(), UserAgent: ctx.Request().Header("User-Agent"),
	})
	return err
}
