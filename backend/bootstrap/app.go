package bootstrap

import (
	contractsconsole "github.com/goravel/framework/contracts/console"
	contractsfoundation "github.com/goravel/framework/contracts/foundation"
	contractsconfiguration "github.com/goravel/framework/contracts/foundation/configuration"
	"github.com/goravel/framework/foundation"
	adminconsole "goravel/app/console"
	adminmiddleware "goravel/app/http/middleware"
	"goravel/config"
	"goravel/routes"
)

func Boot() contractsfoundation.Application {
	return foundation.Setup().
		WithSeeders(Seeders).
		WithMigrations(Migrations).
		WithSchedule(Schedule).
		WithCommands(func() []contractsconsole.Command {
			return []contractsconsole.Command{
				adminconsole.ResourceGeneratorCommand{},
				adminconsole.ModuleCheckCommand{},
				adminconsole.AuditPruneCommand{},
			}
		}).
		WithRouting(func() {
			routes.Web()
			routes.Grpc()
		}).
		WithMiddleware(func(middleware contractsconfiguration.Middleware) {
			middleware.Append(adminmiddleware.CORS())
			middleware.Append(adminmiddleware.HTTPAudit())
		}).
		WithProviders(Providers).
		WithConfig(config.Boot).
		Create()
}
