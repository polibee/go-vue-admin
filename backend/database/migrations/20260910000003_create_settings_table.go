package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"goravel/app/facades"
)

type M20260910000003CreateSettingsTable struct{}

func (m *M20260910000003CreateSettingsTable) Signature() string {
	return "20260910000003_create_settings_table"
}

func (m *M20260910000003CreateSettingsTable) Up() error {
	if facades.Schema().HasTable("settings") {
		return nil
	}

	return facades.Schema().Create("settings", func(table schema.Blueprint) {
		table.String("namespace", 64)
		table.String("key", 128)
		table.Text("value")
		table.String("value_type", 16)
		table.String("description", 255).Nullable()
		table.DateTime("created_at").UseCurrent()
		table.DateTime("updated_at").UseCurrent()
		table.Unique("namespace", "key")
		table.Index("namespace")
	})
}

func (m *M20260910000003CreateSettingsTable) Down() error {
	return facades.Schema().DropIfExists("settings")
}
