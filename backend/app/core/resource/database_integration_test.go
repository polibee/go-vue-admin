package resource

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type resourceDatabaseAcceptanceRow struct {
	ID        string `gorm:"primaryKey"`
	Name      string
	Status    string
	DeletedAt *time.Time
}

// TestConfiguredResourceRepositoryIntegration is opt-in so the normal unit
// test suite never requires a developer database. It exercises the selected
// DB_CONNECTION through the same GORM driver used by generated resources.
func TestConfiguredResourceRepositoryIntegration(t *testing.T) {
	if os.Getenv("RESOURCE_DATABASE_INTEGRATION") != "1" {
		t.Skip("set RESOURCE_DATABASE_INTEGRATION=1 to run against a configured database")
	}

	driver, err := ParseResourceDatabaseDriver(os.Getenv("DB_CONNECTION"))
	require.NoError(t, err)
	host := envOr("DB_HOST", "127.0.0.1")
	port := envOr("DB_PORT", map[ResourceDatabaseDriver]string{
		ResourceDatabaseDriverMySQL:    "3306",
		ResourceDatabaseDriverPostgres: "5432",
	}[driver])
	databaseName := envOr("DB_DATABASE", "postgres")
	username := envOr("DB_USERNAME", "postgres")
	password := os.Getenv("DB_PASSWORD")

	dsn := MySQLResourceDSN(username, password, host, port, databaseName)
	if driver == ResourceDatabaseDriverPostgres {
		dsn = PostgresResourceDSN(username, password, host, port, databaseName)
	}
	database, err := OpenResourceDatabase(driver, dsn, GormDatabaseOptions{MaxIdleConns: 2, MaxOpenConns: 4})
	require.NoError(t, err)
	defer database.Close()
	require.NoError(t, database.sqlDB.Ping())

	table := fmt.Sprintf("codex_resource_acceptance_%d", time.Now().UnixNano())
	tableDB := database.DB.Table(table)
	require.NoError(t, tableDB.AutoMigrate(&resourceDatabaseAcceptanceRow{}))
	defer database.DB.Migrator().DropTable(table)

	repository, err := NewGormResourceRepository(GormResourceOptions{
		DB: database.DB,
		Schema: ResourceSchema{
			Table: table, PrimaryKey: "id",
			Fields:     []string{"id", "name", "status", "deleted_at"},
			Searchable: []string{"name", "status"},
			Sortable:   []string{"name", "status"},
			SoftDelete: true,
		},
	})
	require.NoError(t, err)

	ctx := context.Background()
	_, err = repository.Create(ctx, CoreResourceRecord{"id": "accept-1", "name": "Configured DB", "status": "active"})
	require.NoError(t, err)
	rows, total, err := repository.List(ctx, DemoResourceListQuery{Search: "configured", Page: 1, PerPage: 20})
	require.NoError(t, err)
	require.Equal(t, 1, total)
	require.Equal(t, "accept-1", rows[0]["id"])
	updated, err := repository.Update(ctx, "accept-1", CoreResourceRecord{"name": "Updated DB", "status": "ready"})
	require.NoError(t, err)
	require.Equal(t, "Updated DB", updated["name"])
	require.NoError(t, repository.Delete(ctx, "accept-1"))
	_, err = repository.Get(ctx, "accept-1")
	require.ErrorIs(t, err, ErrResourceNotFound)
}

func envOr(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}
