package audit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestServiceRecordsActorAndChangeSnapshot(t *testing.T) {
	service := NewService(NewMemoryRepository())
	entry, err := service.Record(context.Background(), Entry{
		ActorID: "admin-1", ActorEmail: "admin@example.com", Action: "settings.update",
		ResourceType: "setting", ResourceID: "general.site_name", Before: map[string]any{"value": "Old"}, After: map[string]any{"value": "New"},
		IP: "127.0.0.1", UserAgent: "test-agent",
	})
	require.NoError(t, err)
	require.NotEmpty(t, entry.ID)
	require.NotEmpty(t, entry.CreatedAt)
	require.Equal(t, "settings.update", entry.Action)

	entries, err := service.List(context.Background())
	require.NoError(t, err)
	require.Len(t, entries, 1)
	require.Equal(t, map[string]any{"value": "Old"}, entries[0].Before)
}

func TestServiceRejectsIncompleteEntry(t *testing.T) {
	service := NewService(NewMemoryRepository())
	_, err := service.Record(context.Background(), Entry{Action: "settings.update", ResourceType: "setting"})
	require.ErrorIs(t, err, ErrInvalidEntry)
}
