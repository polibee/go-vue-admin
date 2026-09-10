package resource

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/goravel/framework/contracts/database/db"
	"github.com/goravel/framework/facades"
)

var coreResourceSQLIdentifierPattern = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

type MySQLCoreResourceOptions struct {
	Table      string
	PrimaryKey string
	Fields     []string
	Searchable []string
	Sortable   []string
	SoftDelete bool
}

type coreResourceQuery interface {
	Where(any, ...any) coreResourceQuery
	WhereAny([]string, ...any) coreResourceQuery
	WhereNull(string) coreResourceQuery
	WhereIn(string, []any) coreResourceQuery
	OrderBy(string, ...string) coreResourceQuery
	OrderByDesc(string) coreResourceQuery
	Paginate(int, int, any, *int64) error
	Get(any) error
	Insert(any) error
	Update(any) error
	Delete() error
}

type coreResourceQueryFactory func(context.Context) coreResourceQuery

type MySQLCoreResourceRepository struct {
	options MySQLCoreResourceOptions
	query   coreResourceQueryFactory
}

func NewMySQLCoreResourceRepository(options MySQLCoreResourceOptions) (CoreResourceRepository, error) {
	return newMySQLCoreResourceRepositoryWithQuery(options, func(ctx context.Context) coreResourceQuery {
		return &goravelCoreResourceQuery{query: facades.DB().WithContext(ctx).Table(options.Table)}
	})
}

func newMySQLCoreResourceRepositoryWithQuery(options MySQLCoreResourceOptions, factory coreResourceQueryFactory) (*MySQLCoreResourceRepository, error) {
	if err := validateMySQLCoreResourceOptions(options); err != nil {
		return nil, err
	}
	return &MySQLCoreResourceRepository{options: normalizeMySQLCoreResourceOptions(options), query: factory}, nil
}

func validateMySQLCoreResourceOptions(options MySQLCoreResourceOptions) error {
	for label, value := range map[string]string{"table": options.Table, "primary key": options.PrimaryKey} {
		if !coreResourceSQLIdentifierPattern.MatchString(value) {
			return fmt.Errorf("unsafe resource %s %q", label, value)
		}
	}
	if len(options.Fields) == 0 {
		return errors.New("resource fields cannot be empty")
	}
	for _, field := range append(append([]string{}, options.Fields...), append(options.Searchable, options.Sortable...)...) {
		if !coreResourceSQLIdentifierPattern.MatchString(field) {
			return fmt.Errorf("unsafe resource field %q", field)
		}
	}
	return nil
}

func normalizeMySQLCoreResourceOptions(options MySQLCoreResourceOptions) MySQLCoreResourceOptions {
	options.Fields = uniqueStrings(options.Fields)
	options.Searchable = uniqueStrings(options.Searchable)
	options.Sortable = uniqueStrings(options.Sortable)
	if len(options.Sortable) == 0 {
		options.Sortable = append([]string(nil), options.Fields...)
	}
	return options
}

func (r *MySQLCoreResourceRepository) List(ctx context.Context, query DemoResourceListQuery) ([]CoreResourceRecord, int, error) {
	databaseQuery := r.activeQuery(ctx)
	if query.Search != "" && len(r.options.Searchable) > 0 {
		databaseQuery = databaseQuery.WhereAny(r.options.Searchable, "like", "%"+strings.TrimSpace(query.Search)+"%")
	}
	for field, value := range query.Filters {
		if r.contains(r.options.Fields, field) && strings.TrimSpace(value) != "" {
			databaseQuery = databaseQuery.Where(field, strings.TrimSpace(value))
		}
	}
	sortField := r.options.PrimaryKey
	if r.contains(r.options.Sortable, query.SortField) {
		sortField = query.SortField
	}
	if query.SortDesc {
		databaseQuery = databaseQuery.OrderByDesc(sortField)
	} else {
		databaseQuery = databaseQuery.OrderBy(sortField)
	}
	var rows []CoreResourceRecord
	var total int64
	if err := databaseQuery.Paginate(query.Page, query.PerPage, &rows, &total); err != nil {
		return nil, 0, err
	}
	return rows, int(total), nil
}

func (r *MySQLCoreResourceRepository) Get(ctx context.Context, id string) (CoreResourceRecord, error) {
	var rows []CoreResourceRecord
	if err := r.activeQuery(ctx).Where(r.options.PrimaryKey, id).Get(&rows); err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, ErrResourceNotFound
	}
	return rows[0], nil
}

