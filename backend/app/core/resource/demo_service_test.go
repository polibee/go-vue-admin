package resource

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type fakeDemoRepository struct {
	created DemoResource
	updated DemoResource
	deleted []string
	rows    []DemoResource
}

func (f *fakeDemoRepository) List(context.Context, DemoResourceListQuery) ([]DemoResource, int, error) {
	return f.rows, len(f.rows), nil
}

func (f *fakeDemoRepository) Get(context.Context, string) (DemoResource, error) {
	return f.rows[0], nil
}

func (f *fakeDemoRepository) Create(_ context.Context, input DemoResource) (DemoResource, error) {
	f.created = input
	return input, nil
}

func (f *fakeDemoRepository) Update(_ context.Context, _ string, input DemoResourceUpdate) (DemoResource, error) {
	f.updated = DemoResource{Name: stringValue(input.Name), Status: stringValue(input.Status), Owner: stringValue(input.Owner)}
	return f.updated, nil
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func (f *fakeDemoRepository) Delete(_ context.Context, id string) error {
	f.deleted = append(f.deleted, id)
	return nil
}

func (f *fakeDemoRepository) BulkDelete(_ context.Context, ids []string) error {
	f.deleted = append(f.deleted, ids...)
	return nil
}

func TestDemoResourceServiceCreateNormalizesAndValidatesInput(t *testing.T) {
	repository := &fakeDemoRepository{}
	service := NewDemoResourceService(repository)

	created, err := service.Create(context.Background(), DemoResource{
		ID:     "demo-1",
		Name:   "  示例资源  ",
		Status: "active",
		Owner:  " Platform Admin ",
	})

	require.NoError(t, err)
	require.Equal(t, DemoResource{ID: "demo-1", Name: "示例资源", Status: "active", Owner: "Platform Admin"}, created)
	require.Equal(t, created, repository.created)
}

func TestDemoResourceServiceRejectsUnknownStatus(t *testing.T) {
	service := NewDemoResourceService(&fakeDemoRepository{})

	_, err := service.Create(context.Background(), DemoResource{ID: "demo-1", Name: "示例", Status: "paused", Owner: "Admin"})

	require.ErrorIs(t, err, ErrInvalidResource)
}

func TestDemoResourceServiceListPassesThroughQueryAndRows(t *testing.T) {
	repository := &fakeDemoRepository{rows: []DemoResource{{ID: "demo-1", Name: "示例", Status: "active", Owner: "Admin"}}}
	service := NewDemoResourceService(repository)

	result, total, err := service.List(context.Background(), DemoResourceListQuery{Page: 2, PerPage: 10, Search: "示例"})

	require.NoError(t, err)
	require.Equal(t, 1, total)
	require.Equal(t, repository.rows, result)
}

func TestDemoResourceServicePropagatesRepositoryErrors(t *testing.T) {
	repository := &errorDemoRepository{err: errors.New("database unavailable")}
	service := NewDemoResourceService(repository)

	_, _, err := service.List(context.Background(), DemoResourceListQuery{})

	require.EqualError(t, err, "database unavailable")
}

type errorDemoRepository struct{ err error }

func (r *errorDemoRepository) List(context.Context, DemoResourceListQuery) ([]DemoResource, int, error) {
	return nil, 0, r.err
}

func (r *errorDemoRepository) Get(context.Context, string) (DemoResource, error) { return DemoResource{}, r.err }
func (r *errorDemoRepository) Create(context.Context, DemoResource) (DemoResource, error) {
	return DemoResource{}, r.err
}
func (r *errorDemoRepository) Update(context.Context, string, DemoResourceUpdate) (DemoResource, error) {
	return DemoResource{}, r.err
}
func (r *errorDemoRepository) Delete(context.Context, string) error { return r.err }
func (r *errorDemoRepository) BulkDelete(context.Context, []string) error { return r.err }
