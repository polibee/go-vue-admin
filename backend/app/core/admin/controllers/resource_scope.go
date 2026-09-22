package controllers

import (
	"errors"
	"strconv"
	"strings"

	"github.com/goravel/framework/contracts/database/orm"
	"github.com/goravel/framework/contracts/http"

	"goravel/app/core/resource"
	rbacservices "goravel/app/services/rbac"
)

func resourceID(ctx http.Context) int64 {
	if id := ctx.Request().RouteInt64("id"); id > 0 {
		return id
	}
	parts := strings.Split(strings.Trim(ctx.Request().Path(), "/"), "/")
	if len(parts) == 0 {
		return 0
	}
	id, err := strconv.ParseInt(parts[len(parts)-1], 10, 64)
	if err != nil || id < 1 {
		return 0
	}
	return id
}

func resourceName(ctx http.Context) string {
	if name := ctx.Request().Route("resource"); name != "" {
		return name
	}
	parts := strings.Split(strings.Trim(ctx.Request().Path(), "/"), "/")
	for index, part := range parts {
		if part == "admin" && index+1 < len(parts) {
			return parts[index+1]
		}
	}
	return ""
}

func applyResourceScope(ctx http.Context, query orm.Query, manifest resource.Manifest, action string) (orm.Query, error) {
	return rbacservices.NewResourceScopeService().Apply(ctx, query, manifest, action)
}

func validateResourceActionScope(ctx http.Context, manifest resource.Manifest, action string) error {
	scope, userID, err := rbacservices.NewResourceScopeService().Resolve(ctx, manifest, action)
	if err != nil {
		return err
	}
	_, _, err = rbacservices.OwnerPredicate(manifest, scope, userID)
	return err
}

func resourceScopeError(ctx http.Context, err error) http.Response {
	if errors.Is(err, rbacservices.ErrDataScopeNotAssigned) || errors.Is(err, rbacservices.ErrInvalidOwnerScope) {
		return ctx.Response().Status(403).Json(http.Json{"code": "RBAC_FORBIDDEN"})
	}
	return ctx.Response().Status(500).Json(http.Json{"code": "INTERNAL_ERROR"})
}

func resourceCanAccess(ctx http.Context, manifest resource.Manifest, action string, id int64) (bool, error) {
	return rbacservices.NewResourceScopeService().CanAccess(ctx, manifest, action, id)
}

func enforceResourceCreateOwner(ctx http.Context, manifest resource.Manifest, values map[string]any) error {
	scope, userID, err := rbacservices.NewResourceScopeService().Resolve(ctx, manifest, "create")
	if err != nil {
		return err
	}
	field, _, err := rbacservices.OwnerPredicate(manifest, scope, userID)
	if err != nil {
		return err
	}
	if field != "" {
		values[field] = userID
	}
	return nil
}
