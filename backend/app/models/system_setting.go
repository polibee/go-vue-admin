package models

import "github.com/goravel/framework/database/orm"

type SystemSetting struct {
	orm.Model
	Key         string `json:"key"`
	Value       string `json:"value"`
	ValueType   string `json:"value_type"`
	Group       string `json:"group"`
	Description string `json:"description"`
}
