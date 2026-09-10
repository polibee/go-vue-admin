package resource

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

type demoRepositoryStub struct {
	row CoreResourceRecord
}

func (s *demoRepositoryStub) List(context.Context, DemoResourceListQuery) ([]CoreResourceRecord, int, error) {
	return []CoreResourceRecord{s.row}, 1, nil
}
func (s *demoRepositoryStub) Get(context.Context, string) (CoreResourceRecord, error) {
	return s.row, nil
}
func (s *demoRepositoryStub) Create(context.Context, CoreResourceRecord) (CoreResourceRecord, error) {
	return s.row, nil
}
func (s *demoRepositoryStub) Update(context.Context, string, CoreResourceRecord) (CoreResourceRecord, error) {
	return s.row, nil
}
func (s *demoRepositoryStub) Delete(context.Context, string) error       { return nil }
func (s *demoRepositoryStub) BulkDelete(context.Context, []string) error { return nil }

func TestGormDemoResourceRepositoryAdaptsTheGenericRepository(t *testing.T) {
	demo := NewGormDemoResourceRepository(&demoRepositoryStub{row: CoreResourceRecord{
		"id": "p-1", "name": "Alpha", "status": "active", "owner": "Admin",
	}})
	row, err := demo.Get(context.Background(), "p-1")
	require.NoError(t, err)
	require.Equal(t, DemoResource{ID: "p-1", Name: "Alpha", Status: "active", Owner: "Admin"}, row)
}
