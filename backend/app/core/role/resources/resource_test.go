package resources

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRoleResourcePermissionContract(t *testing.T) {
	require.Equal(t, "roles.view", PermissionView)
	require.Equal(t, "roles.create", PermissionCreate)
	require.Equal(t, "roles.update", PermissionUpdate)
	require.Equal(t, "roles.delete", PermissionDelete)
}
