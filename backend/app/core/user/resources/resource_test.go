package resources

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUserResourcePermissionContract(t *testing.T) {
	require.Equal(t, "users.view", PermissionView)
	require.Equal(t, "users.create", PermissionCreate)
	require.Equal(t, "users.update", PermissionUpdate)
	require.Equal(t, "users.delete", PermissionDelete)
}
