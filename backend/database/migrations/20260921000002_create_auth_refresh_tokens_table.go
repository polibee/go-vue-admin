package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"goravel/app/facades"
)

type M20260921000002CreateAuthRefreshTokensTable struct{}

func (m *M20260921000002CreateAuthRefreshTokensTable) Signature() string {
	return "20260921000002_create_auth_refresh_tokens_table"
}

func (m *M20260921000002CreateAuthRefreshTokensTable) Up() error {
	if facades.Schema().HasTable("auth_refresh_tokens") {
		return nil
	}
	return facades.Schema().Create("auth_refresh_tokens", func(table schema.Blueprint) {
		table.ID()
		table.String("token_hash", 128)
		table.UnsignedBigInteger("user_id")
		table.DateTimeTz("expires_at")
		table.DateTimeTz("created_at").Nullable()
		table.DateTimeTz("updated_at").Nullable()
		table.Unique("token_hash")
		table.Index("user_id")
		table.Index("expires_at")
	})
}

func (m *M20260921000002CreateAuthRefreshTokensTable) Down() error {
	return facades.Schema().DropIfExists("auth_refresh_tokens")
}
