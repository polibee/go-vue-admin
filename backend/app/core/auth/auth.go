package auth

import (
	"context"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
	permissionresources "goravel/app/core/permission/resources"
	roleresources "goravel/app/core/role/resources"
	userresources "goravel/app/core/user/resources"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
	ErrAuthNotConfigured  = errors.New("authentication is not configured")
)

type User struct {
	ID           string   `json:"id"`
	Email        string   `json:"email"`
	Name         string   `json:"name"`
	PasswordHash string   `json:"-"`
	Active       bool     `json:"active"`
	Permissions  []string `json:"permissions,omitempty"`
}

type PublicUser struct {
	ID          string   `json:"id"`
	Email       string   `json:"email"`
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
}

func (u User) Public() PublicUser {
	return PublicUser{ID: u.ID, Email: u.Email, Name: u.Name, Permissions: append([]string(nil), u.Permissions...)}
}

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (User, error)
	FindByID(ctx context.Context, id string) (User, error)
}

type PasswordHasher interface {
	Hash(value string) (string, error)
	Compare(value, hashedValue string) bool
}

type BcryptHasher struct {
	cost int
}

func NewBcryptHasher(cost ...int) *BcryptHasher {
	hashCost := bcrypt.DefaultCost
	if len(cost) > 0 && cost[0] >= bcrypt.MinCost && cost[0] <= bcrypt.MaxCost {
		hashCost = cost[0]
	}

	return &BcryptHasher{cost: hashCost}
}

func (h *BcryptHasher) Hash(value string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(value), h.cost)
	return string(hashed), err
}

func (h *BcryptHasher) Compare(value, hashedValue string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedValue), []byte(value)) == nil
}

type MemoryUserRepository struct {
	byEmail map[string]User
	byID    map[string]User
}

func NewMemoryUserRepository(users ...User) *MemoryUserRepository {
	repository := &MemoryUserRepository{
		byEmail: make(map[string]User, len(users)),
		byID:    make(map[string]User, len(users)),
	}
	for _, user := range users {
		repository.byEmail[strings.ToLower(strings.TrimSpace(user.Email))] = user
		repository.byID[user.ID] = user
	}

	return repository
}

func (r *MemoryUserRepository) FindByEmail(_ context.Context, email string) (User, error) {
	user, ok := r.byEmail[strings.ToLower(strings.TrimSpace(email))]
	if !ok || !user.Active {
		return User{}, ErrUserNotFound
	}
	return user, nil
}

func (r *MemoryUserRepository) FindByID(_ context.Context, id string) (User, error) {
	user, ok := r.byID[id]
	if !ok || !user.Active {
		return User{}, ErrUserNotFound
	}
	return user, nil
}

type Service struct {
	users  UserRepository
	hasher PasswordHasher
}

func NewService(users UserRepository, hasher PasswordHasher) *Service {
	return &Service{users: users, hasher: hasher}
}

func (s *Service) Authenticate(ctx context.Context, email, password string) (User, error) {
	if s == nil || s.users == nil || s.hasher == nil {
		return User{}, ErrAuthNotConfigured
	}

	user, err := s.users.FindByEmail(ctx, email)
	if err != nil || !s.hasher.Compare(password, user.PasswordHash) {
		return User{}, ErrInvalidCredentials
	}

	return user, nil
}

func (s *Service) FindByID(ctx context.Context, id string) (User, error) {
	if s == nil || s.users == nil {
		return User{}, ErrAuthNotConfigured
	}
	return s.users.FindByID(ctx, id)
}

func NewBootstrapService(email, password, name string) (*Service, error) {
	email = strings.TrimSpace(email)
	if email == "" || password == "" {
		return nil, ErrAuthNotConfigured
	}
	if name = strings.TrimSpace(name); name == "" {
		name = "Platform Admin"
	}

	hasher := NewBcryptHasher()
	passwordHash, err := hasher.Hash(password)
	if err != nil {
		return nil, err
	}

	return NewService(
		NewMemoryUserRepository(User{
			ID:           "bootstrap-admin",
			Email:        email,
			Name:         name,
			PasswordHash: passwordHash,
			Active:       true,
			Permissions: []string{
				"dashboard.view",
				"settings.view",
				userresources.PermissionView,
				userresources.PermissionCreate,
				userresources.PermissionUpdate,
				userresources.PermissionDelete,
				roleresources.PermissionView,
				roleresources.PermissionCreate,
				roleresources.PermissionUpdate,
				roleresources.PermissionDelete,
				permissionresources.PermissionView,
				permissionresources.PermissionCreate,
				permissionresources.PermissionUpdate,
				permissionresources.PermissionDelete,
			},
		}),
		hasher,
	), nil
}
