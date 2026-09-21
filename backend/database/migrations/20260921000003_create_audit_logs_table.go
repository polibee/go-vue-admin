package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"goravel/app/facades"
)

type M20260921000003CreateAuditLogsTable struct{}

func (m *M20260921000003CreateAuditLogsTable) Signature() string {
	return "20260921000003_create_audit_logs_table"
}

func (m *M20260921000003CreateAuditLogsTable) Up() error {
	if facades.Schema().HasTable("audit_logs") {
		return nil
	}
	return facades.Schema().Create("audit_logs", func(table schema.Blueprint) {
		table.ID()
		table.UnsignedBigInteger("user_id").Nullable()
		table.String("action")
		table.Jsonb("metadata").Nullable()
		table.DateTimeTz("created_at").Nullable()
		table.DateTimeTz("updated_at").Nullable()
		table.Index("user_id")
		table.Index("action")
		table.Index("created_at")
	})
}

func (m *M20260921000003CreateAuditLogsTable) Down() error {
	return facades.Schema().DropIfExists("audit_logs")
}
