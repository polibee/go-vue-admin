package module

import "errors"

var (
	ErrInvalidModule   = errors.New("invalid module")
	ErrDuplicateModule = errors.New("module is already registered")
	ErrModuleNotFound  = errors.New("module is not registered")
)

// Application is the narrow registration surface exposed to compiled modules.
// HTTP routing, persistence, and framework bootstrap stay outside this package.
type Application interface {
	RegisterService(name string, service any) error
}

// Module is a statically compiled backend capability with a two-phase lifecycle.
type Module interface {
	Name() string
	Register(app Application) error
	Boot(app Application) error
}
