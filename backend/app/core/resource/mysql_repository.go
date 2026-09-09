package resource

import (
	"context"
	"strings"
	"time"

	"github.com/goravel/framework/contracts/database/db"
	"github.com/goravel/framework/facades"
)

const demoResourcesTable = "demo_resources"

type demoResourceQuery interface {
	Where(any, ...any) demoResourceQuery
	WhereLike(string, string) demoResourceQuery
	WhereAny([]string, ...any) demoResourceQuery
	WhereIn(string, []any) demoResourceQuery
	OrderBy(string, ...string) demoResourceQuery
	OrderByDesc(string) demoResourceQuery
	Paginate(int, int, any, *int64) error
	Get(any) error
	Insert(any) error
	Update(any) error
	Delete() error
}

type demoResourceQueryFactory func(context.Context) demoResourceQuery

type MySQLDemoResourceRepository struct {
	query demoResourceQueryFactory
}

func NewMySQLDemoResourceRepository() DemoResourceRepository {
	return newMySQLDemoResourceRepositoryWithQuery(func(ctx context.Context) demoResourceQuery {
		return &goravelDemoResourceQuery{query: facades.DB().WithContext(ctx).Table(demoResourcesTable)}
	})
}

func newMySQLDemoResourceRepositoryWithQuery(factory demoResourceQueryFactory) *MySQLDemoResourceRepository {
	return &MySQLDemoResourceRepository{query: factory}
}

func (r *MySQLDemoResourceRepository) List(ctx context.Context, query DemoResourceListQuery) ([]DemoResource, int, error) {
	databaseQuery := r.query(ctx)
	search := strings.TrimSpace(query.Search)
	if search != "" {
		pattern := "%" + search + "%"
		databaseQuery = databaseQuery.WhereAny([]string{"id", "name", "status", "owner"}, "like", pattern)
	}
	for field, value := range query.Filters {
		if isSortableResourceField(field) && strings.TrimSpace(value) != "" {
			databaseQuery = databaseQuery.Where(field, strings.TrimSpace(value))
		}
	}
	sortField := query.SortField
	if !isSortableResourceField(sortField) {
		sortField = "name"
	}
	if query.SortDesc {
		databaseQuery = databaseQuery.OrderByDesc(sortField)
	} else {
		databaseQuery = databaseQuery.OrderBy(sortField)
	}

	var rows []demoResourceRow
	var total int64
	if err := databaseQuery.Paginate(query.Page, query.PerPage, &rows, &total); err != nil {
		return nil, 0, err
	}
	resources := make([]DemoResource, 0, len(rows))
	for _, row := range rows {
		resources = append(resources, row.resource())
	}
	return resources, int(total), nil
}

func (r *MySQLDemoResourceRepository) Get(ctx context.Context, id string) (DemoResource, error) {
	var rows []demoResourceRow
	if err := r.query(ctx).Where("id", id).Get(&rows); err != nil {
		return DemoResource{}, err
	}
	if len(rows) == 0 {
		return DemoResource{}, ErrResourceNotFound
	}
	return rows[0].resource(), nil
}

func (r *MySQLDemoResourceRepository) Create(ctx context.Context, input DemoResource) (DemoResource, error) {
	if err := r.query(ctx).Insert(map[string]any{
		"id":     input.ID,
		"name":   input.Name,
		"status": input.Status,
		"owner":  input.Owner,
	}); err != nil {
		return DemoResource{}, err
	}
	return r.Get(ctx, input.ID)
}

func (r *MySQLDemoResourceRepository) Update(ctx context.Context, id string, input DemoResourceUpdate) (DemoResource, error) {
	values := make(map[string]any, 3)
	if input.Name != nil {
		values["name"] = *input.Name
	}
	if input.Status != nil {
		values["status"] = *input.Status
	}
	if input.Owner != nil {
		values["owner"] = *input.Owner
	}
	if err := r.query(ctx).Where("id", id).Update(values); err != nil {
		return DemoResource{}, err
	}
	return r.Get(ctx, id)
}

func (r *MySQLDemoResourceRepository) Delete(ctx context.Context, id string) error {
	if _, err := r.Get(ctx, id); err != nil {
		return err
	}
	return r.query(ctx).Where("id", id).Delete()
}

func (r *MySQLDemoResourceRepository) BulkDelete(ctx context.Context, ids []string) error {
	values := make([]any, 0, len(ids))
	for _, id := range ids {
		values = append(values, id)
	}
	return r.query(ctx).WhereIn("id", values).Delete()
}

type demoResourceRow struct {
	ID        string    `mapstructure:"id"`
	Name      string    `mapstructure:"name"`
	Status    string    `mapstructure:"status"`
	Owner     string    `mapstructure:"owner"`
	CreatedAt time.Time `mapstructure:"created_at"`
	UpdatedAt time.Time `mapstructure:"updated_at"`
}

func (r demoResourceRow) resource() DemoResource {
	return DemoResource{
		ID:        r.ID,
		Name:      r.Name,
		Status:    r.Status,
		Owner:     r.Owner,
		CreatedAt: formatResourceTime(r.CreatedAt),
		UpdatedAt: formatResourceTime(r.UpdatedAt),
	}
}

func formatResourceTime(value time.Time) string {
	if value.IsZero() {
		return ""
	}
	return value.UTC().Format(time.RFC3339)
}

type goravelDemoResourceQuery struct {
	query db.Query
}

func (q *goravelDemoResourceQuery) Where(value any, args ...any) demoResourceQuery {
	return &goravelDemoResourceQuery{query: q.query.Where(value, args...)}
}

func (q *goravelDemoResourceQuery) WhereLike(column, value string) demoResourceQuery {
	return &goravelDemoResourceQuery{query: q.query.WhereLike(column, value)}
}

func (q *goravelDemoResourceQuery) WhereAny(columns []string, args ...any) demoResourceQuery {
	return &goravelDemoResourceQuery{query: q.query.WhereAny(columns, args...)}
}

func (q *goravelDemoResourceQuery) WhereIn(column string, values []any) demoResourceQuery {
	return &goravelDemoResourceQuery{query: q.query.WhereIn(column, values)}
}

func (q *goravelDemoResourceQuery) OrderBy(column string, directions ...string) demoResourceQuery {
	return &goravelDemoResourceQuery{query: q.query.OrderBy(column, directions...)}
}

func (q *goravelDemoResourceQuery) OrderByDesc(column string) demoResourceQuery {
	return &goravelDemoResourceQuery{query: q.query.OrderByDesc(column)}
}

func (q *goravelDemoResourceQuery) Paginate(page, limit int, dest any, total *int64) error {
	return q.query.Paginate(page, limit, dest, total)
}

func (q *goravelDemoResourceQuery) Get(dest any) error {
	return q.query.Get(dest)
}

func (q *goravelDemoResourceQuery) Insert(data any) error {
	_, err := q.query.Insert(data)
	return err
}

func (q *goravelDemoResourceQuery) Update(data any) error {
	_, err := q.query.Update(data)
	return err
}

func (q *goravelDemoResourceQuery) Delete() error {
	_, err := q.query.Delete()
	return err
}
