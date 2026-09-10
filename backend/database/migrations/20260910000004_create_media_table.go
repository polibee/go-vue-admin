package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"goravel/app/facades"
)

type M20260910000004CreateMediaTable struct{}

func (m *M20260910000004CreateMediaTable) Signature() string {
	return "20260910000004_create_media_table"
}

func (m *M20260910000004CreateMediaTable) Up() error {
	if facades.Schema().HasTable("media") {
		return nil
	}

	return facades.Schema().Create("media", func(table schema.Blueprint) {
		table.String("id", 36)
		table.String("disk", 32)
		table.String("path", 255)
		table.String("original_name", 255)
		table.String("mime_type", 128)
		table.UnsignedBigInteger("size")
		table.DateTime("created_at").UseCurrent()
		table.Index("mime_type")
		table.Primary("id")
		table.Unique("path")
	})
}

func (m *M20260910000004CreateMediaTable) Down() error {
	return facades.Schema().DropIfExists("media")
}
