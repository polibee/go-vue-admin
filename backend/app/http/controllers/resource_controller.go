package controllers

import (
	"github.com/goravel/framework/contracts/http"

	"goravel/app/resources"
)

type ResourceController struct{}

func NewResourceController() *ResourceController { return &ResourceController{} }

func (r *ResourceController) Index(ctx http.Context) http.Response {
	return ctx.Response().Success().Json(http.Json{"data": resources.AdminRegistry().All()})
}
