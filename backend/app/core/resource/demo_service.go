package resource

import (
	"context"
	"errors"
	"strings"
)

var (
	ErrInvalidResource = errors.New("invalid resource")
	ErrResourceNotFound = errors.New("resource not found")
)

type DemoResource struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	Owner     string `json:"owner"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

type DemoResourceUpdate struct {
	Name   *string `json:"name"`
	Status *string `json:"status"`
	Owner  *string `json:"owner"`
}

type DemoResourceListQuery struct {
	Page       int
	PerPage    int
	Search     string
	Filters    map[string]string
	SortField  string
	SortDesc   bool
}

type DemoResourceRepository interface {
	List(context.Context, DemoResourceListQuery) ([]DemoResource, int, error)
	Get(context.Context, string) (DemoResource, error)
	Create(context.Context, DemoResource) (DemoResource, error)
	Update(context.Context, string, DemoResourceUpdate) (DemoResource, error)
	Delete(context.Context, string) error
	BulkDelete(context.Context, []string) error
}

type DemoResourceService struct {
	repository DemoResourceRepository
}

func NewDemoResourceService(repository DemoResourceRepository) *DemoResourceService {
	return &DemoResourceService{repository: repository}
}

func (s *DemoResourceService) List(ctx context.Context, query DemoResourceListQuery) ([]DemoResource, int, error) {
	if query.Page < 1 {
		query.Page = 1
	}
	if query.PerPage < 1 {
		query.PerPage = 20
	}
	return s.repository.List(ctx, query)
}

func (s *DemoResourceService) Get(ctx context.Context, id string) (DemoResource, error) {
	if strings.TrimSpace(id) == "" {
		return DemoResource{}, ErrInvalidResource
	}
	return s.repository.Get(ctx, strings.TrimSpace(id))
}

func (s *DemoResourceService) Create(ctx context.Context, input DemoResource) (DemoResource, error) {
	input.ID = strings.TrimSpace(input.ID)
	input.Name = strings.TrimSpace(input.Name)
	input.Status = strings.TrimSpace(input.Status)
	input.Owner = strings.TrimSpace(input.Owner)
	if input.ID == "" || input.Name == "" || input.Owner == "" || !validStatus(input.Status) {
		return DemoResource{}, ErrInvalidResource
	}
	return s.repository.Create(ctx, input)
}

func (s *DemoResourceService) Update(ctx context.Context, id string, input DemoResourceUpdate) (DemoResource, error) {
	if strings.TrimSpace(id) == "" {
		return DemoResource{}, ErrInvalidResource
	}
	if input.Name != nil {
		value := strings.TrimSpace(*input.Name)
		input.Name = &value
		if value == "" {
			return DemoResource{}, ErrInvalidResource
		}
	}
	if input.Status != nil {
		value := strings.TrimSpace(*input.Status)
		input.Status = &value
		if !validStatus(value) {
			return DemoResource{}, ErrInvalidResource
		}
	}
	if input.Owner != nil {
		value := strings.TrimSpace(*input.Owner)
		input.Owner = &value
		if value == "" {
			return DemoResource{}, ErrInvalidResource
		}
	}
	if input.Name == nil && input.Status == nil && input.Owner == nil {
		return DemoResource{}, ErrInvalidResource
	}
	return s.repository.Update(ctx, strings.TrimSpace(id), input)
}

func (s *DemoResourceService) Delete(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return ErrInvalidResource
	}
	return s.repository.Delete(ctx, strings.TrimSpace(id))
}

func (s *DemoResourceService) BulkDelete(ctx context.Context, ids []string) error {
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
	return s.repository.BulkDelete(ctx, cleaned)
}

func validStatus(status string) bool {
	return status == "draft" || status == "active"
}
