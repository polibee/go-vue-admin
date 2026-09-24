package controllers

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/goravel/framework/contracts/http"

	adminactions "goravel/app/core/admin/actions"
	"goravel/app/core/resource"
	"goravel/app/facades"
	adminactionregistry "goravel/app/modules/admin/actions"
	"goravel/app/modules/admin/registry"
	useractions "goravel/app/modules/users/actions"
	notificationservices "goravel/app/services/notifications"
	rbacservices "goravel/app/services/rbac"
	userservices "goravel/app/services/users"
)

type resourceActionPayload struct {
	IDs       []int64                 `json:"ids"`
	Selection *adminactions.Selection `json:"selection"`
	Payload   map[string]any          `json:"payload"`
	Params    map[string]any          `json:"params"`
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
	resourceName := resourceName(ctx)
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
	ids := payload.IDs
	if payload.Selection != nil {
		if err := payload.Selection.Validate(); err != nil {
			return actionValidationError(ctx, err)
		}
		ids, err = resolveSelectionIDs(ctx, manifest, *payload.Selection)
		if err != nil {
			return actionValidationError(ctx, err)
		}
	}
	request, err := adminactions.NormalizeRequest(adminactions.Request{Action: actionName, IDs: ids, Payload: payload.Payload})
	if err != nil {
		return actionValidationError(ctx, err)
	}
	if action.Kind == "builtin-delete" {
		result, executeErr := executeBuiltinDelete(ctx, manifest, request.IDs)
		if executeErr != nil {
			return actionHandlerError(ctx, executeErr)
		}
		recordManagementAudit(ctx, "resource.action", map[string]any{"resource": resourceName, "action": actionName, "requested": result.Requested, "succeeded": result.Succeeded, "failed": result.Failed})
		notifyResourceAction(ctx, resourceName, actionName, result)
		return ctx.Response().Success().Json(http.Json{"data": result})
	}
	if action.Kind == "builtin-update" {
		result, executeErr := executeBuiltinUpdate(ctx, manifest, request)
		if executeErr != nil {
			return actionHandlerError(ctx, executeErr)
		}
		recordManagementAudit(ctx, "resource.action", map[string]any{"resource": resourceName, "action": actionName, "requested": result.Requested, "succeeded": result.Succeeded, "failed": result.Failed})
		notifyResourceAction(ctx, resourceName, actionName, result)
		return ctx.Response().Success().Json(http.Json{"data": result})
	}
	if action.Kind == "builtin-restore" || action.Kind == "builtin-force-delete" {
		result, executeErr := executeBuiltinTrashAction(ctx, manifest, action.Kind, request.IDs)
		if executeErr != nil {
			return actionHandlerError(ctx, executeErr)
		}
		recordManagementAudit(ctx, "resource.action", map[string]any{"resource": resourceName, "action": actionName, "requested": result.Requested, "succeeded": result.Succeeded, "failed": result.Failed})
		notifyResourceAction(ctx, resourceName, actionName, result)
		return ctx.Response().Success().Json(http.Json{"data": result})
	}
	handler, err := adminactionregistry.Registry().FindForPayload(action.Kind, action.Payload)
	if err != nil {
		if errors.Is(err, adminactions.ErrPayloadContract) {
			return actionValidationError(ctx, adminactions.ErrPayloadContract)
		}
		return ctx.Response().Status(500).Json(http.Json{"code": "ACTION_HANDLER_NOT_FOUND"})
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
	notifyResourceAction(ctx, resourceName, actionName, result)
	return ctx.Response().Success().Json(http.Json{"data": result})
}

func notifyResourceAction(ctx http.Context, resourceName, actionName string, result adminactions.Result) {
	identity, err := facades.Auth(ctx).ID()
	if err != nil {
		return
	}
	userID, err := strconv.ParseUint(identity, 10, 32)
	if err != nil || userID == 0 {
		return
	}
	notificationservices.NewNotificationService().PublishBestEffort(uint(userID), notificationservices.NotificationInput{
		Type:  notificationservices.TypeResourceAction,
		Title: "批量操作已完成",
		Body:  fmt.Sprintf("%s / %s：请求 %d 条，成功 %d 条，失败 %d 条。", resourceName, actionName, result.Requested, result.Succeeded, result.Failed),
		URL:   "/admin/" + resourceName,
	})
}

func executeBuiltinUpdate(ctx http.Context, manifest resource.Manifest, request adminactions.Request) (adminactions.Result, error) {
	result := adminactions.Result{Action: request.Action, Requested: len(request.IDs)}
	if len(request.Payload) == 0 {
		return result, adminactions.ErrPayloadContract
	}
	policies, err := resourceFieldPoliciesFor(ctx, manifest, "update")
	if err != nil {
		return result, err
	}
	values, err := rbacservices.NewFieldPermissionService().ValidateWritablePayload(request.Payload, manifest, policies)
	if err != nil {
		return result, err
	}
	for _, id := range request.IDs {
		allowed, accessErr := resourceCanAccess(ctx, manifest, "update", id)
		if accessErr != nil {
			return result, accessErr
		}
		if !allowed {
			result.Skipped++
			result.Skips = append(result.Skips, adminactions.Failure{ID: id, Code: "OUT_OF_SCOPE"})
			continue
		}
		if _, updateErr := facades.Orm().Query().Table(manifest.Table).Where("id = ?", id).Update(values); updateErr != nil {
			result.Failed++
			result.Failures = append(result.Failures, adminactions.Failure{ID: id, Code: "UPDATE_FAILED"})
			continue
		}
		result.Succeeded++
	}
	return result, nil
}

func executeBuiltinDelete(ctx http.Context, manifest resource.Manifest, ids []int64) (adminactions.Result, error) {
	result := adminactions.Result{Action: "bulk-delete", Requested: len(ids)}
	for _, id := range ids {
		allowed, err := resourceCanAccess(ctx, manifest, "delete", id)
		if err != nil {
			return result, err
		}
		if !allowed {
			result.Skipped++
			result.Skips = append(result.Skips, adminactions.Failure{ID: id, Code: "OUT_OF_SCOPE"})
			continue
		}
		var deleteErr error
		switch resourceDeleteStrategy(manifest) {
		case "users-service":
			deleteErr = userservices.NewUserService().Delete(id)
		case "roles-service":
			deleteErr = rbacservices.NewRoleService().Delete(id)
		default:
			_, deleteErr = facades.Orm().Query().Table(manifest.Table).Where("id = ?", id).Delete()
		}
		if deleteErr != nil {
			result.Failed++
			result.Failures = append(result.Failures, adminactions.Failure{ID: id, Code: "DELETE_FAILED"})
			continue
		}
		result.Succeeded++
	}
	return result, nil
}

func findManifestAction(manifest resource.Manifest, name string) (resource.Action, bool) {
	for _, action := range manifest.Actions {
		if action.Name == name {
			return action, true
		}
	}
	if manifest.SoftDelete && (name == "restore" || name == "force-delete") {
		permission := ""
		for _, action := range manifest.Actions {
			if action.Name == "delete" {
				permission = action.Permission
				break
			}
		}
		if permission == "" {
			return resource.Action{}, false
		}
		kind := "builtin-restore"
		if name == "force-delete" {
			kind = "builtin-force-delete"
		}
		return resource.Action{Name: name, Label: name, Kind: kind, Permission: permission, Batch: true, Payload: "trash"}, true
	}
	return resource.Action{}, false
}

func executeBuiltinTrashAction(ctx http.Context, manifest resource.Manifest, kind string, ids []int64) (adminactions.Result, error) {
	result := adminactions.Result{Action: kind, Requested: len(ids)}
	if !manifest.SoftDelete {
		return result, adminactions.ErrActionNotBatch
	}
	for _, id := range ids {
		allowed, err := resourceCanAccess(ctx, manifest, "delete", id)
		if err != nil {
			return result, err
		}
		if !allowed {
			result.Skipped++
			result.Skips = append(result.Skips, adminactions.Failure{ID: id, Code: "OUT_OF_SCOPE"})
			continue
		}
		var actionErr error
		if kind == "builtin-restore" {
			_, actionErr = facades.Orm().Query().Table(manifest.Table).Where("id = ?", id).Update(map[string]any{"deleted_at": nil})
		} else {
			_, actionErr = facades.Orm().Query().Table(manifest.Table).Where("id = ?", id).Delete()
		}
		if actionErr != nil {
			result.Failed++
			result.Failures = append(result.Failures, adminactions.Failure{ID: id, Code: "TRASH_ACTION_FAILED"})
			continue
		}
		result.Succeeded++
	}
	return result, nil
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
	if errors.Is(err, adminactions.ErrEmptyIDs) || errors.Is(err, adminactions.ErrInvalidID) || errors.Is(err, adminactions.ErrTooManyIDs) || errors.Is(err, adminactions.ErrActionNotBatch) || errors.Is(err, adminactions.ErrPayloadContract) || errors.Is(err, adminactions.ErrSelectionContract) {
		return ctx.Response().Status(422).Json(http.Json{"code": "VALIDATION_ERROR"})
	}
	return ctx.Response().Status(500).Json(http.Json{"code": "INTERNAL_ERROR"})
}

func actionHandlerError(ctx http.Context, err error) http.Response {
	if errors.Is(err, useractions.ErrInvalidStatus) || errors.Is(err, useractions.ErrUnknownParameter) || errors.Is(err, rbacservices.ErrFieldPermissionDenied) || errors.Is(err, rbacservices.ErrFieldPolicyExpansion) || errors.Is(err, rbacservices.ErrFieldNotFound) || errors.Is(err, adminactions.ErrPayloadContract) {
		return ctx.Response().Status(422).Json(http.Json{"code": "VALIDATION_ERROR"})
	}
	return ctx.Response().Status(500).Json(http.Json{"code": "INTERNAL_ERROR"})
}
