package resource

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
)

func TestOpenGormResourceDatabaseOwnsConfiguredConnectionPool(t *testing.T) {
	database, err := OpenGormResourceDatabase(sqlite.Open("file::memory:?cache=shared"), GormDatabaseOptions{
		MaxIdleConns: 3, MaxOpenConns: 7,
	})
	require.NoError(t, err)
	sqlDatabase, err := database.DB.DB()
	require.NoError(t, err)
	require.NoError(t, sqlDatabase.Ping())
	require.NoError(t, database.Close())
}
