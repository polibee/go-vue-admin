package controllers

import (
	"errors"
	"strconv"

	"github.com/goravel/framework/contracts/http"

	"goravel/app/facades"
	notificationservices "goravel/app/services/notifications"
)

type NotificationController struct{}

func NewNotificationController() *NotificationController { return &NotificationController{} }

func normalizeNotificationQuery(page, perPage int) (int, int) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}
	if perPage > 100 {
		perPage = 100
	}
	return page, perPage
}

func notificationUserID(ctx http.Context) (uint, error) {
	identity, err := facades.Auth(ctx).ID()
	if err != nil {
		return 0, err
	}
	parsed, err := strconv.ParseUint(identity, 10, 32)
	if err != nil || parsed == 0 {
		return 0, errors.New("invalid authenticated user")
	}
	return uint(parsed), nil
}

func (n *NotificationController) Index(ctx http.Context) http.Response {
	userID, err := notificationUserID(ctx)
	if err != nil {
		return ctx.Response().Status(401).Json(http.Json{"code": "AUTH_UNAUTHORIZED"})
	}
	page, perPage := normalizeNotificationQuery(positiveInt(ctx.Request().Query("page", "1"), 1), positiveInt(ctx.Request().Query("per_page", "20"), 20))
	unreadOnly := ctx.Request().Query("unread") == "true"
	result, err := notificationservices.NewNotificationService().List(userID, page, perPage, unreadOnly)
	if err != nil {
		return ctx.Response().Status(500).Json(http.Json{"code": "NOTIFICATIONS_ERROR"})
	}
	return ctx.Response().Success().Json(http.Json{
		"data": result.Data,
		"meta": http.Json{"page": result.Page, "per_page": result.PerPage, "total": result.Total, "last_page": result.LastPage},
	})
}

func (n *NotificationController) UnreadCount(ctx http.Context) http.Response {
	userID, err := notificationUserID(ctx)
	if err != nil {
		return ctx.Response().Status(401).Json(http.Json{"code": "AUTH_UNAUTHORIZED"})
	}
	count, err := notificationservices.NewNotificationService().UnreadCount(userID)
	if err != nil {
		return ctx.Response().Status(500).Json(http.Json{"code": "NOTIFICATIONS_ERROR"})
	}
	return ctx.Response().Success().Json(http.Json{"data": http.Json{"count": count}})
}

func (n *NotificationController) MarkRead(ctx http.Context) http.Response {
	userID, err := notificationUserID(ctx)
	if err != nil {
		return ctx.Response().Status(401).Json(http.Json{"code": "AUTH_UNAUTHORIZED"})
	}
	id := uint(positiveInt(ctx.Request().Route("id"), 0))
	if id == 0 {
		return ctx.Response().Status(422).Json(http.Json{"code": "VALIDATION_ERROR"})
	}
	if err := notificationservices.NewNotificationService().MarkRead(userID, id); err != nil {
		if errors.Is(err, notificationservices.ErrNotificationNotFound) {
			return ctx.Response().Status(404).Json(http.Json{"code": "NOTIFICATION_NOT_FOUND"})
		}
		return ctx.Response().Status(500).Json(http.Json{"code": "NOTIFICATIONS_ERROR"})
	}
	return ctx.Response().Success().Json(http.Json{"data": http.Json{"id": id, "read": true}})
}

func (n *NotificationController) MarkAllRead(ctx http.Context) http.Response {
	userID, err := notificationUserID(ctx)
	if err != nil {
		return ctx.Response().Status(401).Json(http.Json{"code": "AUTH_UNAUTHORIZED"})
	}
	updated, err := notificationservices.NewNotificationService().MarkAllRead(userID)
	if err != nil {
		return ctx.Response().Status(500).Json(http.Json{"code": "NOTIFICATIONS_ERROR"})
	}
	return ctx.Response().Success().Json(http.Json{"data": http.Json{"updated": updated}})
}
