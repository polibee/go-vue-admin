package routes

import (
	"github.com/google/uuid"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"goravel/app/core/shared/response"
)

func Health() {
	facades.Route().Get("/api/health", func(ctx http.Context) http.Response {
		requestID := ctx.Request().Header("X-Request-ID", uuid.NewString())
		ctx.Response().Header("X-Request-ID", requestID)
		return ctx.Response().Success().Json(response.Success(map[string]string{
			"status":  "ok",
			"service": "go-vue-admin",
		}, response.Meta{RequestID: requestID}))
	})
}
