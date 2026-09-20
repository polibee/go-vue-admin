package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"goravel/app/facades"
)

type M20260920000003AddRBACTimestamps struct{}

func (m *M20260920000003AddRBACTimestamps) Signature() string {
	return "20260920000003_add_rbac_timestamps"
}

func (m *M20260920000003AddRBACTimestamps) Up() error {
	for _, tableName := range []string{"roles", "permissions"} {
		if err := facades.Schema().Table(tableName, func(table schema.Blueprint) {
			table.DateTimeTz("created_at").Nullable()
			table.DateTimeTz("updated_at").Nullable()
		}); err != nil {
			return err
		}
	}
	return nil
}

func (m *M20260920000003AddRBACTimestamps) Down() error {
	return nil
}
