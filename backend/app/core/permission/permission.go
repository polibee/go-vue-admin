package permission

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrForbidden         = errors.New("permission denied")
	ErrInvalidPermission = errors.New("invalid permission")
)

type Permission struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type Repository interface {
	Find(name string) (Permission, error)
	List() []Permission
}

type MemoryRepository struct {
	permissions map[string]Permission
}

func NewMemoryRepository(permissions ...Permission) *MemoryRepository {
	repository := &MemoryRepository{permissions: make(map[string]Permission, len(permissions))}
	for _, item := range permissions {
		if normalized, err := Normalize(item.Name); err == nil {
			item.Name = normalized
			repository.permissions[normalized] = item
		}
	}
	return repository
}

func (r *MemoryRepository) Find(name string) (Permission, error) {
	normalized, err := Normalize(name)
	if err != nil {
		return Permission{}, err
	}
	item, ok := r.permissions[normalized]
	if !ok {
		return Permission{}, fmt.Errorf("permission %q not found", normalized)
	}
	return item, nil
}

func (r *MemoryRepository) List() []Permission {
	items := make([]Permission, 0, len(r.permissions))
	for _, item := range r.permissions {
		items = append(items, item)
	}
	return items
}

type Authorizer struct{}

func NewAuthorizer() Authorizer { return Authorizer{} }

func (Authorizer) Allows(granted []string, required string) bool {
	required, err := Normalize(required)
	if err != nil {
		return false
	}
	for _, candidate := range granted {
		candidate = strings.TrimSpace(candidate)
		if candidate == required || candidate == "*" {
			return true
		}
		parts := strings.Split(candidate, ".")
		if len(parts) == 2 && parts[1] == "*" && parts[0] == strings.Split(required, ".")[0] {
			return true
		}
	}
	return false
}

func (a Authorizer) Require(granted []string, required string) error {
	if !a.Allows(granted, required) {
		return ErrForbidden
	}
	return nil
}

func Normalize(value string) (string, error) {
	value = strings.TrimSpace(value)
	parts := strings.Split(value, ".")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", fmt.Errorf("%w: %q", ErrInvalidPermission, value)
	}
	return value, nil
}
