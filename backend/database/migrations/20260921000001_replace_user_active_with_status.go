package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"goravel/app/facades"
)

type M20260921000001ReplaceUserActiveWithStatus struct{}

func (m *M20260921000001ReplaceUserActiveWithStatus) Signature() string {
	return "20260921000001_replace_user_active_with_status"
}

func (m *M20260921000001ReplaceUserActiveWithStatus) Up() error {
	if !facades.Schema().HasColumn("users", "status") {
		if err := facades.Schema().Table("users", func(table schema.Blueprint) {
			table.String("status").Default("active")
		}); err != nil {
			return err
		}
	}
	if !facades.Schema().HasColumn("users", "is_active") {
		return nil
	}
	if _, err := facades.Orm().Query().Table("users").Where("is_active = ?", false).Update("status", "disabled"); err != nil {
		return err
	}
	return facades.Schema().Table("users", func(table schema.Blueprint) {
		table.DropColumn("is_active")
	})
}

func (m *M20260921000001ReplaceUserActiveWithStatus) Down() error {
	if err := facades.Schema().Table("users", func(table schema.Blueprint) {
		table.Boolean("is_active").Default(true)
	}); err != nil {
		return err
	}
	if _, err := facades.Orm().Query().Table("users").Where("status <> ?", "active").Update("is_active", false); err != nil {
		return err
	}
	return facades.Schema().Table("users", func(table schema.Blueprint) {
		table.DropColumn("status")
	})
}
