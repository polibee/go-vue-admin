package models

import "github.com/goravel/framework/database/orm"

type User struct {
	orm.Model
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"`
	IsActive bool   `json:"is_active"`
	Locale   string `json:"locale"`
}

func (u User) Public() map[string]any {
	return map[string]any{
		"id":         u.ID,
		"name":       u.Name,
		"email":      u.Email,
		"is_active":  u.IsActive,
		"locale":     u.Locale,
		"created_at": u.CreatedAt,
		"updated_at": u.UpdatedAt,
	}
}
