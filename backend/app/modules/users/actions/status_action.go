package actions

import (
	"errors"
	"fmt"
	"strings"

	adminactions "goravel/app/core/admin/actions"
	userservices "goravel/app/services/users"

	"github.com/goravel/framework/contracts/http"
)

var (
	ErrUnknownParameter = errors.New("unknown user status action parameter")
	ErrInvalidStatus    = errors.New("invalid user status")
)

type SetStatusHandler struct{}

func NewSetStatusHandler() *SetStatusHandler { return &SetStatusHandler{} }

func (h *SetStatusHandler) Kind() string { return "user-status" }

func ParseStatusParams(params map[string]any) (string, error) {
	if len(params) != 1 {
		for name := range params {
			if name != "status" {
				return "", ErrUnknownParameter
			}
		}
	}
	value, ok := params["status"].(string)
	if !ok {
		return "", ErrInvalidStatus
	}
	status := strings.TrimSpace(value)
	switch status {
	case "active", "disabled", "locked":
		return status, nil
	default:
		return "", ErrInvalidStatus
	}
}

func (h *SetStatusHandler) Execute(_ http.Context, request adminactions.Request) (adminactions.Result, error) {
	status, err := ParseStatusParams(request.Params)
	if err != nil {
		return adminactions.Result{}, err
	}
	result := adminactions.Result{Action: request.Action, Requested: len(request.IDs)}
	for _, id := range request.IDs {
		if err := userservices.NewUserService().BulkSetStatus([]int64{id}, status); err != nil {
			result.Failed++
			result.Failures = append(result.Failures, adminactions.Failure{ID: id, Code: userStatusErrorCode(err)})
			continue
		}
		result.Succeeded++
	}
	return result, nil
}

func userStatusErrorCode(err error) string {
	if errors.Is(err, userservices.ErrInvalidUser) {
		return "VALIDATION_ERROR"
	}
	if errors.Is(err, userservices.ErrUserNotFound) {
		return "RESOURCE_NOT_FOUND"
	}
	if errors.Is(err, userservices.ErrLastAdmin) {
		return "LAST_ADMIN_PROTECTED"
	}
	return fmt.Sprintf("%T", err)
}
