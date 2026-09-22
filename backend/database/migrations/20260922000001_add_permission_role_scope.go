package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"goravel/app/facades"
)

type M20260922000001AddPermissionRoleScope struct{}

func (m *M20260922000001AddPermissionRoleScope) Signature() string {
	return "20260922000001_add_permission_role_scope"
}

func (m *M20260922000001AddPermissionRoleScope) Up() error {
	return facades.Schema().Table("permission_role", func(table schema.Blueprint) {
		table.String("scope").Default("all")
	})
}

func (m *M20260922000001AddPermissionRoleScope) Down() error {
	return facades.Schema().Table("permission_role", func(table schema.Blueprint) {
		table.DropColumn("scope")
	})
}
