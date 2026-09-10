package bootstrap

import (
	"testing"

	"github.com/stretchr/testify/require"

	"goravel/app/core/module"
)

type bootstrapModule struct {
	name   string
	booted *bool
}

func (m bootstrapModule) Name() string { return m.name }

func (m bootstrapModule) Register(app module.Application) error {
	return app.RegisterService(m.name, m.name)
}

func (m bootstrapModule) Boot(_ module.Application) error {
	*m.booted = true
	return nil
}

func TestBootModulesRunsRegisteredLifecycleInOrder(t *testing.T) {
	booted := false
	registry := module.NewRegistry()
	require.NoError(t, registry.Register(bootstrapModule{name: "example", booted: &booted}))
	app := module.NewApplication()

	require.NoError(t, BootModules(registry, app))
	service, ok := app.Service("example")
	require.True(t, ok)
	require.Equal(t, "example", service)
	require.True(t, booted)
}
