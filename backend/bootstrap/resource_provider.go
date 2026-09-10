package bootstrap

import (
	"fmt"
	"github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/facades"

	"goravel/app/core/resource"
)

type ResourceServiceProvider struct{}

func (p *ResourceServiceProvider) Register(_ foundation.Application) {
	if facades.Config().GetString("resource.provider", "memory") != string(resource.ResourceProviderMySQL) {
		return
	}
	config := facades.Config()
	dsn := resource.MySQLResourceDSN(
		fmt.Sprint(config.Env("DB_USERNAME")),
		fmt.Sprint(config.Env("DB_PASSWORD")),
		fmt.Sprint(config.Env("DB_HOST", "127.0.0.1")),
		fmt.Sprint(config.Env("DB_PORT", "3306")),
		fmt.Sprint(config.Env("DB_DATABASE")),
	)
	_ = resource.ConfigureApplicationResourceDatabase(dsn, resource.GormDatabaseOptions{MaxIdleConns: 10, MaxOpenConns: 100})
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
