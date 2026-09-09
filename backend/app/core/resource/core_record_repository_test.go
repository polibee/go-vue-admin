package resource

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMemoryCoreResourceRepositorySupportsSearchSortAndMutation(t *testing.T) {
	repository := NewMemoryCoreResourceRepository(
		CoreResourceRecord{"id": "user-2", "name": "Beta", "active": true},
		CoreResourceRecord{"id": "user-1", "name": "Alpha", "active": true},
	)

	rows, total, err := repository.List(context.Background(), DemoResourceListQuery{
		Page: 1, PerPage: 1, Search: "alp", SortField: "name",
	})
	require.NoError(t, err)
	require.Equal(t, 1, total)
	require.Equal(t, "user-1", rows[0]["id"])

	created, err := repository.Create(context.Background(), CoreResourceRecord{"id": "user-3", "name": "Gamma"})
	require.NoError(t, err)
	require.Equal(t, "Gamma", created["name"])

	updated, err := repository.Update(context.Background(), "user-3", CoreResourceRecord{"name": "Updated"})
	require.NoError(t, err)
	require.Equal(t, "Updated", updated["name"])
	require.NoError(t, repository.BulkDelete(context.Background(), []string{"user-1", "user-3"}))
	_, err = repository.Get(context.Background(), "user-1")
	require.ErrorIs(t, err, ErrResourceNotFound)
}
