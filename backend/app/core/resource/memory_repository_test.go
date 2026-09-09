package resource

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMemoryDemoResourceRepositoryListFiltersSortsAndPaginates(t *testing.T) {
	repository := NewMemoryDemoResourceRepository(
		DemoResource{ID: "demo-1", Name: "Alpha", Status: "active", Owner: "Platform"},
		DemoResource{ID: "demo-2", Name: "Beta", Status: "draft", Owner: "Platform"},
		DemoResource{ID: "demo-3", Name: "Gamma", Status: "active", Owner: "Other"},
	)

	rows, total, err := repository.List(context.Background(), DemoResourceListQuery{
		Page:      1,
		PerPage:   1,
		Filters:   map[string]string{"status": "active"},
		SortField: "name",
		SortDesc:  true,
	})

	require.NoError(t, err)
	require.Equal(t, 2, total)
	require.Equal(t, []DemoResource{{ID: "demo-3", Name: "Gamma", Status: "active", Owner: "Other"}}, rows)
}

func TestMemoryDemoResourceRepositoryReturnsNotFound(t *testing.T) {
	repository := NewMemoryDemoResourceRepository()

	_, err := repository.Get(context.Background(), "missing")

	require.ErrorIs(t, err, ErrResourceNotFound)
}
