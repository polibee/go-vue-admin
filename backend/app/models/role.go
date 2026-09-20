package models

import "github.com/goravel/framework/database/orm"

type Role struct {
	orm.Model
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
}
