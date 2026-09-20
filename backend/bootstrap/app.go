package bootstrap

import (
	contractsfoundation "github.com/goravel/framework/contracts/foundation"
	contractsconfiguration "github.com/goravel/framework/contracts/foundation/configuration"
	"github.com/goravel/framework/foundation"
	adminmiddleware "goravel/app/http/middleware"
	"goravel/config"
	"goravel/routes"
)

func Boot() contractsfoundation.Application {
	return foundation.Setup().
		WithSeeders(Seeders).
		WithMigrations(Migrations).
		WithRouting(func() {
			routes.Web()
			routes.Grpc()
		}).
		WithMiddleware(func(middleware contractsconfiguration.Middleware) {
			middleware.Append(adminmiddleware.CORS())
		}).
		WithProviders(Providers).
		WithConfig(config.Boot).
		Create()
}
