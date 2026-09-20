package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"goravel/app/facades"
)

type M20260920000002CreateRBACTables struct{}

func (m *M20260920000002CreateRBACTables) Signature() string {
	return "20260920000002_create_rbac_tables"
}

func (m *M20260920000002CreateRBACTables) Up() error {
	if !facades.Schema().HasTable("roles") {
		if err := facades.Schema().Create("roles", func(table schema.Blueprint) {
			table.ID()
			table.String("name")
			table.String("display_name")
			table.DateTimeTz("created_at").Nullable()
			table.DateTimeTz("updated_at").Nullable()
			table.Unique("name")
		}); err != nil {
			return err
		}
	}
	if !facades.Schema().HasTable("permissions") {
		if err := facades.Schema().Create("permissions", func(table schema.Blueprint) {
			table.ID()
			table.String("name")
			table.String("display_name")
			table.DateTimeTz("created_at").Nullable()
			table.DateTimeTz("updated_at").Nullable()
			table.Unique("name")
		}); err != nil {
			return err
		}
	}
	if !facades.Schema().HasTable("role_user") {
		if err := facades.Schema().Create("role_user", func(table schema.Blueprint) {
			table.ID()
			table.Integer("role_id")
			table.Integer("user_id")
			table.Unique("role_id", "user_id")
		}); err != nil {
			return err
		}
	}
	if !facades.Schema().HasTable("permission_role") {
		if err := facades.Schema().Create("permission_role", func(table schema.Blueprint) {
			table.ID()
			table.Integer("permission_id")
			table.Integer("role_id")
			table.Unique("permission_id", "role_id")
		}); err != nil {
			return err
		}
	}
	return nil
}

func (m *M20260920000002CreateRBACTables) Down() error {
	for _, table := range []string{"permission_role", "role_user", "permissions", "roles"} {
		if err := facades.Schema().DropIfExists(table); err != nil {
			return err
		}
	}
	return nil
}
