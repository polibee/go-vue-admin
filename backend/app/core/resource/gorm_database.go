package resource

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

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
	applicationResourceDatabase.Lock()
	defer applicationResourceDatabase.Unlock()
	if applicationResourceDatabase.database != nil || applicationResourceDatabase.err != nil {
		return applicationResourceDatabase.err
	}
	database, err := OpenMySQLResourceDatabase(dsn, options)
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

func MySQLResourceDSN(username, password, host, port, database string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=UTC", username, password, host, port, database)
}

func (database *GormResourceDatabase) Close() error {
	if database == nil || database.sqlDB == nil {
		return nil
	}
	return database.sqlDB.Close()
}
