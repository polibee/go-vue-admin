package resources

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPermissionResourcePermissionContract(t *testing.T) {
	require.Equal(t, "permissions.view", PermissionView)
	require.Equal(t, "permissions.create", PermissionCreate)
	require.Equal(t, "permissions.update", PermissionUpdate)
	require.Equal(t, "permissions.delete", PermissionDelete)
}
