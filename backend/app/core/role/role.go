package role

import (
	"context"
	"fmt"
)

type Role struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
}

type Repository interface {
	FindByID(ctx context.Context, id string) (Role, error)
}

type MemoryRepository struct {
	byID map[string]Role
}

func NewMemoryRepository(roles ...Role) *MemoryRepository {
	repository := &MemoryRepository{byID: make(map[string]Role, len(roles))}
	for _, item := range roles {
		item.Permissions = append([]string(nil), item.Permissions...)
		repository.byID[item.ID] = item
	}
	return repository
}

func (r *MemoryRepository) FindByID(_ context.Context, id string) (Role, error) {
	item, ok := r.byID[id]
	if !ok {
		return Role{}, fmt.Errorf("role %q not found", id)
	}
	return item, nil
}
