package setting

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
)

type ValueType string

const (
	ValueTypeString  ValueType = "string"
	ValueTypeBoolean ValueType = "boolean"
	ValueTypeInteger ValueType = "integer"
	ValueTypeNumber  ValueType = "number"
	ValueTypeJSON    ValueType = "json"
)

var (
	ErrInvalidSetting   = errors.New("invalid setting")
	ErrInvalidValueType = errors.New("invalid setting value type")
	ErrInvalidValue     = errors.New("invalid setting value")
	ErrNotFound         = errors.New("setting not found")
)

type Setting struct {
	Namespace   string    `json:"namespace"`
	Key         string    `json:"key"`
	Value       any       `json:"value"`
	ValueType   ValueType `json:"value_type"`
	Description string    `json:"description,omitempty"`
	UpdatedAt   string    `json:"updated_at,omitempty"`
}

type Repository interface {
	List(context.Context, string) ([]Setting, error)
	Get(context.Context, string, string) (Setting, error)
	Upsert(context.Context, Setting) (Setting, error)
	Delete(context.Context, string, string) error
}

type Service struct{ repository Repository }

func NewService(repository Repository) *Service { return &Service{repository: repository} }

func (s *Service) List(ctx context.Context, namespace string) ([]Setting, error) {
	if s == nil || s.repository == nil {
		return nil, ErrInvalidSetting
	}
	return s.repository.List(ctx, strings.TrimSpace(namespace))
}

func (s *Service) Get(ctx context.Context, namespace, key string) (Setting, error) {
	if s == nil || s.repository == nil {
		return Setting{}, ErrInvalidSetting
	}
	namespace, key = normalizeIdentity(namespace, key)
	if namespace == "" || key == "" {
		return Setting{}, ErrInvalidSetting
	}
	return s.repository.Get(ctx, namespace, key)
}

func (s *Service) Upsert(ctx context.Context, input Setting) (Setting, error) {
	input.Namespace, input.Key = normalizeIdentity(input.Namespace, input.Key)
	if s == nil || s.repository == nil || input.Namespace == "" || input.Key == "" {
		return Setting{}, ErrInvalidSetting
	}
	if !isValueType(input.ValueType) {
		return Setting{}, ErrInvalidValueType
	}
	if err := validateValue(input.Value, input.ValueType); err != nil {
		return Setting{}, err
	}
	input.Description = strings.TrimSpace(input.Description)
	return s.repository.Upsert(ctx, input)
}

func (s *Service) Delete(ctx context.Context, namespace, key string) error {
	if s == nil || s.repository == nil {
		return ErrInvalidSetting
	}
	namespace, key = normalizeIdentity(namespace, key)
	if namespace == "" || key == "" {
		return ErrInvalidSetting
	}
	return s.repository.Delete(ctx, namespace, key)
}

type MemoryRepository struct {
	mu   sync.RWMutex
	rows map[string]Setting
}

func NewMemoryRepository(initial ...Setting) *MemoryRepository {
	repository := &MemoryRepository{rows: make(map[string]Setting, len(initial))}
	for _, item := range initial {
		item.Namespace, item.Key = normalizeIdentity(item.Namespace, item.Key)
		if item.Namespace != "" && item.Key != "" {
			repository.rows[settingID(item.Namespace, item.Key)] = clone(item)
		}
	}
	return repository
}

func (r *MemoryRepository) List(_ context.Context, namespace string) ([]Setting, error) {
	namespace = strings.TrimSpace(namespace)
	r.mu.RLock()
	items := make([]Setting, 0, len(r.rows))
	for _, item := range r.rows {
		if namespace == "" || item.Namespace == namespace {
			items = append(items, clone(item))
		}
	}
	r.mu.RUnlock()
	for i := 0; i < len(items); i++ {
		for j := i + 1; j < len(items); j++ {
			if items[j].Namespace < items[i].Namespace || (items[j].Namespace == items[i].Namespace && items[j].Key < items[i].Key) {
				items[i], items[j] = items[j], items[i]
			}
		}
	}
	return items, nil
}

func (r *MemoryRepository) Get(_ context.Context, namespace, key string) (Setting, error) {
	r.mu.RLock()
	item, ok := r.rows[settingID(namespace, key)]
	r.mu.RUnlock()
	if !ok {
		return Setting{}, ErrNotFound
	}
	return clone(item), nil
}

func (r *MemoryRepository) Upsert(_ context.Context, input Setting) (Setting, error) {
	input.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	r.mu.Lock()
	r.rows[settingID(input.Namespace, input.Key)] = clone(input)
	r.mu.Unlock()
	return clone(input), nil
}

func (r *MemoryRepository) Delete(_ context.Context, namespace, key string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	id := settingID(namespace, key)
	if _, ok := r.rows[id]; !ok {
		return ErrNotFound
	}
	delete(r.rows, id)
	return nil
}

func normalizeIdentity(namespace, key string) (string, string) {
	return strings.TrimSpace(namespace), strings.TrimSpace(key)
}

func settingID(namespace, key string) string { return namespace + "\x00" + key }

func isValueType(valueType ValueType) bool {
	switch valueType {
	case ValueTypeString, ValueTypeBoolean, ValueTypeInteger, ValueTypeNumber, ValueTypeJSON:
		return true
	default:
		return false
	}
}

func validateValue(value any, valueType ValueType) error {
	switch valueType {
	case ValueTypeString:
		if _, ok := value.(string); !ok {
			return fmt.Errorf("%w: expected string", ErrInvalidValue)
		}
	case ValueTypeBoolean:
		if _, ok := value.(bool); !ok {
			return fmt.Errorf("%w: expected boolean", ErrInvalidValue)
		}
	case ValueTypeInteger:
		switch number := value.(type) {
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			return nil
		case float64:
			if number == float64(int64(number)) {
				return nil
			}
		default:
			return fmt.Errorf("%w: expected integer", ErrInvalidValue)
		}
	case ValueTypeNumber:
		switch value.(type) {
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
			return nil
		default:
			return fmt.Errorf("%w: expected number", ErrInvalidValue)
		}
	case ValueTypeJSON:
		if !json.Valid(mustJSON(value)) {
			return fmt.Errorf("%w: expected JSON value", ErrInvalidValue)
		}
	}
	return nil
}

func mustJSON(value any) []byte {
	data, _ := json.Marshal(value)
	return data
}

func clone(item Setting) Setting { return item }
