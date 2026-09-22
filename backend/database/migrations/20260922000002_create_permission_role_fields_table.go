package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"goravel/app/facades"
)

type M20260922000002CreatePermissionRoleFieldsTable struct{}

func (m *M20260922000002CreatePermissionRoleFieldsTable) Signature() string {
	return "20260922000002_create_permission_role_fields_table"
}

func (m *M20260922000002CreatePermissionRoleFieldsTable) Up() error {
	if facades.Schema().HasTable("permission_role_field") {
		return nil
	}
	return facades.Schema().Create("permission_role_field", func(table schema.Blueprint) {
		table.ID()
		table.Integer("role_id")
		table.Integer("permission_id")
		table.String("field_name")
		table.Boolean("readable").Default(false)
		table.Boolean("writable").Default(false)
		table.DateTimeTz("created_at").Nullable()
		table.DateTimeTz("updated_at").Nullable()
		table.Unique("role_id", "permission_id", "field_name")
	})
}

func (m *M20260922000002CreatePermissionRoleFieldsTable) Down() error {
	return facades.Schema().DropIfExists("permission_role_field")
}
