package bootstrap

import (
	"github.com/goravel/framework/contracts/database/schema"

	"goravel/database/migrations"
)

func Migrations() []schema.Migration {
	return []schema.Migration{
		&migrations.M20210101000001CreateJobsTable{},
		&migrations.M20260920000001CreateUsersTable{},
		&migrations.M20260920000002CreateRBACTables{},
		&migrations.M20260920000003AddRBACTimestamps{},
		&migrations.M20260921000001ReplaceUserActiveWithStatus{},
		&migrations.M20260921000002CreateAuthRefreshTokensTable{},
		&migrations.M20260921000003CreateAuditLogsTable{},
		&migrations.M20260921000004CreateAuthLoginAttemptsTable{},
	}
}
