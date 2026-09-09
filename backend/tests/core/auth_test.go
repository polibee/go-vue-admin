package core_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"goravel/app/core/auth"
)

func TestAuthenticationServiceAcceptsValidCredentials(t *testing.T) {
	hasher := auth.NewBcryptHasher()
	hashedPassword, err := hasher.Hash("correct horse battery staple")
	require.NoError(t, err)

	service := auth.NewService(
		auth.NewMemoryUserRepository(auth.User{
			ID:           "user-1",
			Email:        "admin@example.com",
			Name:         "Platform Admin",
			PasswordHash: hashedPassword,
			Active:       true,
		}),
		hasher,
	)

	user, err := service.Authenticate(context.Background(), "admin@example.com", "correct horse battery staple")
	require.NoError(t, err)
	require.Equal(t, "user-1", user.ID)
	require.Equal(t, "admin@example.com", user.Email)
	require.NotEmpty(t, user.PasswordHash)
}

func TestAuthenticationServiceRejectsInvalidCredentialsWithoutLeakingWhichFieldFailed(t *testing.T) {
	hasher := auth.NewBcryptHasher()
	hashedPassword, err := hasher.Hash("correct horse battery staple")
	require.NoError(t, err)

	service := auth.NewService(
		auth.NewMemoryUserRepository(auth.User{
			ID:           "user-1",
			Email:        "admin@example.com",
			Name:         "Platform Admin",
			PasswordHash: hashedPassword,
			Active:       true,
		}),
		hasher,
	)

	_, err = service.Authenticate(context.Background(), "admin@example.com", "wrong-password")
	require.ErrorIs(t, err, auth.ErrInvalidCredentials)

	_, err = service.Authenticate(context.Background(), "missing@example.com", "wrong-password")
	require.ErrorIs(t, err, auth.ErrInvalidCredentials)
}

func TestMemoryUserRepositoryFindsOnlyActiveUsers(t *testing.T) {
	service := auth.NewService(
		auth.NewMemoryUserRepository(
			auth.User{ID: "active", Email: "active@example.com", Active: true},
			auth.User{ID: "disabled", Email: "disabled@example.com", Active: false},
		),
		auth.NewBcryptHasher(),
	)

	user, err := service.FindByID(context.Background(), "active")
	require.NoError(t, err)
	require.Equal(t, "active", user.ID)

	_, err = service.FindByID(context.Background(), "disabled")
	require.ErrorIs(t, err, auth.ErrUserNotFound)
}
