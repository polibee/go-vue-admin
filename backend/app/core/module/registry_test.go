package module

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type fakeApplication struct {
	services map[string]any
}

func (a *fakeApplication) RegisterService(name string, service any) error {
	if a.services == nil {
		a.services = make(map[string]any)
	}
	a.services[name] = service
	return nil
}

type fakeModule struct {
	name  string
	order *[]string
}

func (m fakeModule) Name() string { return m.name }

func (m fakeModule) Register(app Application) error {
	*m.order = append(*m.order, "register:"+m.name)
	return app.RegisterService(m.name, m.name)
}

func (m fakeModule) Boot(_ Application) error {
	*m.order = append(*m.order, "boot:"+m.name)
	return nil
}

func TestRegistryPreservesRegistrationOrderAcrossLifecycle(t *testing.T) {
	order := []string{}
	registry := NewRegistry()
	require.NoError(t, registry.Register(fakeModule{name: "alpha", order: &order}))
	require.NoError(t, registry.Register(fakeModule{name: "beta", order: &order}))

	app := &fakeApplication{}
	require.NoError(t, registry.RegisterAll(app))
	require.NoError(t, registry.BootAll(app))
	require.Equal(t, []string{"register:alpha", "register:beta", "boot:alpha", "boot:beta"}, order)
	require.Equal(t, []string{"alpha", "beta"}, registry.Names())
}

func TestRegistryRejectsDuplicateModuleNames(t *testing.T) {
	registry := NewRegistry()
	order := []string{}
	require.NoError(t, registry.Register(fakeModule{name: "alpha", order: &order}))
	require.ErrorIs(t, registry.Register(fakeModule{name: "alpha", order: &order}), ErrDuplicateModule)
}

func TestRegistryRemovalDoesNotAffectOtherModules(t *testing.T) {
	order := []string{}
	registry := NewRegistry()
	require.NoError(t, registry.Register(fakeModule{name: "alpha", order: &order}))
	require.NoError(t, registry.Register(fakeModule{name: "beta", order: &order}))
	require.NoError(t, registry.Remove("alpha"))
	require.Equal(t, []string{"beta"}, registry.Names())

	app := &fakeApplication{}
	require.NoError(t, registry.RegisterAll(app))
	require.NoError(t, registry.BootAll(app))
	require.Equal(t, []string{"register:beta", "boot:beta"}, order)
	require.ErrorIs(t, registry.Remove("alpha"), ErrModuleNotFound)
}

func TestApplicationStoresNamedServicesAndRejectsDuplicates(t *testing.T) {
	app := NewApplication()
	require.NoError(t, app.RegisterService("mailer", "service"))
	require.ErrorIs(t, app.RegisterService("mailer", "other"), ErrDuplicateService)

	service, ok := app.Service("mailer")
	require.True(t, ok)
	require.Equal(t, "service", service)
}
