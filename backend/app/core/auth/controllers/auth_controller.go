package controllers

import (
	"errors"
	"strings"

	"github.com/goravel/framework/contracts/http"

	"goravel/app/facades"
	"goravel/app/models"
	auditservices "goravel/app/services/audit"
	authservices "goravel/app/services/auth"
	rbacservices "goravel/app/services/rbac"
)

type AuthController struct{}

func loginAllowedForStatus(status string) bool {
	return status == "active"
}

func NewAuthController() *AuthController {
	return &AuthController{}
}

func (r *AuthController) Login(ctx http.Context) http.Response {
	email := strings.TrimSpace(ctx.Request().Input("email"))
	password := ctx.Request().Input("password")
	if email == "" || password == "" {
		return ctx.Response().Status(422).Json(http.Json{
			"code":    "VALIDATION_ERROR",
			"message": "email and password are required",
		})
	}
	rateLimiter := authservices.NewLoginRateLimiter()
	allowed, err := rateLimiter.Allow(email, ctx.Request().Ip())
	if err != nil {
		return ctx.Response().Status(503).Json(http.Json{"code": "AUTH_RATE_LIMIT_STORE_UNAVAILABLE"})
	}
	if !allowed {
		return ctx.Response().Status(429).Json(http.Json{"code": "AUTH_RATE_LIMITED"})
	}

	var user models.User
	if err := facades.Orm().Query().Where("email = ?", email).First(&user); err != nil || !facades.Hash().Check(password, user.Password) {
		return ctx.Response().Status(401).Json(http.Json{
			"code":    "AUTH_INVALID_CREDENTIALS",
			"message": "invalid credentials",
		})
	}
	if !loginAllowedForStatus(user.Status) {
		return ctx.Response().Status(403).Json(http.Json{
			"code":    "AUTH_USER_DISABLED",
			"message": "user is disabled",
		})
	}

	token, err := facades.Auth(ctx).Login(&user)
	if err != nil {
		return ctx.Response().Status(500).Json(http.Json{
			"code":    "AUTH_TOKEN_ERROR",
			"message": "could not create access token",
		})
	}
	refreshToken, err := authservices.NewDurableRefreshTokenService().Issue(user.ID)
	if err != nil {
		if errors.Is(err, authservices.ErrRefreshTokenStoreUnavailable) {
			return sessionStoreUnavailable(ctx)
		}
		return ctx.Response().Status(500).Json(http.Json{
			"code":    "AUTH_REFRESH_TOKEN_ERROR",
			"message": "could not create refresh token",
		})
	}
	recordAudit(user.ID, "auth.login", map[string]any{"method": "password"})
	rateLimiter.Reset(email, ctx.Request().Ip())
	publicUser, err := authUserPublic(&user)
	if err != nil {
		return ctx.Response().Status(500).Json(http.Json{"code": "RBAC_PERMISSIONS_ERROR"})
	}

	response := ctx.Response().Cookie(refreshTokenCookie(refreshToken))
	return response.Success().Json(http.Json{
		"data": http.Json{
			"access_token": token,
			"token_type":   "Bearer",
			"user":         publicUser,
		},
	})
}

func (r *AuthController) Me(ctx http.Context) http.Response {
	if err := r.parseToken(ctx); err != nil {
		return unauthorized(ctx)
	}

	var user models.User
	if err := facades.Auth(ctx).User(&user); err != nil {
		return unauthorized(ctx)
	}

	publicUser, err := authUserPublic(&user)
	if err != nil {
		return ctx.Response().Status(500).Json(http.Json{"code": "RBAC_PERMISSIONS_ERROR"})
	}
	return ctx.Response().Success().Json(http.Json{"data": publicUser})
}

func (r *AuthController) Logout(ctx http.Context) http.Response {
	authservices.NewDurableRefreshTokenService().Revoke(ctx.Request().Cookie(authservices.RefreshTokenCookieName))
	if err := r.parseToken(ctx); err != nil {
		return ctx.Response().WithoutCookie(authservices.RefreshTokenCookieName).Status(204).Json(nil)
	}
	var user models.User
	if err := facades.Auth(ctx).User(&user); err != nil {
		return unauthorized(ctx)
	}
	response := ctx.Response().WithoutCookie(authservices.RefreshTokenCookieName)
	if err := facades.Auth(ctx).Logout(); err != nil {
		return response.Status(503).Json(http.Json{
			"code":    "AUTH_SESSION_STORE_UNAVAILABLE",
			"message": "authentication session storage is unavailable",
		})
	}
	recordAudit(user.ID, "auth.logout", nil)

	return response.NoContent(204)
}

