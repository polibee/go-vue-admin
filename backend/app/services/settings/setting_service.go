package settingsservices

import (
	"encoding/json"
	"errors"
	"regexp"
	"strconv"
	"strings"

	"goravel/app/facades"
	"goravel/app/models"
)

var (
	ErrInvalidSetting = errors.New("invalid system setting")
)

var settingKeyPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9._-]{0,159}$`)

type SettingService struct{}

func NewSettingService() *SettingService { return &SettingService{} }

func normalizeSettingType(valueType string) string {
	valueType = strings.ToLower(strings.TrimSpace(valueType))
	if valueType == "" {
		return "string"
	}
	return valueType
}

func validateSetting(key, value, valueType string) error {
	key = strings.TrimSpace(key)
	valueType = normalizeSettingType(valueType)
	if !settingKeyPattern.MatchString(key) {
		return ErrInvalidSetting
	}
	switch valueType {
	case "string":
		return nil
	case "boolean":
		if value == "true" || value == "false" {
			return nil
		}
	case "integer":
		if _, err := strconv.ParseInt(value, 10, 64); err == nil {
			return nil
		}
	case "json":
		if json.Valid([]byte(value)) {
			return nil
		}
	}
	return ErrInvalidSetting
}

func (s *SettingService) List() ([]models.SystemSetting, error) {
	var settings []models.SystemSetting
	if err := facades.Orm().Query().OrderBy("group").OrderBy("key").Get(&settings); err != nil {
		return nil, err
	}
	return settings, nil
}

func (s *SettingService) Upsert(key, value, valueType, group, description string) (*models.SystemSetting, error) {
	key = strings.TrimSpace(key)
	valueType = normalizeSettingType(valueType)
	group = strings.TrimSpace(group)
	if group == "" {
		group = "general"
	}
	if err := validateSetting(key, value, valueType); err != nil {
		return nil, err
	}

	var existing []models.SystemSetting
	if err := facades.Orm().Query().Where("key = ?", key).Get(&existing); err != nil {
		return nil, err
	}
	values := map[string]any{
		"key":         key,
		"value":       value,
		"value_type":  valueType,
		"group":       group,
		"description": strings.TrimSpace(description),
	}
	if len(existing) == 0 {
		if err := facades.Orm().Query().Table("system_settings").Create(&values); err != nil {
			return nil, err
		}
	} else if _, err := facades.Orm().Query().Table("system_settings").Where("key = ?", key).Update(values); err != nil {
		return nil, err
	}

	var setting models.SystemSetting
	if err := facades.Orm().Query().Where("key = ?", key).First(&setting); err != nil {
		return nil, err
	}
	return &setting, nil
}
