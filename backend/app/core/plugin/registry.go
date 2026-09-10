package plugin

import (
	"fmt"
	"strings"
	"sync"
)

type registryEntry struct {
	plugin   BuiltinPlugin
	manifest PluginManifest
	state    PluginState
}

type Registry struct {
	mu      sync.RWMutex
	plugins map[string]*registryEntry
	order   []string
}

func NewRegistry() *Registry {
	return &Registry{plugins: make(map[string]*registryEntry)}
}

func (r *Registry) Register(plugin BuiltinPlugin) error {
	if r == nil || plugin == nil {
		return ErrInvalidPlugin
	}
	manifest := cloneManifest(plugin.Manifest())
	if err := validateManifest(manifest); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.plugins[manifest.ID]; exists {
		return fmt.Errorf("%w: %s", ErrDuplicatePlugin, manifest.ID)
	}
	r.plugins[manifest.ID] = &registryEntry{
		plugin:   plugin,
		manifest: manifest,
		state:    PluginDisabled,
	}
	r.order = append(r.order, manifest.ID)
	return nil
}

func (r *Registry) List() []PluginInfo {
	if r == nil {
		return nil
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	items := make([]PluginInfo, 0, len(r.order))
	for _, id := range r.order {
		entry := r.plugins[id]
		items = append(items, PluginInfo{Manifest: cloneManifest(entry.manifest), State: entry.state})
	}
	return items
}

func (r *Registry) State(id string) PluginState {
	if r == nil {
		return ""
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	entry, ok := r.plugins[strings.TrimSpace(id)]
	if !ok {
		return ""
	}
	return entry.state
}

func (r *Registry) Enable(id string) error {
	if r == nil {
		return ErrPluginNotFound
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	entry, ok := r.plugins[strings.TrimSpace(id)]
	if !ok {
		return fmt.Errorf("%w: %s", ErrPluginNotFound, id)
	}
	if entry.state == PluginEnabled {
		return nil
	}
	if err := r.validateDependenciesLocked(entry); err != nil {
		return err
	}
	if err := entry.plugin.Enable(); err != nil {
		return fmt.Errorf("enable plugin %q: %w", entry.manifest.ID, err)
	}
	entry.state = PluginEnabled
	return nil
}

func (r *Registry) Disable(id string) error {
	if r == nil {
		return ErrPluginNotFound
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	entry, ok := r.plugins[strings.TrimSpace(id)]
	if !ok {
		return fmt.Errorf("%w: %s", ErrPluginNotFound, id)
	}
	if entry.state == PluginDisabled {
		return nil
	}
	for _, otherID := range r.order {
		other := r.plugins[otherID]
		if other.state != PluginEnabled {
			continue
		}
		for _, dependency := range other.manifest.Dependencies {
			if dependency.ID == entry.manifest.ID {
				return fmt.Errorf("%w: %s", ErrDependencyActive, other.manifest.ID)
			}
		}
	}
	if err := entry.plugin.Disable(); err != nil {
		return fmt.Errorf("disable plugin %q: %w", entry.manifest.ID, err)
	}
	entry.state = PluginDisabled
	return nil
}

func (r *Registry) validateDependenciesLocked(entry *registryEntry) error {
	for _, dependency := range entry.manifest.Dependencies {
		dependencyEntry, exists := r.plugins[dependency.ID]
		if !exists {
			return fmt.Errorf("%w: %s", ErrDependencyNotFound, dependency.ID)
		}
		if dependencyEntry.state != PluginEnabled {
			return fmt.Errorf("%w: %s", ErrDependencyNotEnabled, dependency.ID)
		}
	}
	return nil
}
