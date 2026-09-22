package models

import (
	"time"

	"github.com/goravel/framework/database/orm"
)

type Notification struct {
	orm.Model
	UserID uint       `json:"user_id"`
	Type   string     `json:"type"`
	Title  string     `json:"title"`
	Body   string     `json:"body"`
	URL    string     `json:"url"`
	ReadAt *time.Time `json:"read_at"`
}
