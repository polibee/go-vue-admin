package resource

import (
	"context"
	"strings"
)

type CoreResourceService struct {
	repository CoreResourceRepository
}

func NewCoreResourceService(repository CoreResourceRepository) *CoreResourceService {
	return &CoreResourceService{repository: repository}
}

func (s *CoreResourceService) List(ctx context.Context, query DemoResourceListQuery) ([]CoreResourceRecord, int, error) {
	if s == nil || s.repository == nil {
		return nil, 0, ErrInvalidResource
	}
	if query.Page < 1 {
		query.Page = 1
	}
	if query.PerPage < 1 {
		query.PerPage = 20
	}
	return s.repository.List(ctx, query)
}

func (s *CoreResourceService) Get(ctx context.Context, id string) (CoreResourceRecord, error) {
	if s == nil || s.repository == nil || strings.TrimSpace(id) == "" {
		return nil, ErrInvalidResource
	}
	return s.repository.Get(ctx, strings.TrimSpace(id))
}

func (s *CoreResourceService) Create(ctx context.Context, input CoreResourceRecord) (CoreResourceRecord, error) {
	if s == nil || s.repository == nil || coreResourceID(input) == "" {
		return nil, ErrInvalidResource
	}
	return s.repository.Create(ctx, input)
}

func (s *CoreResourceService) Update(ctx context.Context, id string, input CoreResourceRecord) (CoreResourceRecord, error) {
	if s == nil || s.repository == nil || strings.TrimSpace(id) == "" || len(input) == 0 {
		return nil, ErrInvalidResource
	}
	return s.repository.Update(ctx, strings.TrimSpace(id), input)
}

func (s *CoreResourceService) Delete(ctx context.Context, id string) error {
	if s == nil || s.repository == nil || strings.TrimSpace(id) == "" {
		return ErrInvalidResource
	}
	return s.repository.Delete(ctx, strings.TrimSpace(id))
}

func (s *CoreResourceService) BulkDelete(ctx context.Context, ids []string) error {
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
