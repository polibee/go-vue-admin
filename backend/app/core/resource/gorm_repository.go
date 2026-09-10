package resource

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"gorm.io/gorm"
)

var resourceSQLIdentifierPattern = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

type ResourceSchema struct {
	Table      string
	PrimaryKey string
	Fields     []string
	Searchable []string
	Sortable   []string
	SoftDelete bool
}

type GormResourceOptions struct {
	DB     *gorm.DB
	Schema ResourceSchema
}

type GormResourceRepository struct {
	db     *gorm.DB
	schema ResourceSchema
}

func NewGormResourceRepository(options GormResourceOptions) (*GormResourceRepository, error) {
	if options.DB == nil {
		return nil, errors.New("resource gorm database cannot be nil")
	}
	schema, err := normalizeResourceSchema(options.Schema)
	if err != nil {
		return nil, err
	}
	return &GormResourceRepository{db: options.DB, schema: schema}, nil
}

func normalizeResourceSchema(schema ResourceSchema) (ResourceSchema, error) {
	for label, value := range map[string]string{"table": schema.Table, "primary key": schema.PrimaryKey} {
		if !resourceSQLIdentifierPattern.MatchString(value) {
			return ResourceSchema{}, fmt.Errorf("unsafe resource %s %q", label, value)
		}
	}
	if len(schema.Fields) == 0 {
		return ResourceSchema{}, errors.New("resource fields cannot be empty")
	}
	for _, field := range append(append([]string{}, schema.Fields...), append(schema.Searchable, schema.Sortable...)...) {
		if !resourceSQLIdentifierPattern.MatchString(field) {
			return ResourceSchema{}, fmt.Errorf("unsafe resource field %q", field)
		}
	}
	schema.Fields = uniqueStrings(schema.Fields)
	schema.Searchable = uniqueStrings(schema.Searchable)
	schema.Sortable = uniqueStrings(schema.Sortable)
	if len(schema.Sortable) == 0 {
		schema.Sortable = append([]string(nil), schema.Fields...)
	}
	return schema, nil
}

func (r *GormResourceRepository) List(ctx context.Context, query DemoResourceListQuery) ([]CoreResourceRecord, int, error) {
	databaseQuery := r.activeQuery(ctx)
	if search := strings.TrimSpace(query.Search); search != "" && len(r.schema.Searchable) > 0 {
		pattern := "%" + search + "%"
		conditions := make([]string, 0, len(r.schema.Searchable))
		arguments := make([]any, 0, len(r.schema.Searchable))
		for _, field := range r.schema.Searchable {
			conditions = append(conditions, "LOWER("+field+") LIKE LOWER(?)")
			arguments = append(arguments, pattern)
		}
		databaseQuery = databaseQuery.Where("("+strings.Join(conditions, " OR ")+")", arguments...)
	}
	for field, value := range query.Filters {
		if r.contains(r.schema.Fields, field) && strings.TrimSpace(value) != "" {
			databaseQuery = databaseQuery.Where(field+" = ?", strings.TrimSpace(value))
		}
	}
	var total int64
	if err := databaseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	page, perPage := query.Page, query.PerPage
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}
	sortField := r.schema.PrimaryKey
	if r.contains(r.schema.Sortable, query.SortField) {
		sortField = query.SortField
	}
	direction := " ASC"
	if query.SortDesc {
		direction = " DESC"
	}
	var rows []map[string]any
	err := databaseQuery.Order(sortField + direction).Offset((page - 1) * perPage).Limit(perPage).Find(&rows).Error
	return coreResourceRecords(rows), int(total), err
}

func (r *GormResourceRepository) Get(ctx context.Context, id string) (CoreResourceRecord, error) {
	var row map[string]any
	err := r.activeQuery(ctx).Where(r.schema.PrimaryKey+" = ?", strings.TrimSpace(id)).Take(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrResourceNotFound
	}
	if err != nil {
		return nil, err
	}
	return CoreResourceRecord(row), nil
}

func (r *GormResourceRepository) Create(ctx context.Context, input CoreResourceRecord) (CoreResourceRecord, error) {
	values := r.allowedValues(input, true)
	if strings.TrimSpace(coreResourceValue(input, r.schema.PrimaryKey)) == "" {
		return nil, ErrInvalidResource
	}
	if err := r.db.WithContext(ctx).Table(r.schema.Table).Create(map[string]any(values)).Error; err != nil {
		return nil, err
	}
	return r.Get(ctx, coreResourceValue(input, r.schema.PrimaryKey))
}

func (r *GormResourceRepository) Update(ctx context.Context, id string, input CoreResourceRecord) (CoreResourceRecord, error) {
	values := r.allowedValues(input, false)
	if len(values) == 0 {
		return nil, ErrInvalidResource
	}
	if err := r.activeQuery(ctx).Where(r.schema.PrimaryKey+" = ?", id).Updates(map[string]any(values)).Error; err != nil {
		return nil, err
	}
	return r.Get(ctx, id)
}

func (r *GormResourceRepository) Delete(ctx context.Context, id string) error {
	if _, err := r.Get(ctx, id); err != nil {
		return err
	}
	query := r.activeQuery(ctx).Where(r.schema.PrimaryKey+" = ?", id)
	if r.schema.SoftDelete {
		return query.Updates(map[string]any{"deleted_at": time.Now().UTC()}).Error
	}
	return query.Delete(nil).Error
}

func (r *GormResourceRepository) BulkDelete(ctx context.Context, ids []string) error {
	cleaned := make([]string, 0, len(ids))
	seen := make(map[string]struct{}, len(ids))
	for _, id := range ids {
		id = strings.TrimSpace(id)
		if id == "" {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		cleaned = append(cleaned, id)
	}
	if len(cleaned) == 0 {
		return ErrInvalidResource
	}
	query := r.activeQuery(ctx).Where(r.schema.PrimaryKey+" IN ?", cleaned)
	if r.schema.SoftDelete {
		return query.Updates(map[string]any{"deleted_at": time.Now().UTC()}).Error
	}
	return query.Delete(nil).Error
}

func (r *GormResourceRepository) activeQuery(ctx context.Context) *gorm.DB {
	query := r.db.WithContext(ctx).Table(r.schema.Table)
	if r.schema.SoftDelete {
		query = query.Where("deleted_at IS NULL")
	}
	return query
}

func (r *GormResourceRepository) allowedValues(input CoreResourceRecord, includePrimaryKey bool) CoreResourceRecord {
	allowed := make(map[string]struct{}, len(r.schema.Fields))
	for _, field := range r.schema.Fields {
		allowed[field] = struct{}{}
	}
	values := make(CoreResourceRecord)
	for field, value := range input {
		if field == r.schema.PrimaryKey && !includePrimaryKey {
			continue
		}
		if _, ok := allowed[field]; ok {
			values[field] = value
		}
	}
	return values
}

func (r *GormResourceRepository) contains(values []string, candidate string) bool {
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
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	return result
}

func coreResourceRecords(rows []map[string]any) []CoreResourceRecord {
	result := make([]CoreResourceRecord, 0, len(rows))
	for _, row := range rows {
		result = append(result, CoreResourceRecord(row))
	}
	return result
}
