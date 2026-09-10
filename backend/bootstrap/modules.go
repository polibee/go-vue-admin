package bootstrap

import "goravel/app/core/module"

// Modules returns the compiled module registry used by application startup.
// Callers register built-in modules before BootModules is invoked.
func Modules() *module.Registry {
	return module.NewRegistry()
}

func BootModules(registry *module.Registry, app module.Application) error {
	if registry == nil || app == nil {
		return module.ErrInvalidModule
	}
	if err := registry.RegisterAll(app); err != nil {
		return err
	}
	return registry.BootAll(app)
}