func (r *AuthController) LogoutAll(ctx http.Context) http.Response {
	if err := r.parseToken(ctx); err != nil {
		return unauthorized(ctx)
	}
	var user models.User
	if err := facades.Auth(ctx).User(&user); err != nil {
		return unauthorized(ctx)
	}
	response := ctx.Response().WithoutCookie(authservices.RefreshTokenCookieName)
	if err := authservices.NewDurableRefreshTokenService().RevokeAll(user.ID); err != nil {
		return response.Status(503).Json(http.Json{
			"code":    "AUTH_SESSION_STORE_UNAVAILABLE",
			"message": "authentication session storage is unavailable",
		})
	}
	if err := facades.Auth(ctx).Logout(); err != nil {
		return response.Status(503).Json(http.Json{
			"code":    "AUTH_SESSION_STORE_UNAVAILABLE",
			"message": "authentication session storage is unavailable",
		})
	}
	recordAudit(user.ID, "auth.logout_all", nil)
	return response.NoContent(204)
}

func (r *AuthController) Refresh(ctx http.Context) http.Response {
	refreshService := authservices.NewDurableRefreshTokenService()
	userID, err := refreshService.Consume(ctx.Request().Cookie(authservices.RefreshTokenCookieName))
	if err != nil {
		if errors.Is(err, authservices.ErrRefreshTokenStoreUnavailable) {
			return sessionStoreUnavailable(ctx)
		}
		return unauthorized(ctx)
	}

	var user models.User
	if err := facades.Orm().Query().Find(&user, userID); err != nil || !loginAllowedForStatus(user.Status) {
		return unauthorized(ctx)
	}
	token, err := facades.Auth(ctx).Login(&user)
	if err != nil {
		return unauthorized(ctx)
	}
	rotatedToken, err := refreshService.Issue(user.ID)
	if err != nil {
		if errors.Is(err, authservices.ErrRefreshTokenStoreUnavailable) {
			return sessionStoreUnavailable(ctx)
		}
		return ctx.Response().Status(500).Json(http.Json{
			"code":    "AUTH_REFRESH_TOKEN_ERROR",
			"message": "could not rotate refresh token",
		})
	}
	recordAudit(user.ID, "auth.refresh", nil)

	return ctx.Response().Cookie(refreshTokenCookie(rotatedToken)).Success().Json(http.Json{
		"data": http.Json{
			"access_token": token,
			"token_type":   "Bearer",
		},
	})
}

func refreshTokenCookie(token string) http.Cookie {
	return http.Cookie{
		Name:     authservices.RefreshTokenCookieName,
		Value:    token,
		Path:     "/api/v1/auth",
		MaxAge:   30 * 24 * 60 * 60,
		SameSite: "Lax",
		Secure:   facades.Config().GetString("app.env", "production") == "production",
		HttpOnly: true,
	}
}

func (r *AuthController) parseToken(ctx http.Context) error {
	return parseAuthToken(ctx)
}

func unauthorized(ctx http.Context) http.Response {
	return ctx.Response().Status(401).Json(http.Json{
		"code":    "AUTH_UNAUTHORIZED",
		"message": "authentication required",
	})
}

func recordAudit(userID uint, action string, metadata map[string]any) {
	if err := auditservices.NewAuditService().Record(userID, action, metadata); err != nil {
		facades.Log().Errorf("audit record failed action=%s user_id=%d error=%v", action, userID, err)
	}
}

func authUserPublic(user *models.User) (map[string]any, error) {
	publicUser := user.Public()
	permissions, err := rbacservices.NewRBACService().PermissionsForUser(user.ID)
	if err != nil {
		return nil, err
	}
	publicUser["permissions"] = permissions
	return publicUser, nil
}

func sessionStoreUnavailable(ctx http.Context) http.Response {
	return ctx.Response().Status(503).Json(http.Json{
		"code":    "AUTH_SESSION_STORE_UNAVAILABLE",
		"message": "authentication session storage is unavailable",
	})
}
