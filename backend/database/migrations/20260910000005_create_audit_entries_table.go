package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"goravel/app/facades"
)

type M20260910000005CreateAuditEntriesTable struct{}

func (m *M20260910000005CreateAuditEntriesTable) Signature() string {
	return "20260910000005_create_audit_entries_table"
}

func (m *M20260910000005CreateAuditEntriesTable) Up() error {
	if facades.Schema().HasTable("audit_entries") {
		return nil
	}

	return facades.Schema().Create("audit_entries", func(table schema.Blueprint) {
		table.String("id", 36)
		table.String("actor_id", 64).Nullable()
		table.String("actor_email", 191).Nullable()
		table.String("action", 128)
		table.String("resource_type", 128)
		table.String("resource_id", 128).Nullable()
		table.LongText("before_state").Nullable()
		table.LongText("after_state").Nullable()
		table.String("ip", 64).Nullable()
		table.String("user_agent", 512).Nullable()
		table.DateTime("created_at").UseCurrent()
		table.Primary("id")
		table.Index("resource_type", "resource_id")
		table.Index("created_at")
	})
}

func (m *M20260910000005CreateAuditEntriesTable) Down() error {
	return facades.Schema().DropIfExists("audit_entries")
}
