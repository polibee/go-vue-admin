package bootstrap

import (
	contractsfoundation "github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/contracts/foundation/configuration"
	"github.com/goravel/framework/foundation"
	frameworkhttp "github.com/goravel/framework/http/middleware"
	"github.com/goravel/framework/session/middleware"

	"goravel/config"
	"goravel/routes"
)

func Boot() contractsfoundation.Application {
	return foundation.Setup().
		WithMigrations(Migrations).
		WithMiddleware(func(handler configuration.Middleware) {
			handler.Append(middleware.StartSession(), frameworkhttp.VerifyCsrfToken())
		}).
		WithRouting(func() {
			routes.Web()
			routes.Health()
			routes.Grpc()
		}).
		WithProviders(Providers).
		WithConfig(config.Boot).
		Create()
}
