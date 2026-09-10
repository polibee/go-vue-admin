package plugin

import (
	"errors"
	"reflect"
	"testing"
)

type fakePlugin struct {
	manifest PluginManifest
	enables  int
	disables int
}

func (p *fakePlugin) Manifest() PluginManifest { return p.manifest }

func (p *fakePlugin) Enable() error {
	p.enables++
	return nil
}

func (p *fakePlugin) Disable() error {
	p.disables++
	return nil
}

func testManifest(id string) PluginManifest {
	return PluginManifest{
		ID:              id,
		Name:            id + " plugin",
		Version:         "1.0.0",
		Runtime:         RuntimeBuiltin,
		UICompatibility: UICompatibilityShadcnVue,
		Permissions:     []string{id + ".view"},
		Menus: []PluginMenu{{
			ID:         id,
			Label:      id,
			Route:      "/admin/" + id,
			Permission: id + ".view",
		}},
	}
}

func TestRegistryListsBuiltinPluginsInRegistrationOrder(t *testing.T) {
	registry := NewRegistry()
	alpha := &fakePlugin{manifest: testManifest("alpha")}
	beta := &fakePlugin{manifest: testManifest("beta")}

	if err := registry.Register(alpha); err != nil {
		t.Fatal(err)
	}
	if err := registry.Register(beta); err != nil {
		t.Fatal(err)
	}

	items := registry.List()
	if got := []string{items[0].Manifest.ID, items[1].Manifest.ID}; !reflect.DeepEqual(got, []string{"alpha", "beta"}) {
		t.Fatalf("plugin order = %#v", got)
	}
	if items[0].State != PluginDisabled || items[1].State != PluginDisabled {
		t.Fatalf("new plugins should be disabled: %#v", items)
	}
}

func TestRegistryRejectsNonBuiltinAndDuplicatePlugins(t *testing.T) {
	registry := NewRegistry()
	plugin := &fakePlugin{manifest: testManifest("alpha")}
	if err := registry.Register(plugin); err != nil {
		t.Fatal(err)
	}

	duplicate := &fakePlugin{manifest: testManifest("alpha")}
	if err := registry.Register(duplicate); !errors.Is(err, ErrDuplicatePlugin) {
		t.Fatalf("duplicate error = %v", err)
	}

	external := testManifest("external")
	external.Runtime = RuntimeExternal
	if err := registry.Register(&fakePlugin{manifest: external}); !errors.Is(err, ErrBuiltinOnly) {
		t.Fatalf("external runtime error = %v", err)
	}
}

func TestRegistryEnableAndDisablePluginLifecycle(t *testing.T) {
	registry := NewRegistry()
	plugin := &fakePlugin{manifest: testManifest("alpha")}
	if err := registry.Register(plugin); err != nil {
		t.Fatal(err)
	}

	if err := registry.Enable("alpha"); err != nil {
		t.Fatal(err)
	}
	if plugin.enables != 1 || registry.State("alpha") != PluginEnabled {
		t.Fatalf("enable state = %d/%q", plugin.enables, registry.State("alpha"))
	}
	if err := registry.Enable("alpha"); err != nil {
		t.Fatal(err)
	}
	if plugin.enables != 1 {
		t.Fatalf("re-enable called lifecycle %d times", plugin.enables)
	}

	if err := registry.Disable("alpha"); err != nil {
		t.Fatal(err)
	}
	if plugin.disables != 1 || registry.State("alpha") != PluginDisabled {
		t.Fatalf("disable state = %d/%q", plugin.disables, registry.State("alpha"))
	}
}

func TestRegistryRequiresEnabledDependenciesAndProtectsDependents(t *testing.T) {
	registry := NewRegistry()
	alpha := &fakePlugin{manifest: testManifest("alpha")}
	betaManifest := testManifest("beta")
	betaManifest.Dependencies = []PluginDependency{{ID: "alpha", Version: ">=1.0.0"}}
	beta := &fakePlugin{manifest: betaManifest}
	if err := registry.Register(alpha); err != nil {
		t.Fatal(err)
	}
	if err := registry.Register(beta); err != nil {
		t.Fatal(err)
	}

	if err := registry.Enable("beta"); !errors.Is(err, ErrDependencyNotEnabled) {
		t.Fatalf("dependency error = %v", err)
	}
	if err := registry.Enable("alpha"); err != nil {
		t.Fatal(err)
	}
	if err := registry.Enable("beta"); err != nil {
		t.Fatal(err)
	}
	if err := registry.Disable("alpha"); !errors.Is(err, ErrDependencyActive) {
		t.Fatalf("active dependency error = %v", err)
	}
}
