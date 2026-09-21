package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"goravel/app/facades"
)

type M20260921000005CreateSystemSettingsTable struct{}

func (m *M20260921000005CreateSystemSettingsTable) Signature() string {
	return "20260921000005_create_system_settings_table"
}

func (m *M20260921000005CreateSystemSettingsTable) Up() error {
	if facades.Schema().HasTable("system_settings") {
		return nil
	}
	return facades.Schema().Create("system_settings", func(table schema.Blueprint) {
		table.ID()
		table.String("key", 160)
		table.Text("value")
		table.String("value_type", 32).Default("string")
		table.String("group", 80).Default("general")
		table.String("description", 255).Nullable()
		table.DateTimeTz("created_at").Nullable()
		table.DateTimeTz("updated_at").Nullable()
		table.Unique("key")
		table.Index("group")
	})
}

func (m *M20260921000005CreateSystemSettingsTable) Down() error {
	return facades.Schema().DropIfExists("system_settings")
}
