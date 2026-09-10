package setting

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMemoryRepositoryUpsertsTypedSettingsByNamespaceAndKey(t *testing.T) {
	repository := NewMemoryRepository(Setting{
		Namespace: "general", Key: "site_name", Value: "Go Vue Admin", ValueType: ValueTypeString,
	})
	service := NewService(repository)

	updated, err := service.Upsert(context.Background(), Setting{
		Namespace: "general", Key: "site_name", Value: "Platform Admin", ValueType: ValueTypeString,
	})
	require.NoError(t, err)
	require.Equal(t, "Platform Admin", updated.Value)

	items, err := service.List(context.Background(), "general")
	require.NoError(t, err)
	require.Len(t, items, 1)
	require.Equal(t, "Platform Admin", items[0].Value)
}

func TestServiceValidatesTypedValues(t *testing.T) {
	service := NewService(NewMemoryRepository())

	_, err := service.Upsert(context.Background(), Setting{
		Namespace: "general", Key: "maintenance", Value: "yes", ValueType: ValueTypeBoolean,
	})
	require.ErrorIs(t, err, ErrInvalidValue)

	_, err = service.Upsert(context.Background(), Setting{
		Namespace: "general", Key: "site_name", Value: "name", ValueType: "unsupported",
	})
	require.ErrorIs(t, err, ErrInvalidValueType)
}

func TestServiceRejectsMissingIdentityAndDeletesOneSetting(t *testing.T) {
	service := NewService(NewMemoryRepository(Setting{
		Namespace: "general", Key: "site_name", Value: "Go Vue Admin", ValueType: ValueTypeString,
	}))

	_, err := service.Upsert(context.Background(), Setting{Value: "missing"})
	require.ErrorIs(t, err, ErrInvalidSetting)

	require.NoError(t, service.Delete(context.Background(), "general", "site_name"))
	_, err = service.Get(context.Background(), "general", "site_name")
	require.ErrorIs(t, err, ErrNotFound)
}
