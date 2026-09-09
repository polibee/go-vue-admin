package core_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"goravel/app/core/menu"
	"goravel/app/core/permission"
	"goravel/app/core/role"
	"goravel/app/core/user"
)

func TestUserRolesResolveToResourceActionPermissions(t *testing.T) {
	users := user.NewMemoryRepository(user.User{
		ID:      "user-1",
		RoleIDs: []string{"operator"},
		Active:  true,
	})
	roles := role.NewMemoryRepository(role.Role{
		ID:          "operator",
		Name:        "Operator",
		Permissions: []string{"users.view", "users.update"},
	})

	service := user.NewService(users, roles)
	permissions, err := service.Permissions(context.Background(), "user-1")

	require.NoError(t, err)
	require.ElementsMatch(t, []string{"users.view", "users.update"}, permissions)
}

func TestAuthorizerRejectsForbiddenResourceActions(t *testing.T) {
	authorizer := permission.NewAuthorizer()

	require.NoError(t, authorizer.Require([]string{"users.view"}, "users.view"))
	require.ErrorIs(t, authorizer.Require([]string{"users.view"}, "users.delete"), permission.ErrForbidden)
}

func TestMenuRegistryFiltersItemsWithUnauthorizedPermissions(t *testing.T) {
	registry := menu.NewRegistry(
		menu.Item{ID: "users", Label: "用户", Route: "/admin/users", Permission: "users.view"},
		menu.Item{ID: "settings", Label: "设置", Route: "/admin/settings", Permission: "settings.view"},
	)

	items := registry.ForPermissions([]string{"settings.view"})

	require.Len(t, items, 1)
	require.Equal(t, "settings", items[0].ID)
}
