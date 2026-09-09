package resource

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

type recordingResourceQuery struct {
	orderField string
	desc       bool
	page       int
	perPage    int
	total      int64
}

func (q *recordingResourceQuery) Where(any, ...any) demoResourceQuery         { return q }
func (q *recordingResourceQuery) WhereLike(string, string) demoResourceQuery  { return q }
func (q *recordingResourceQuery) WhereAny([]string, ...any) demoResourceQuery { return q }
func (q *recordingResourceQuery) WhereIn(string, []any) demoResourceQuery     { return q }
func (q *recordingResourceQuery) OrderBy(field string, _ ...string) demoResourceQuery {
	q.orderField = field
	return q
}
func (q *recordingResourceQuery) OrderByDesc(field string) demoResourceQuery {
	q.orderField = field
	q.desc = true
	return q
}
func (q *recordingResourceQuery) Paginate(page, perPage int, _ any, total *int64) error {
	q.page = page
	q.perPage = perPage
	*total = q.total
	return nil
}
func (q *recordingResourceQuery) Get(any) error    { return nil }
func (q *recordingResourceQuery) Insert(any) error { return nil }
func (q *recordingResourceQuery) Update(any) error { return nil }
func (q *recordingResourceQuery) Delete() error    { return nil }

func TestMySQLDemoResourceRepositoryUsesSafeSortAndPagination(t *testing.T) {
	query := &recordingResourceQuery{total: 4}
	repository := newMySQLDemoResourceRepositoryWithQuery(func(context.Context) demoResourceQuery {
		return query
	})

	_, total, err := repository.List(context.Background(), DemoResourceListQuery{
		Page:      2,
		PerPage:   10,
		SortField: "name;drop table demo_resources",
		SortDesc:  true,
	})

	require.NoError(t, err)
	require.Equal(t, 4, total)
	require.Equal(t, 2, query.page)
	require.Equal(t, 10, query.perPage)
	require.Equal(t, "name", query.orderField)
	require.True(t, query.desc)
}
