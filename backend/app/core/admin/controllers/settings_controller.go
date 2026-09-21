package controllers

import (
	"errors"

	"github.com/goravel/framework/contracts/http"

	settingsservices "goravel/app/services/settings"
)

type SettingsController struct{}

func NewSettingsController() *SettingsController { return &SettingsController{} }

type settingPayload struct {
	Value       string `json:"value"`
	ValueType   string `json:"value_type"`
	Group       string `json:"group"`
	Description string `json:"description"`
}

func (s *SettingsController) Index(ctx http.Context) http.Response {
	settings, err := settingsservices.NewSettingService().List()
	if err != nil {
		return ctx.Response().Status(500).Json(http.Json{"code": "SETTINGS_ERROR"})
	}
	return ctx.Response().Success().Json(http.Json{"data": settings})
}

func (s *SettingsController) Upsert(ctx http.Context) http.Response {
	var payload settingPayload
	if err := ctx.Request().Bind(&payload); err != nil {
		return ctx.Response().Status(422).Json(http.Json{"code": "SETTINGS_INVALID"})
	}
	setting, err := settingsservices.NewSettingService().Upsert(ctx.Request().Route("key"), payload.Value, payload.ValueType, payload.Group, payload.Description)
	if err != nil {
		if errors.Is(err, settingsservices.ErrInvalidSetting) {
			return ctx.Response().Status(422).Json(http.Json{"code": "SETTINGS_INVALID"})
		}
		return ctx.Response().Status(500).Json(http.Json{"code": "SETTINGS_ERROR"})
	}
	return ctx.Response().Success().Json(http.Json{"data": setting})
}
