package actions

import (
	adminactions "goravel/app/core/admin/actions"
	useractions "goravel/app/modules/users/actions"
)

var applicationRegistry = func() *adminactions.Registry {
	registry := adminactions.NewRegistry()
	_ = registry.Register(useractions.NewSetStatusHandler())
	return registry
}()

// Registry returns the application-level batch Action Handler registry.
// Modules should register custom handlers during application bootstrap.
func Registry() *adminactions.Registry { return applicationRegistry }

// Register adds a module-owned batch Action Handler to the application registry.
func Register(handler adminactions.Handler) error { return applicationRegistry.Register(handler) }
