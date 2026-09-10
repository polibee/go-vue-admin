package audit

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidEntry = errors.New("invalid audit entry")
	ErrNotFound     = errors.New("audit entry not found")
)

type Entry struct {
	ID           string `json:"id"`
	ActorID      string `json:"actor_id"`
	ActorEmail   string `json:"actor_email"`
	Action       string `json:"action"`
	ResourceType string `json:"resource_type"`
	ResourceID   string `json:"resource_id"`
	Before       any    `json:"before"`
	After        any    `json:"after"`
	IP           string `json:"ip"`
	UserAgent    string `json:"user_agent"`
	CreatedAt    string `json:"created_at"`
}

type Repository interface {
	List(context.Context) ([]Entry, error)
	Get(context.Context, string) (Entry, error)
	Create(context.Context, Entry) (Entry, error)
}

type Service struct{ repository Repository }

func NewService(repository Repository) *Service { return &Service{repository: repository} }

func (s *Service) Record(ctx context.Context, entry Entry) (Entry, error) {
	if s == nil || s.repository == nil || strings.TrimSpace(entry.Action) == "" || strings.TrimSpace(entry.ResourceType) == "" || strings.TrimSpace(entry.ResourceID) == "" {
		return Entry{}, ErrInvalidEntry
	}
	entry.Action = strings.TrimSpace(entry.Action)
	entry.ResourceType = strings.TrimSpace(entry.ResourceType)
	entry.ResourceID = strings.TrimSpace(entry.ResourceID)
	if entry.ID == "" {
		entry.ID = uuid.NewString()
	}
	if entry.CreatedAt == "" {
		entry.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	}
	return s.repository.Create(ctx, entry)
}

func (s *Service) List(ctx context.Context) ([]Entry, error) {
	if s == nil || s.repository == nil {
		return nil, ErrInvalidEntry
	}
	return s.repository.List(ctx)
}

func (s *Service) Get(ctx context.Context, id string) (Entry, error) {
	if s == nil || s.repository == nil || strings.TrimSpace(id) == "" {
		return Entry{}, ErrInvalidEntry
	}
	return s.repository.Get(ctx, strings.TrimSpace(id))
}

type MemoryRepository struct {
	mu   sync.RWMutex
	rows map[string]Entry
}

func NewMemoryRepository(initial ...Entry) *MemoryRepository {
	rows := make(map[string]Entry, len(initial))
	for _, entry := range initial {
		if strings.TrimSpace(entry.ID) != "" {
			rows[entry.ID] = entry
		}
	}
	return &MemoryRepository{rows: rows}
}

func (r *MemoryRepository) List(_ context.Context) ([]Entry, error) {
	r.mu.RLock()
	items := make([]Entry, 0, len(r.rows))
	for _, entry := range r.rows {
		items = append(items, entry)
	}
	r.mu.RUnlock()
	for i := 0; i < len(items); i++ {
		for j := i + 1; j < len(items); j++ {
			if items[j].CreatedAt > items[i].CreatedAt {
				items[i], items[j] = items[j], items[i]
			}
		}
	}
	return items, nil
}

func (r *MemoryRepository) Get(_ context.Context, id string) (Entry, error) {
	r.mu.RLock()
	entry, ok := r.rows[id]
	r.mu.RUnlock()
	if !ok {
		return Entry{}, ErrNotFound
	}
	return entry, nil
}

func (r *MemoryRepository) Create(_ context.Context, entry Entry) (Entry, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.rows[entry.ID]; exists {
		return Entry{}, ErrInvalidEntry
	}
	r.rows[entry.ID] = entry
	return entry, nil
}
