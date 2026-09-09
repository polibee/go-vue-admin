package controllers

import (
	"context"
	"errors"
	"os"
	"strings"

	"github.com/goravel/framework/contracts/http"
	contractssession "github.com/goravel/framework/contracts/session"
	"github.com/goravel/framework/facades"
	frameworksession "github.com/goravel/framework/session"

	"goravel/app/core/auth"
	apierrors "goravel/app/core/shared/errors"
	"goravel/app/core/shared/response"
)

const SessionUserIDKey = "auth.user_id"

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthController struct {
	service   *auth.Service
	initError error
}

func NewAuthController() *AuthController {
	config := facades.Config()
	email := os.Getenv("AUTH_BOOTSTRAP_EMAIL")
	if email == "" {
		email = config.GetString("auth.bootstrap.email")
	}
	password := os.Getenv("AUTH_BOOTSTRAP_PASSWORD")
	if password == "" {
		password = config.GetString("auth.bootstrap.password")
	}
	name := os.Getenv("AUTH_BOOTSTRAP_NAME")
	if name == "" {
		name = config.GetString("auth.bootstrap.name", "Platform Admin")
	}
	service, err := auth.NewBootstrapService(
		email,
		password,
		name,
	)

	return &AuthController{service: service, initError: err}
}

func NewAuthControllerWithService(service *auth.Service) *AuthController {
	return &AuthController{service: service}
}

func (c *AuthController) CSRF(ctx http.Context) http.Response {
	session := requestSession(ctx)
	if session == nil {
		return c.error(ctx, 503, "AUTH_SESSION_UNAVAILABLE", "会话服务不可用")
	}

	token := session.Token()
	ctx.Response().Header("X-CSRF-TOKEN", token)
	return c.success(ctx, map[string]string{"token": token})
}

func (c *AuthController) Login(ctx http.Context) http.Response {
	if c.initError != nil || c.service == nil {
		return c.error(ctx, 503, "AUTH_NOT_CONFIGURED", "认证服务尚未配置")
	}

	var input LoginRequest
	if err := ctx.Request().Bind(&input); err != nil {
		return c.error(ctx, 400, "INVALID_REQUEST", "登录请求格式无效")
	}
	input.Email = strings.TrimSpace(input.Email)

	user, err := c.service.Authenticate(requestContext(ctx), input.Email, input.Password)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return c.error(ctx, 401, "INVALID_CREDENTIALS", "邮箱或密码错误")
		}
		return c.error(ctx, 503, "AUTH_UNAVAILABLE", "认证服务暂不可用")
	}

	session := requestSession(ctx)
	if session == nil {
		return c.error(ctx, 503, "AUTH_SESSION_UNAVAILABLE", "会话服务不可用")
	}
	if err := session.Regenerate(true); err != nil {
		return c.error(ctx, 500, "SESSION_ROTATION_FAILED", "登录会话创建失败")
	}
	session.Put(SessionUserIDKey, user.ID)
	frameworksession.WriteCookie(ctx, session)
	ctx.Response().Header("X-CSRF-TOKEN", session.Token())

	return c.success(ctx, user.Public())
}

func (c *AuthController) Me(ctx http.Context) http.Response {
	user, err := c.CurrentUser(ctx)
	if err != nil {
		return c.error(ctx, 401, "UNAUTHENTICATED", "请先登录")
	}
	return c.success(ctx, user.Public())
}

func (c *AuthController) Logout(ctx http.Context) http.Response {
	session := requestSession(ctx)
	if session == nil {
		return c.error(ctx, 503, "AUTH_SESSION_UNAVAILABLE", "会话服务不可用")
	}
	if err := session.Invalidate(); err != nil {
		return c.error(ctx, 500, "SESSION_REVOKE_FAILED", "退出登录失败")
	}
	frameworksession.WriteCookie(ctx, session)

	return c.success(ctx, map[string]bool{"logged_out": true})
}

func (c *AuthController) CurrentUser(ctx http.Context) (auth.User, error) {
	if c.initError != nil || c.service == nil {
		return auth.User{}, auth.ErrAuthNotConfigured
	}
	session := requestSession(ctx)
	if session == nil {
		return auth.User{}, auth.ErrUserNotFound
	}
	userID, ok := session.Get(SessionUserIDKey).(string)
	if !ok || userID == "" {
		return auth.User{}, auth.ErrUserNotFound
	}
	return c.service.FindByID(requestContext(ctx), userID)
}

func requestSession(ctx http.Context) contractssession.Session {
	if ctx == nil || ctx.Request() == nil || !ctx.Request().HasSession() {
		return nil
	}
	return ctx.Request().Session()
}

func requestContext(ctx http.Context) context.Context {
	if ctx != nil && ctx.Request() != nil && ctx.Request().Origin() != nil {
		return ctx.Request().Origin().Context()
	}
	return context.Background()
}

func (c *AuthController) success(ctx http.Context, data any) http.Response {
	requestID := requestID(ctx)
	return ctx.Response().Success().Json(response.Success(data, response.Meta{RequestID: requestID}))
}

func (c *AuthController) error(ctx http.Context, status int, code, message string) http.Response {
	requestID := requestID(ctx)
	return ctx.Response().Status(status).Json(apierrors.New(code, message, map[string]string{"request_id": requestID}))
}

func requestID(ctx http.Context) string {
	requestID := ctx.Request().Header("X-Request-ID")
	if requestID == "" {
		requestID = "auth-request"
	}
	ctx.Response().Header("X-Request-ID", requestID)
	return requestID
}
