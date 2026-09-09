package resource

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseResourceProviderMode(t *testing.T) {
	mode, err := ParseResourceProviderMode("")
	require.NoError(t, err)
	require.Equal(t, ResourceProviderMemory, mode)

	mode, err = ParseResourceProviderMode(" memory ")
	require.NoError(t, err)
	require.Equal(t, ResourceProviderMemory, mode)

	mode, err = ParseResourceProviderMode("MYSQL")
	require.NoError(t, err)
	require.Equal(t, ResourceProviderMySQL, mode)

	_, err = ParseResourceProviderMode("redis")
	require.Error(t, err)
}
