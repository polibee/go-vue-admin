package resource

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ResourceDatabaseDriver string

const (
	ResourceDatabaseDriverMySQL    ResourceDatabaseDriver = "mysql"
	ResourceDatabaseDriverPostgres ResourceDatabaseDriver = "postgres"
)

func ParseResourceDatabaseDriver(value string) (ResourceDatabaseDriver, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "", string(ResourceDatabaseDriverMySQL):
		return ResourceDatabaseDriverMySQL, nil
	case string(ResourceDatabaseDriverPostgres), "postgresql", "pgsql":
		return ResourceDatabaseDriverPostgres, nil
	default:
		return "", fmt.Errorf("unsupported resource database driver %q", value)
	}
}

type GormDatabaseOptions struct {
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

type GormResourceDatabase struct {
	DB    *gorm.DB
	sqlDB *sql.DB
}

var applicationResourceDatabase struct {
	sync.RWMutex
	database *GormResourceDatabase
	err      error
}

// ConfigureApplicationResourceDatabase creates the single application-owned
// pool used by generated resources. It is safe to call once during bootstrap.
func ConfigureApplicationResourceDatabase(dsn string, options GormDatabaseOptions) error {
	return ConfigureApplicationResourceDatabaseWithDriver(ResourceDatabaseDriverMySQL, dsn, options)
}

func ConfigureApplicationResourceDatabaseWithDriver(driver ResourceDatabaseDriver, dsn string, options GormDatabaseOptions) error {
	applicationResourceDatabase.Lock()
	defer applicationResourceDatabase.Unlock()
	if applicationResourceDatabase.database != nil || applicationResourceDatabase.err != nil {
		return applicationResourceDatabase.err
	}
	database, err := OpenResourceDatabase(driver, dsn, options)
	if err != nil {
		applicationResourceDatabase.err = err
		return err
	}
	applicationResourceDatabase.database = database
	return nil
}

func ApplicationResourceDatabase() *GormResourceDatabase {
	applicationResourceDatabase.RLock()
	defer applicationResourceDatabase.RUnlock()
	return applicationResourceDatabase.database
}

func ApplicationResourceDatabaseError() error {
	applicationResourceDatabase.RLock()
	defer applicationResourceDatabase.RUnlock()
	return applicationResourceDatabase.err
}

func CloseApplicationResourceDatabase() error {
	applicationResourceDatabase.Lock()
	defer applicationResourceDatabase.Unlock()
	if applicationResourceDatabase.database == nil {
		return nil
	}
	err := applicationResourceDatabase.database.Close()
	applicationResourceDatabase.database = nil
	applicationResourceDatabase.err = nil
	return err
}

func OpenGormResourceDatabase(dialector gorm.Dialector, options GormDatabaseOptions) (*GormResourceDatabase, error) {
	database, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, err
	}
	sqlDatabase, err := database.DB()
	if err != nil {
		return nil, err
	}
	if options.MaxIdleConns > 0 {
		sqlDatabase.SetMaxIdleConns(options.MaxIdleConns)
	}
	if options.MaxOpenConns > 0 {
		sqlDatabase.SetMaxOpenConns(options.MaxOpenConns)
	}
	if options.ConnMaxLifetime > 0 {
		sqlDatabase.SetConnMaxLifetime(options.ConnMaxLifetime)
	}
	if options.ConnMaxIdleTime > 0 {
		sqlDatabase.SetConnMaxIdleTime(options.ConnMaxIdleTime)
	}
	return &GormResourceDatabase{DB: database, sqlDB: sqlDatabase}, nil
}

func OpenMySQLResourceDatabase(dsn string, options GormDatabaseOptions) (*GormResourceDatabase, error) {
	if dsn == "" {
		return nil, errors.New("mysql resource database dsn cannot be empty")
	}
	return OpenGormResourceDatabase(mysql.Open(dsn), options)
}

func OpenPostgresResourceDatabase(dsn string, options GormDatabaseOptions) (*GormResourceDatabase, error) {
	if dsn == "" {
		return nil, errors.New("postgres resource database dsn cannot be empty")
	}
	return OpenGormResourceDatabase(postgres.Open(dsn), options)
}

func OpenResourceDatabase(driver ResourceDatabaseDriver, dsn string, options GormDatabaseOptions) (*GormResourceDatabase, error) {
	switch driver {
	case ResourceDatabaseDriverMySQL:
		return OpenMySQLResourceDatabase(dsn, options)
	case ResourceDatabaseDriverPostgres:
		return OpenPostgresResourceDatabase(dsn, options)
	default:
		return nil, fmt.Errorf("unsupported resource database driver %q", driver)
	}
}

func MySQLResourceDSN(username, password, host, port, database string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=UTC", username, password, host, port, database)
}

func PostgresResourceDSN(username, password, host, port, database string) string {
	return PostgresResourceDSNWithSSLMode(username, password, host, port, database, "disable")
}

func PostgresResourceDSNWithSSLMode(username, password, host, port, database, sslMode string) string {
	if strings.TrimSpace(sslMode) == "" { sslMode = "disable" }
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=UTC", host, port, username, password, database, sslMode)
}

func (database *GormResourceDatabase) Close() error {
	if database == nil || database.sqlDB == nil {
		return nil
	}
	return database.sqlDB.Close()
}
