package bootstrap

import (
	"fmt"

	"github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/facades"

	"goravel/app/core/resource"
)

type ResourceServiceProvider struct{}

func (p *ResourceServiceProvider) Register(_ foundation.Application) {
	mode, err := resource.ParseResourceProviderMode(facades.Config().GetString("resource.provider", "memory"))
	if err != nil || mode == resource.ResourceProviderMemory {
		return
	}
	config := facades.Config()
	driver, err := resource.ParseResourceDatabaseDriver(config.GetString("database.default", fmt.Sprint(config.Env("DB_CONNECTION", "mysql"))))
	if err != nil {
		_ = resource.ConfigureApplicationResourceDatabaseWithDriver(resource.ResourceDatabaseDriver("invalid"), "", resource.GormDatabaseOptions{})
		return
	}
	username := fmt.Sprint(config.Env("DB_USERNAME"))
	password := fmt.Sprint(config.Env("DB_PASSWORD"))
	host := fmt.Sprint(config.Env("DB_HOST", "127.0.0.1"))
	port := fmt.Sprint(config.Env("DB_PORT", "3306"))
	database := fmt.Sprint(config.Env("DB_DATABASE"))
	dsn := resource.MySQLResourceDSN(username, password, host, port, database)
	if driver == resource.ResourceDatabaseDriverPostgres {
		dsn = resource.PostgresResourceDSN(username, password, host, port, database)
	}
	_ = resource.ConfigureApplicationResourceDatabaseWithDriver(driver, dsn, resource.GormDatabaseOptions{MaxIdleConns: 10, MaxOpenConns: 100})
}

func (p *ResourceServiceProvider) Boot(_ foundation.Application) {}

func (p *ResourceServiceProvider) Runners(_ foundation.Application) []foundation.Runner {
	return []foundation.Runner{resourceDatabaseRunner{}}
}

type resourceDatabaseRunner struct{}

func (resourceDatabaseRunner) Signature() string { return "resource-database" }
func (resourceDatabaseRunner) ShouldRun() bool   { return true }
func (resourceDatabaseRunner) Run() error        { return nil }
func (resourceDatabaseRunner) Shutdown() error   { return resource.CloseApplicationResourceDatabase() }
