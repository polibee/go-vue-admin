package module

import (
	"fmt"
	"strings"
	"sync"
)

type Registry struct {
	mu      sync.RWMutex
	modules map[string]Module
	order   []string
}

func NewRegistry() *Registry {
	return &Registry{modules: make(map[string]Module)}
}

func (r *Registry) Register(item Module) error {
	if r == nil || item == nil {
		return ErrInvalidModule
	}
	name := strings.TrimSpace(item.Name())
	if name == "" {
		return ErrInvalidModule
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.modules[name]; exists {
		return fmt.Errorf("%w: %s", ErrDuplicateModule, name)
	}
	r.modules[name] = item
	r.order = append(r.order, name)
	return nil
}

func (r *Registry) Remove(name string) error {
	if r == nil {
		return ErrModuleNotFound
	}
	name = strings.TrimSpace(name)
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.modules[name]; !exists {
		return fmt.Errorf("%w: %s", ErrModuleNotFound, name)
	}
	delete(r.modules, name)
	for index, registered := range r.order {
		if registered == name {
			r.order = append(r.order[:index], r.order[index+1:]...)
			break
		}
	}
	return nil
}

func (r *Registry) Get(name string) (Module, bool) {
	if r == nil {
		return nil, false
	}
	r.mu.RLock()
	item, ok := r.modules[strings.TrimSpace(name)]
	r.mu.RUnlock()
	return item, ok
}

func (r *Registry) Names() []string {
	if r == nil {
		return nil
	}
	r.mu.RLock()
	names := append([]string(nil), r.order...)
	r.mu.RUnlock()
	return names
}

func (r *Registry) RegisterAll(app Application) error {
	for _, item := range r.snapshot() {
		if err := item.Register(app); err != nil {
			return fmt.Errorf("register module %q: %w", item.Name(), err)
		}
	}
	return nil
}

func (r *Registry) BootAll(app Application) error {
	for _, item := range r.snapshot() {
		if err := item.Boot(app); err != nil {
			return fmt.Errorf("boot module %q: %w", item.Name(), err)
		}
	}
	return nil
}

func (r *Registry) snapshot() []Module {
	if r == nil {
		return nil
	}
	r.mu.RLock()
	items := make([]Module, 0, len(r.order))
	for _, name := range r.order {
		items = append(items, r.modules[name])
	}
	r.mu.RUnlock()
	return items
}
