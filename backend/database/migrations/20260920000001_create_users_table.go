package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"goravel/app/facades"
)

type M20260920000001CreateUsersTable struct{}

func (m *M20260920000001CreateUsersTable) Signature() string {
	return "20260920000001_create_users_table"
}

func (m *M20260920000001CreateUsersTable) Up() error {
	if facades.Schema().HasTable("users") {
		return nil
	}

	return facades.Schema().Create("users", func(table schema.Blueprint) {
		table.ID()
		table.String("name")
		table.String("email")
		table.String("password")
		table.Boolean("is_active").Default(true)
		table.String("locale").Default("zh-CN")
		table.DateTimeTz("created_at").Nullable()
		table.DateTimeTz("updated_at").Nullable()
		table.Unique("email")
	})
}

func (m *M20260920000001CreateUsersTable) Down() error {
	return facades.Schema().DropIfExists("users")
}
