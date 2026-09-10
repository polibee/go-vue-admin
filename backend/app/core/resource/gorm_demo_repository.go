package resource

import (
	"context"
	"fmt"
	"time"
)

type GormDemoResourceRepository struct {
	repository ResourceRepository
}

func NewGormDemoResourceRepository(repository ResourceRepository) DemoResourceRepository {
	return &GormDemoResourceRepository{repository: repository}
}

func (r *GormDemoResourceRepository) List(ctx context.Context, query DemoResourceListQuery) ([]DemoResource, int, error) {
	rows, total, err := r.repository.List(ctx, query)
	if err != nil {
		return nil, 0, err
	}
	result := make([]DemoResource, 0, len(rows))
	for _, row := range rows {
		result = append(result, demoResourceFromRecord(row))
	}
	return result, total, nil
}

func (r *GormDemoResourceRepository) Get(ctx context.Context, id string) (DemoResource, error) {
	row, err := r.repository.Get(ctx, id)
	if err != nil {
		return DemoResource{}, err
	}
	return demoResourceFromRecord(row), nil
}

func (r *GormDemoResourceRepository) Create(ctx context.Context, input DemoResource) (DemoResource, error) {
	row, err := r.repository.Create(ctx, demoResourceRecord(input))
	if err != nil {
		return DemoResource{}, err
	}
	return demoResourceFromRecord(row), nil
}

func (r *GormDemoResourceRepository) Update(ctx context.Context, id string, input DemoResourceUpdate) (DemoResource, error) {
	values := CoreResourceRecord{}
	if input.Name != nil {
		values["name"] = *input.Name
	}
	if input.Status != nil {
		values["status"] = *input.Status
	}
	if input.Owner != nil {
		values["owner"] = *input.Owner
	}
	row, err := r.repository.Update(ctx, id, values)
	if err != nil {
		return DemoResource{}, err
	}
	return demoResourceFromRecord(row), nil
}

func (r *GormDemoResourceRepository) Delete(ctx context.Context, id string) error {
	return r.repository.Delete(ctx, id)
}

func (r *GormDemoResourceRepository) BulkDelete(ctx context.Context, ids []string) error {
	return r.repository.BulkDelete(ctx, ids)
}

func demoResourceRecord(value DemoResource) CoreResourceRecord {
	return CoreResourceRecord{"id": value.ID, "name": value.Name, "status": value.Status, "owner": value.Owner}
}

func demoResourceFromRecord(value CoreResourceRecord) DemoResource {
	return DemoResource{
		ID: stringResourceValue(value["id"]), Name: stringResourceValue(value["name"]),
		Status: stringResourceValue(value["status"]), Owner: stringResourceValue(value["owner"]),
		CreatedAt: resourceTimeValue(value["created_at"]), UpdatedAt: resourceTimeValue(value["updated_at"]),
	}
}

func stringResourceValue(value any) string {
	if value == nil {
		return ""
	}
	return fmt.Sprint(value)
}

func resourceTimeValue(value any) string {
	if timestamp, ok := value.(time.Time); ok && !timestamp.IsZero() {
		return timestamp.UTC().Format(time.RFC3339)
	}
	return stringResourceValue(value)
}
