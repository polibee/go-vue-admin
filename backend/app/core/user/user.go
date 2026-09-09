package user

import (
	"context"
	"fmt"

	"goravel/app/core/role"
)

type User struct {
	ID      string   `json:"id"`
	RoleIDs []string `json:"role_ids"`
	Active  bool     `json:"active"`
}

type Repository interface {
	FindByID(ctx context.Context, id string) (User, error)
}

type MemoryRepository struct {
	byID map[string]User
}

func NewMemoryRepository(users ...User) *MemoryRepository {
	repository := &MemoryRepository{byID: make(map[string]User, len(users))}
	for _, item := range users {
		item.RoleIDs = append([]string(nil), item.RoleIDs...)
		repository.byID[item.ID] = item
	}
	return repository
}

func (r *MemoryRepository) FindByID(_ context.Context, id string) (User, error) {
	item, ok := r.byID[id]
	if !ok || !item.Active {
		return User{}, fmt.Errorf("user %q not found", id)
	}
	return item, nil
}

type Service struct {
	users Repository
	roles role.Repository
}

func NewService(users Repository, roles role.Repository) *Service {
	return &Service{users: users, roles: roles}
}

func (s *Service) Permissions(ctx context.Context, userID string) ([]string, error) {
	if s == nil || s.users == nil || s.roles == nil {
		return nil, fmt.Errorf("user authorization is not configured")
	}
	currentUser, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	seen := make(map[string]struct{})
	permissions := make([]string, 0)
	for _, roleID := range currentUser.RoleIDs {
		currentRole, err := s.roles.FindByID(ctx, roleID)
		if err != nil {
			return nil, err
		}
		for _, item := range currentRole.Permissions {
			if _, exists := seen[item]; exists {
				continue
			}
			seen[item] = struct{}{}
			permissions = append(permissions, item)
		}
	}
	return permissions, nil
}
