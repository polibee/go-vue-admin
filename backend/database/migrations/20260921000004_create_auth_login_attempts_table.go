package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"goravel/app/facades"
)

type M20260921000004CreateAuthLoginAttemptsTable struct{}

func (m *M20260921000004CreateAuthLoginAttemptsTable) Signature() string {
	return "20260921000004_create_auth_login_attempts_table"
}

func (m *M20260921000004CreateAuthLoginAttemptsTable) Up() error {
	if facades.Schema().HasTable("auth_login_attempts") {
		return nil
	}
	return facades.Schema().Create("auth_login_attempts", func(table schema.Blueprint) {
		table.ID()
		table.String("key_hash", 64)
		table.Integer("attempts")
		table.DateTimeTz("window_started_at")
		table.DateTimeTz("updated_at").Nullable()
		table.Unique("key_hash")
		table.Index("updated_at")
	})
}

func (m *M20260921000004CreateAuthLoginAttemptsTable) Down() error {
	return facades.Schema().DropIfExists("auth_login_attempts")
}
