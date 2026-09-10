package controllers

import (
	"errors"
	"os"
	"strings"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"goravel/app/core/audit"
	"goravel/app/core/media"
	"goravel/app/core/permission"
	"goravel/app/core/resource"
	apierrors "goravel/app/core/shared/errors"
	"goravel/app/core/shared/response"
)

type MediaController struct {
	auth      *AuthController
	service   *media.Service
	storage   media.Storage
	audit     *audit.Service
	initError error
}

func NewConfiguredMediaController(auth *AuthController) *MediaController {
	return NewConfiguredMediaControllerWithAudit(auth, nil)
}

func NewConfiguredMediaControllerWithAudit(auth *AuthController, auditService *audit.Service) *MediaController {
	mode, err := resource.ParseResourceProviderMode(facades.Config().GetString("resource.provider", "memory"))
	if err != nil {
		return &MediaController{auth: auth, audit: auditService, initError: err}
	}
	var repository media.Repository
	switch mode {
	case resource.ResourceProviderMemory:
		repository = media.NewMemoryRepository()
	case resource.ResourceProviderMySQL:
		repository = media.NewMySQLRepository()
	default:
		return &MediaController{auth: auth, audit: auditService, initError: errors.New("media provider is not configured")}
	}
	storage := media.NewConfiguredLocalStorage()
	return &MediaController{auth: auth, service: media.NewService(repository, storage), storage: storage, audit: auditService}
}

func (c *MediaController) Index(ctx http.Context) http.Response {
	if response := c.authorize(ctx, "view"); response != nil {
		return response
	}
	items, err := c.service.List(ctx)
	if err != nil {
		return c.storageError(ctx, err)
	}
	return c.success(ctx, items)
}

func (c *MediaController) Store(ctx http.Context) http.Response {
	if response := c.authorize(ctx, "upload"); response != nil {
		return response
	}
	uploaded, err := ctx.Request().File("file")
	if err != nil {
		return c.error(ctx, 400, "INVALID_MEDIA", "请选择要上传的文件", nil)
	}
	file, err := os.Open(uploaded.File())
	if err != nil {
		return c.error(ctx, 400, "INVALID_MEDIA", "上传文件不可读取", nil)
	}
	defer file.Close()
	mimeType, err := uploaded.MimeType()
	if err != nil {
		return c.error(ctx, 400, "INVALID_MEDIA", "上传文件类型不可识别", nil)
	}
	size, err := uploaded.Size()
	if err != nil {
		return c.error(ctx, 400, "INVALID_MEDIA", "上传文件大小不可识别", nil)
	}
	item, err := c.service.Upload(ctx, media.UploadInput{
		OriginalName: uploaded.GetClientOriginalName(), MIMEType: mimeType, Size: size, Content: file,
	})
	if err != nil {
		return c.domainError(ctx, err)
	}
	if err := recordAudit(ctx, c.auth, c.audit, "media.upload", "media", item.ID, nil, item); err != nil {
		return c.storageError(ctx, err)
	}
	return ctx.Response().Status(201).Json(response.Success(item, response.Meta{RequestID: mediaRequestID(ctx)}))
}

func (c *MediaController) Destroy(ctx http.Context) http.Response {
	if response := c.authorize(ctx, "delete"); response != nil {
		return response
	}
	id := strings.TrimSpace(ctx.Request().Route("id"))
	before, _ := c.service.Get(ctx, id)
	if err := c.service.Delete(ctx, id); err != nil {
		return c.domainError(ctx, err)
	}
	if err := recordAudit(ctx, c.auth, c.audit, "media.delete", "media", id, before, nil); err != nil {
		return c.storageError(ctx, err)
	}
	return c.success(ctx, map[string]bool{"deleted": true})
}

func (c *MediaController) Preview(ctx http.Context) http.Response {
	if response := c.authorize(ctx, "view"); response != nil {
		return response
	}
	item, err := c.service.Get(ctx, strings.TrimSpace(ctx.Request().Route("id")))
	if err != nil {
		return c.domainError(ctx, err)
	}
	path := c.storage.Path(item.Path)
	if path == "" {
		return c.storageError(ctx, media.ErrInvalidMedia)
	}
	return ctx.Response().File(path)
}

func (c *MediaController) authorize(ctx http.Context, action string) http.Response {
	if c.initError != nil || c.service == nil || c.storage == nil {
		return c.error(ctx, 503, "MEDIA_PROVIDER_NOT_CONFIGURED", "媒体存储尚未配置", nil)
	}
	if c.auth == nil {
		return c.error(ctx, 401, "UNAUTHENTICATED", "请先登录", nil)
	}
	user, err := c.auth.CurrentUser(ctx)
	if err != nil {
		return c.error(ctx, 401, "UNAUTHENTICATED", "请先登录", nil)
	}
	if permission.NewAuthorizer().Require(user.Permissions, "media."+action) != nil {
		return c.error(ctx, 403, "FORBIDDEN", "没有执行该操作的权限", map[string]string{"permission": "media." + action})
	}
	return nil
}

func (c *MediaController) domainError(ctx http.Context, err error) http.Response {
	switch {
	case errors.Is(err, media.ErrInvalidMedia), errors.Is(err, media.ErrMediaTooLarge):
		return c.error(ctx, 400, "INVALID_MEDIA", "媒体文件无效或超过大小限制", map[string]int64{"max_bytes": media.MaxUploadSize})
	case errors.Is(err, media.ErrNotFound):
		return c.error(ctx, 404, "MEDIA_NOT_FOUND", "媒体文件不存在", nil)
	default:
		return c.storageError(ctx, err)
	}
}

func (c *MediaController) storageError(ctx http.Context, err error) http.Response {
	return c.error(ctx, 503, "MEDIA_STORAGE_UNAVAILABLE", "媒体存储暂不可用", map[string]string{"reason": err.Error()})
}

func (c *MediaController) success(ctx http.Context, data any) http.Response {
	return ctx.Response().Success().Json(response.Success(data, response.Meta{RequestID: mediaRequestID(ctx)}))
}

func (c *MediaController) error(ctx http.Context, status int, code, message string, details any) http.Response {
	return ctx.Response().Status(status).Json(apierrors.New(code, message, details))
}

func mediaRequestID(ctx http.Context) string {
	requestID := ctx.Request().Header("X-Request-ID")
	if requestID == "" {
		requestID = "media-request"
	}
	ctx.Response().Header("X-Request-ID", requestID)
	return requestID
}
