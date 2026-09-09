package resource

import (
	"context"
	"errors"
	"sort"
	"strings"
	"sync"
)

var ErrResourceExists = errors.New("resource already exists")

type MemoryDemoResourceRepository struct {
	mu   sync.RWMutex
	rows map[string]DemoResource
}

func NewMemoryDemoResourceRepository(initial ...DemoResource) DemoResourceRepository {
	rows := make(map[string]DemoResource, len(initial))
	for _, item := range initial {
		if item.ID != "" {
			rows[item.ID] = cloneDemoResource(item)
		}
	}
	return &MemoryDemoResourceRepository{rows: rows}
}

func (r *MemoryDemoResourceRepository) List(_ context.Context, query DemoResourceListQuery) ([]DemoResource, int, error) {
	r.mu.RLock()
	rows := make([]DemoResource, 0, len(r.rows))
	for _, item := range r.rows {
		rows = append(rows, cloneDemoResource(item))
	}
	r.mu.RUnlock()

	rows = filterDemoResources(rows, query)
	sortDemoResources(rows, query.SortField, query.SortDesc)
	total := len(rows)
	page := query.Page
	if page < 1 {
		page = 1
	}
	perPage := query.PerPage
	if perPage < 1 {
		perPage = 20
	}
	start := (page - 1) * perPage
	if start >= len(rows) {
		return []DemoResource{}, total, nil
	}
	end := start + perPage
	if end > len(rows) {
		end = len(rows)
	}
	return rows[start:end], total, nil
}

func (r *MemoryDemoResourceRepository) Get(_ context.Context, id string) (DemoResource, error) {
	r.mu.RLock()
	item, ok := r.rows[id]
	r.mu.RUnlock()
	if !ok {
		return DemoResource{}, ErrResourceNotFound
	}
	return cloneDemoResource(item), nil
}

func (r *MemoryDemoResourceRepository) Create(_ context.Context, input DemoResource) (DemoResource, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.rows[input.ID]; ok {
		return DemoResource{}, ErrResourceExists
	}
	r.rows[input.ID] = cloneDemoResource(input)
	return cloneDemoResource(input), nil
}

func (r *MemoryDemoResourceRepository) Update(_ context.Context, id string, input DemoResourceUpdate) (DemoResource, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	item, ok := r.rows[id]
	if !ok {
		return DemoResource{}, ErrResourceNotFound
	}
	if input.Name != nil {
		item.Name = *input.Name
	}
	if input.Status != nil {
		item.Status = *input.Status
	}
	if input.Owner != nil {
		item.Owner = *input.Owner
	}
	r.rows[id] = item
	return cloneDemoResource(item), nil
}

func (r *MemoryDemoResourceRepository) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.rows[id]; !ok {
		return ErrResourceNotFound
	}
	delete(r.rows, id)
	return nil
}

func (r *MemoryDemoResourceRepository) BulkDelete(_ context.Context, ids []string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, id := range ids {
		delete(r.rows, id)
	}
	return nil
}

func filterDemoResources(rows []DemoResource, query DemoResourceListQuery) []DemoResource {
	search := strings.ToLower(strings.TrimSpace(query.Search))
	filtered := rows[:0]
	for _, item := range rows {
		if search != "" && !strings.Contains(strings.ToLower(strings.Join([]string{item.ID, item.Name, item.Status, item.Owner}, " ")), search) {
			continue
		}
		matches := true
		for field, expected := range query.Filters {
			if resourceField(item, field) != expected {
				matches = false
				break
			}
		}
		if matches {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

func sortDemoResources(rows []DemoResource, field string, desc bool) {
	if !isSortableResourceField(field) {
		field = "name"
	}
	sort.SliceStable(rows, func(i, j int) bool {
		left, right := resourceField(rows[i], field), resourceField(rows[j], field)
		if desc {
			return left > right
		}
		return left < right
	})
}

func isSortableResourceField(field string) bool {
	switch field {
	case "id", "name", "status", "owner":
		return true
	default:
		return false
	}
}

func resourceField(item DemoResource, field string) string {
	switch field {
	case "id":
		return item.ID
	case "status":
		return item.Status
	case "owner":
		return item.Owner
	default:
		return item.Name
	}
}

func cloneDemoResource(item DemoResource) DemoResource {
	return item
}