func (r *MySQLCoreResourceRepository) Create(ctx context.Context, input CoreResourceRecord) (CoreResourceRecord, error) {
	values := r.allowedValues(input, true)
	if strings.TrimSpace(coreResourceValue(input, r.options.PrimaryKey)) == "" {
		return nil, ErrInvalidResource
	}
	if err := r.query(ctx).Insert(values); err != nil {
		return nil, err
	}
	return r.Get(ctx, coreResourceValue(input, r.options.PrimaryKey))
}

func (r *MySQLCoreResourceRepository) Update(ctx context.Context, id string, input CoreResourceRecord) (CoreResourceRecord, error) {
	values := r.allowedValues(input, false)
	if len(values) == 0 {
		return nil, ErrInvalidResource
	}
	if err := r.query(ctx).Where(r.options.PrimaryKey, id).Update(values); err != nil {
		return nil, err
	}
	return r.Get(ctx, id)
}

func (r *MySQLCoreResourceRepository) Delete(ctx context.Context, id string) error {
	if _, err := r.Get(ctx, id); err != nil {
		return err
	}
	query := r.query(ctx).Where(r.options.PrimaryKey, id)
	if r.options.SoftDelete {
		return query.Update(map[string]any{"deleted_at": time.Now().UTC()})
	}
	return query.Delete()
}

func (r *MySQLCoreResourceRepository) BulkDelete(ctx context.Context, ids []string) error {
	values := make([]any, 0, len(ids))
	for _, id := range ids {
		if strings.TrimSpace(id) != "" {
			values = append(values, strings.TrimSpace(id))
		}
	}
	if len(values) == 0 {
		return ErrInvalidResource
	}
	query := r.activeQuery(ctx).WhereIn(r.options.PrimaryKey, values)
	if r.options.SoftDelete {
		return query.Update(map[string]any{"deleted_at": time.Now().UTC()})
	}
	return query.Delete()
}

func (r *MySQLCoreResourceRepository) activeQuery(ctx context.Context) coreResourceQuery {
	query := r.query(ctx)
	if r.options.SoftDelete {
		return query.WhereNull("deleted_at")
	}
	return query
}

func (r *MySQLCoreResourceRepository) allowedValues(input CoreResourceRecord, includePrimaryKey bool) CoreResourceRecord {
	allowed := make(map[string]struct{}, len(r.options.Fields))
	for _, field := range r.options.Fields {
		allowed[field] = struct{}{}
	}
	values := make(CoreResourceRecord)
	for field, value := range input {
		if field == r.options.PrimaryKey && !includePrimaryKey {
			continue
		}
		if _, ok := allowed[field]; ok {
			values[field] = value
		}
	}
	return values
}

func (r *MySQLCoreResourceRepository) contains(values []string, candidate string) bool {
	for _, value := range values {
		if value == candidate {
			return true
		}
	}
	return false
}

func uniqueStrings(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))
	for _, value := range values {
		if _, ok := seen[value]; ok || value == "" {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	return result
}

type goravelCoreResourceQuery struct{ query db.Query }

func (q *goravelCoreResourceQuery) Where(value any, args ...any) coreResourceQuery {
	return &goravelCoreResourceQuery{query: q.query.Where(value, args...)}
}
func (q *goravelCoreResourceQuery) WhereAny(columns []string, args ...any) coreResourceQuery {
	return &goravelCoreResourceQuery{query: q.query.WhereAny(columns, args...)}
}
func (q *goravelCoreResourceQuery) WhereNull(column string) coreResourceQuery {
	return &goravelCoreResourceQuery{query: q.query.WhereNull(column)}
}
func (q *goravelCoreResourceQuery) WhereIn(column string, values []any) coreResourceQuery {
	return &goravelCoreResourceQuery{query: q.query.WhereIn(column, values)}
}
func (q *goravelCoreResourceQuery) OrderBy(column string, directions ...string) coreResourceQuery {
	return &goravelCoreResourceQuery{query: q.query.OrderBy(column, directions...)}
}
func (q *goravelCoreResourceQuery) OrderByDesc(column string) coreResourceQuery {
	return &goravelCoreResourceQuery{query: q.query.OrderByDesc(column)}
}
func (q *goravelCoreResourceQuery) Paginate(page, limit int, dest any, total *int64) error {
	return q.query.Paginate(page, limit, dest, total)
}
func (q *goravelCoreResourceQuery) Get(dest any) error { return q.query.Get(dest) }
func (q *goravelCoreResourceQuery) Insert(data any) error {
	_, err := q.query.Insert(data)
	return err
}
func (q *goravelCoreResourceQuery) Update(data any) error {
	_, err := q.query.Update(data)
	return err
}
func (q *goravelCoreResourceQuery) Delete() error {
	_, err := q.query.Delete()
	return err
}
