package controllers

import (
	"strings"

	"github.com/goravel/framework/contracts/http"

	"goravel/app/facades"
	"goravel/app/models"
)

type AuthController struct{}

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
	if !user.IsActive {
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

	return ctx.Response().Success().Json(http.Json{
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
	if err := r.parseToken(ctx); err != nil {
		return unauthorized(ctx)
	}
	if err := facades.Auth(ctx).Logout(); err != nil {
		return unauthorized(ctx)
	}

	return ctx.Response().NoContent(204)
}

func (r *AuthController) Refresh(ctx http.Context) http.Response {
	if err := r.parseToken(ctx); err != nil {
		return unauthorized(ctx)
	}
	token, err := facades.Auth(ctx).Refresh()
	if err != nil {
		return unauthorized(ctx)
	}

	return ctx.Response().Success().Json(http.Json{
		"data": http.Json{
			"access_token": token,
			"token_type":   "Bearer",
		},
	})
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
