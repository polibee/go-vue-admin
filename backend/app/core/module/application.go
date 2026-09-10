package module

import (
	"fmt"
	"strings"
	"sync"
)

var ErrDuplicateService = fmt.Errorf("module service is already registered")

type ApplicationContainer struct {
	mu       sync.RWMutex
	services map[string]any
}

func NewApplication() *ApplicationContainer {
	return &ApplicationContainer{services: make(map[string]any)}
}

func (a *ApplicationContainer) RegisterService(name string, service any) error {
	if a == nil || strings.TrimSpace(name) == "" || service == nil {
		return ErrInvalidModule
	}
	name = strings.TrimSpace(name)
	a.mu.Lock()
	defer a.mu.Unlock()
	if _, exists := a.services[name]; exists {
		return fmt.Errorf("%w: %s", ErrDuplicateService, name)
	}
	a.services[name] = service
	return nil
}

func (a *ApplicationContainer) Service(name string) (any, bool) {
	if a == nil {
		return nil, false
	}
	a.mu.RLock()
	service, ok := a.services[strings.TrimSpace(name)]
	a.mu.RUnlock()
	return service, ok
}
