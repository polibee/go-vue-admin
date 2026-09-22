package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"goravel/app/facades"
)

type M20260922000005CreateNotificationsTable struct{}

func (m *M20260922000005CreateNotificationsTable) Signature() string {
	return "20260922000005_create_notifications_table"
}

func (m *M20260922000005CreateNotificationsTable) Up() error {
	if facades.Schema().HasTable("notifications") {
		return nil
	}
	return facades.Schema().Create("notifications", func(table schema.Blueprint) {
		table.ID()
		table.UnsignedBigInteger("user_id")
		table.String("type", 80)
		table.String("title", 255)
		table.Text("body")
		table.String("url", 500).Nullable()
		table.DateTimeTz("read_at").Nullable()
		table.DateTimeTz("created_at").Nullable()
		table.DateTimeTz("updated_at").Nullable()
		table.Index("user_id")
		table.Index("read_at")
		table.Index("created_at")
	})
}

func (m *M20260922000005CreateNotificationsTable) Down() error {
	return facades.Schema().DropIfExists("notifications")
}
