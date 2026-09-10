package resource

import (
	"database/sql"
	"errors"
	"fmt"
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
