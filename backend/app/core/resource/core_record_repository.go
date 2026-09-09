package resource

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
)

type CoreResourceRecord map[string]any

type CoreResourceRepository interface {
	List(context.Context, DemoResourceListQuery) ([]CoreResourceRecord, int, error)
	Get(context.Context, string) (CoreResourceRecord, error)
	Create(context.Context, CoreResourceRecord) (CoreResourceRecord, error)
	Update(context.Context, string, CoreResourceRecord) (CoreResourceRecord, error)
	Delete(context.Context, string) error
	BulkDelete(context.Context, []string) error
}

type MemoryCoreResourceRepository struct {
	mu   sync.RWMutex
	rows map[string]CoreResourceRecord
}

func NewMemoryCoreResourceRepository(initial ...CoreResourceRecord) CoreResourceRepository {
	rows := make(map[string]CoreResourceRecord, len(initial))
	for _, item := range initial {
		if id := coreResourceID(item); id != "" {
			rows[id] = cloneCoreResourceRecord(item)
		}
	}
	return &MemoryCoreResourceRepository{rows: rows}
}

func (r *MemoryCoreResourceRepository) List(_ context.Context, query DemoResourceListQuery) ([]CoreResourceRecord, int, error) {
	r.mu.RLock()
	rows := make([]CoreResourceRecord, 0, len(r.rows))
	for _, item := range r.rows {
		rows = append(rows, cloneCoreResourceRecord(item))
	}
	r.mu.RUnlock()

	search := strings.ToLower(strings.TrimSpace(query.Search))
	filtered := rows[:0]
	for _, item := range rows {
		if search != "" && !strings.Contains(strings.ToLower(coreResourceSearchText(item)), search) {
			continue
		}
		matches := true
		for field, expected := range query.Filters {
			if coreResourceValue(item, field) != expected {
				matches = false
				break
			}
		}
		if matches {
			filtered = append(filtered, item)
		}
	}

	field := query.SortField
	if field == "" {
		field = "id"
	}
	sort.SliceStable(filtered, func(i, j int) bool {
		left, right := coreResourceValue(filtered[i], field), coreResourceValue(filtered[j], field)
		if query.SortDesc {
			return left > right
		}
		return left < right
	})

	total := len(filtered)
	page, perPage := query.Page, query.PerPage
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}
	start := (page - 1) * perPage
	if start >= len(filtered) {
		return []CoreResourceRecord{}, total, nil
	}
	end := start + perPage
	if end > len(filtered) {
		end = len(filtered)
	}
	return filtered[start:end], total, nil
}

func (r *MemoryCoreResourceRepository) Get(_ context.Context, id string) (CoreResourceRecord, error) {
	r.mu.RLock()
	item, ok := r.rows[id]
	r.mu.RUnlock()
	if !ok {
		return nil, ErrResourceNotFound
	}
	return cloneCoreResourceRecord(item), nil
}

func (r *MemoryCoreResourceRepository) Create(_ context.Context, input CoreResourceRecord) (CoreResourceRecord, error) {
	id := coreResourceID(input)
	if id == "" {
		return nil, ErrInvalidResource
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.rows[id]; ok {
		return nil, ErrResourceExists
	}
	r.rows[id] = cloneCoreResourceRecord(input)
	return cloneCoreResourceRecord(input), nil
}

func (r *MemoryCoreResourceRepository) Update(_ context.Context, id string, input CoreResourceRecord) (CoreResourceRecord, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	item, ok := r.rows[id]
	if !ok {
		return nil, ErrResourceNotFound
	}
	for key, value := range input {
		if key != "id" {
			item[key] = value
		}
	}
	r.rows[id] = item
	return cloneCoreResourceRecord(item), nil
}

func (r *MemoryCoreResourceRepository) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.rows[id]; !ok {
		return ErrResourceNotFound
	}
	delete(r.rows, id)
	return nil
}

func (r *MemoryCoreResourceRepository) BulkDelete(_ context.Context, ids []string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, id := range ids {
		delete(r.rows, strings.TrimSpace(id))
	}
	return nil
}

func coreResourceID(item CoreResourceRecord) string {
	return strings.TrimSpace(coreResourceValue(item, "id"))
}

func coreResourceSearchText(item CoreResourceRecord) string {
	values := make([]string, 0, len(item))
	for _, value := range item {
		values = append(values, stringifyCoreResourceValue(value))
	}
	return strings.Join(values, " ")
}

func coreResourceValue(item CoreResourceRecord, field string) string {
	return stringifyCoreResourceValue(item[field])
}

func stringifyCoreResourceValue(value any) string {
	switch value := value.(type) {
	case nil:
		return ""
	case []string:
		return strings.Join(value, ",")
	default:
		return strings.TrimSpace(fmt.Sprint(value))
	}
}

func cloneCoreResourceRecord(item CoreResourceRecord) CoreResourceRecord {
	clone := make(CoreResourceRecord, len(item))
	for key, value := range item {
		switch value := value.(type) {
		case []string:
			clone[key] = append([]string(nil), value...)
		default:
			clone[key] = value
		}
	}
	return clone
}
