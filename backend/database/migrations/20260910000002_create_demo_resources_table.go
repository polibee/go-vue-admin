package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"goravel/app/facades"
)

type M20260910000002CreateDemoResourcesTable struct{}

func (m *M20260910000002CreateDemoResourcesTable) Signature() string {
	return "20260910000002_create_demo_resources_table"
}

func (m *M20260910000002CreateDemoResourcesTable) Up() error {
	if facades.Schema().HasTable("demo_resources") {
		return nil
	}

	return facades.Schema().Create("demo_resources", func(table schema.Blueprint) {
		table.String("id", 64)
		table.String("name", 191)
		table.String("status", 32)
		table.String("owner", 191)
		table.DateTime("created_at").UseCurrent()
		table.DateTime("updated_at").UseCurrent()
		table.Primary("id")
		table.Index("status")
		table.Index("owner")
	})
}

func (m *M20260910000002CreateDemoResourcesTable) Down() error {
	return facades.Schema().DropIfExists("demo_resources")
}
