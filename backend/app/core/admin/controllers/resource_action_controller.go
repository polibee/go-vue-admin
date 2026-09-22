package controllers

import (
	"errors"

	"github.com/goravel/framework/contracts/http"

	adminactions "goravel/app/core/admin/actions"
	"goravel/app/core/resource"
	"goravel/app/modules/admin/registry"
	useractions "goravel/app/modules/users/actions"
	rbacservices "goravel/app/services/rbac"
)

type resourceActionPayload struct {
	IDs     []int64        `json:"ids"`
	Payload map[string]any `json:"payload"`
	Params  map[string]any `json:"params"`
}

type actionFailureInput struct {
	ID   int64
	Code string
}

type actionResultInput struct {
	Succeeded int
	Failures  []adminactions.Failure
	Skips     []adminactions.Failure
}

func (r *ResourceController) Action(ctx http.Context) http.Response {
	resourceName := ctx.Request().Route("resource")
	actionName := ctx.Request().Route("action")
	manifest, err := registry.AdminRegistry().Find(resourceName)
	if err != nil || manifest.Table == "" {
		return ctx.Response().Status(404).Json(http.Json{"code": "RESOURCE_NOT_FOUND"})
	}
	action, ok := findManifestAction(manifest, actionName)
	if !ok || action.Permission == "" {
		return ctx.Response().Status(404).Json(http.Json{"code": "ACTION_NOT_FOUND"})
	}
	if !action.Batch || action.Kind == "" || action.Payload == "" {
		return actionValidationError(ctx, adminactions.ErrActionNotBatch)
	}
	allowed, err := rbacservices.NewRBACService().UserHasPermission(ctx, action.Permission)
	if err != nil {
		return ctx.Response().Status(401).Json(http.Json{"code": "AUTH_UNAUTHORIZED"})
	}
	if !allowed {
		return ctx.Response().Status(403).Json(http.Json{"code": "RBAC_FORBIDDEN"})
	}
	var payload resourceActionPayload
	if err := ctx.Request().Bind(&payload); err != nil {
		return ctx.Response().Status(422).Json(http.Json{"code": "VALIDATION_ERROR"})
	}
	if payload.Params != nil {
		return actionValidationError(ctx, adminactions.ErrPayloadContract)
	}
	request, err := adminactions.NormalizeRequest(adminactions.Request{Action: actionName, IDs: payload.IDs, Payload: payload.Payload})
	if err != nil {
		return actionValidationError(ctx, err)
	}
	handler, err := newAdminActionRegistry().Find(action.Kind)
	if err != nil {
		return ctx.Response().Status(500).Json(http.Json{"code": "ACTION_HANDLER_NOT_FOUND"})
	}
	if handler.Payload() != action.Payload {
		return actionValidationError(ctx, adminactions.ErrPayloadContract)
	}
	accessibleIDs := make([]int64, 0, len(request.IDs))
	scopeFailures := make([]actionFailureInput, 0)
	for _, id := range request.IDs {
		canAccess, scopeErr := resourceCanAccess(ctx, manifest, actionName, id)
		if scopeErr != nil {
			return resourceScopeError(ctx, scopeErr)
		}
		if !canAccess {
			scopeFailures = append(scopeFailures, actionFailureInput{ID: id, Code: "OUT_OF_SCOPE"})
			continue
		}
		accessibleIDs = append(accessibleIDs, id)
	}
	request.IDs = accessibleIDs
	handlerResult := actionResultInput{}
	if len(accessibleIDs) > 0 {
		result, executeErr := handler.Execute(ctx, request)
		if executeErr != nil {
			return actionHandlerError(ctx, executeErr)
		}
		handlerResult.Succeeded = result.Succeeded
		handlerResult.Failures = result.Failures
		handlerResult.Skips = result.Skips
	}
	result := mergeActionResults(actionName, len(request.IDs)+len(scopeFailures), scopeFailures, handlerResult)
	recordManagementAudit(ctx, "resource.action", map[string]any{"resource": resourceName, "action": actionName, "requested": result.Requested, "succeeded": result.Succeeded, "failed": result.Failed})
	return ctx.Response().Success().Json(http.Json{"data": result})
}

func findManifestAction(manifest resource.Manifest, name string) (resource.Action, bool) {
	for _, action := range manifest.Actions {
		if action.Name == name {
			return action, true
		}
	}
	return resource.Action{}, false
}

func newAdminActionRegistry() *adminactions.Registry {
	registry := adminactions.NewRegistry()
	_ = registry.Register(useractions.NewSetStatusHandler())
	return registry
}

func mergeActionResults(action string, requested int, scopeFailures []actionFailureInput, handlerResult actionResultInput) adminactions.Result {
	failures := append([]adminactions.Failure(nil), handlerResult.Failures...)
	skips := make([]adminactions.Failure, 0, len(scopeFailures)+len(handlerResult.Skips))
	for _, skip := range scopeFailures {
		skips = append(skips, adminactions.Failure{ID: skip.ID, Code: skip.Code})
	}
	skips = append(skips, handlerResult.Skips...)
	return adminactions.Result{Action: action, Requested: requested, Succeeded: handlerResult.Succeeded, Failed: len(failures), Skipped: len(skips), Failures: failures, Skips: skips}
}

func actionValidationError(ctx http.Context, err error) http.Response {
	if errors.Is(err, adminactions.ErrEmptyIDs) || errors.Is(err, adminactions.ErrInvalidID) || errors.Is(err, adminactions.ErrTooManyIDs) || errors.Is(err, adminactions.ErrActionNotBatch) || errors.Is(err, adminactions.ErrPayloadContract) {
		return ctx.Response().Status(422).Json(http.Json{"code": "VALIDATION_ERROR"})
	}
	return ctx.Response().Status(500).Json(http.Json{"code": "INTERNAL_ERROR"})
}

func actionHandlerError(ctx http.Context, err error) http.Response {
	if errors.Is(err, useractions.ErrInvalidStatus) || errors.Is(err, useractions.ErrUnknownParameter) {
		return ctx.Response().Status(422).Json(http.Json{"code": "VALIDATION_ERROR"})
	}
	return ctx.Response().Status(500).Json(http.Json{"code": "INTERNAL_ERROR"})
}
