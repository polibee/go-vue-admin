package controllers

import (
	"errors"
	"strings"

	"github.com/goravel/framework/contracts/http"

	"goravel/app/facades"
	"goravel/app/models"
	"goravel/app/services"
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
	refreshToken, err := services.NewRefreshTokenService(facades.Cache()).Issue(user.ID)
	if err != nil {
		if errors.Is(err, services.ErrRefreshTokenStoreUnavailable) {
			return sessionStoreUnavailable(ctx)
		}
		return ctx.Response().Status(500).Json(http.Json{
			"code":    "AUTH_REFRESH_TOKEN_ERROR",
			"message": "could not create refresh token",
		})
	}

	response := ctx.Response().Cookie(refreshTokenCookie(refreshToken))
	return response.Success().Json(http.Json{
		"data": http.Json{
			"access_token": token,
			"token_type":   "Bearer",
			"user":         user.Public(),
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

	return ctx.Response().Success().Json(http.Json{"data": user.Public()})
}

func (r *AuthController) Logout(ctx http.Context) http.Response {
	services.NewRefreshTokenService(facades.Cache()).Revoke(ctx.Request().Cookie(services.RefreshTokenCookieName))
	if err := r.parseToken(ctx); err != nil {
		return ctx.Response().WithoutCookie(services.RefreshTokenCookieName).Status(204).Json(nil)
	}
	if err := facades.Auth(ctx).Logout(); err != nil {
		return unauthorized(ctx)
	}

	return ctx.Response().WithoutCookie(services.RefreshTokenCookieName).NoContent(204)
}

func (r *AuthController) Refresh(ctx http.Context) http.Response {
	refreshService := services.NewRefreshTokenService(facades.Cache())
	userID, err := refreshService.Consume(ctx.Request().Cookie(services.RefreshTokenCookieName))
	if err != nil {
		if errors.Is(err, services.ErrRefreshTokenStoreUnavailable) {
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
		if errors.Is(err, services.ErrRefreshTokenStoreUnavailable) {
			return sessionStoreUnavailable(ctx)
		}
		return ctx.Response().Status(500).Json(http.Json{
			"code":    "AUTH_REFRESH_TOKEN_ERROR",
			"message": "could not rotate refresh token",
		})
	}

	return ctx.Response().Cookie(refreshTokenCookie(rotatedToken)).Success().Json(http.Json{
		"data": http.Json{
			"access_token": token,
			"token_type":   "Bearer",
		},
	})
}

func refreshTokenCookie(token string) http.Cookie {
	return http.Cookie{
		Name:     services.RefreshTokenCookieName,
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

func sessionStoreUnavailable(ctx http.Context) http.Response {
	return ctx.Response().Status(503).Json(http.Json{
		"code":    "AUTH_SESSION_STORE_UNAVAILABLE",
		"message": "authentication session storage is unavailable",
	})
}
