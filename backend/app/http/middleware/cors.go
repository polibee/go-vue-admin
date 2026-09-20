package middleware

import (
	httpcontract "github.com/goravel/framework/contracts/http"

	"goravel/app/http/corspolicy"
)

type corsMiddleware struct{}

func (corsMiddleware) Signature() string {
	return "admin:cors"
}

func (corsMiddleware) Handle(ctx httpcontract.Context) {
	origin := ctx.Request().Header("Origin")
	if !corspolicy.AllowedOrigin(origin) {
		ctx.Request().Next()
		return
	}

	ctx.Response().Header("Access-Control-Allow-Origin", origin)
	ctx.Response().Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
	ctx.Response().Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
	ctx.Response().Header("Access-Control-Max-Age", "600")
	ctx.Response().Header("Vary", "Origin")

	if ctx.Request().Method() == "OPTIONS" {
		ctx.Response().NoContent(204).Abort()
		return
	}

	ctx.Request().Next()
}

func CORS() httpcontract.Middleware {
	return corsMiddleware{}
}
