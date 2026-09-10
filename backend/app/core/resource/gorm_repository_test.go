package resource

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestGormResourceRepositorySupportsCRUDFilteringPaginationAndSoftDelete(t *testing.T) {
	database, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	require.NoError(t, err)
	type product struct {
		ID        string `gorm:"primaryKey"`
		Name      string
		Status    string
		DeletedAt *int64
	}
	require.NoError(t, database.AutoMigrate(&product{}))
	require.NoError(t, database.Create(&product{ID: "p-1", Name: "Alpha", Status: "active"}).Error)
	require.NoError(t, database.Create(&product{ID: "p-2", Name: "Beta", Status: "draft"}).Error)
	require.NoError(t, database.Create(&product{ID: "p-3", Name: "Archived", Status: "active"}).Error)

	repository, err := NewGormResourceRepository(GormResourceOptions{
		DB: database,
		Schema: ResourceSchema{
			Table: "products", PrimaryKey: "id",
			Fields:     []string{"id", "name", "status", "deleted_at"},
			Searchable: []string{"name", "status"}, Sortable: []string{"name", "status"}, SoftDelete: true,
		},
	})
	require.NoError(t, err)

	rows, total, err := repository.List(context.Background(), DemoResourceListQuery{Search: "a", Page: 1, PerPage: 1, SortField: "name"})
	require.NoError(t, err)
	require.Equal(t, 3, total)
	require.Len(t, rows, 1)
	require.Equal(t, "Alpha", rows[0]["name"])

	created, err := repository.Create(context.Background(), CoreResourceRecord{"id": "p-4", "name": "Gamma", "status": "active", "ignored": true})
	require.NoError(t, err)
	require.Equal(t, "Gamma", created["name"])

	updated, err := repository.Update(context.Background(), "p-4", CoreResourceRecord{"name": "Gamma Updated", "ignored": true})
	require.NoError(t, err)
	require.Equal(t, "Gamma Updated", updated["name"])

	require.NoError(t, repository.Delete(context.Background(), "p-4"))
	_, err = repository.Get(context.Background(), "p-4")
	require.ErrorIs(t, err, ErrResourceNotFound)

	require.NoError(t, repository.BulkDelete(context.Background(), []string{"p-1", "p-2"}))
	rows, total, err = repository.List(context.Background(), DemoResourceListQuery{Page: 1, PerPage: 20})
	require.NoError(t, err)
	require.Equal(t, 1, total)
	require.Len(t, rows, 1)
}

func TestGormResourceRepositoryRejectsUnsafeSchema(t *testing.T) {
	_, err := NewGormResourceRepository(GormResourceOptions{
		DB:     &gorm.DB{},
		Schema: ResourceSchema{Table: "products;drop", PrimaryKey: "id", Fields: []string{"id"}},
	})
	require.Error(t, err)
}
