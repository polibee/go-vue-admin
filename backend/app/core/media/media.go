package media

import (
	"context"
	"errors"
	"io"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

const MaxUploadSize int64 = 10 * 1024 * 1024

var (
	ErrInvalidMedia  = errors.New("invalid media")
	ErrMediaTooLarge = errors.New("media file is too large")
	ErrNotFound      = errors.New("media not found")
)

type Media struct {
	ID           string `json:"id"`
	Disk         string `json:"disk"`
	Path         string `json:"path"`
	OriginalName string `json:"original_name"`
	MIMEType     string `json:"mime_type"`
	Size         int64  `json:"size"`
	URL          string `json:"url"`
	CreatedAt    string `json:"created_at,omitempty"`
}

type UploadInput struct {
	OriginalName string
	MIMEType     string
	Size         int64
	Content      io.Reader
}

type Storage interface {
	Put(context.Context, string, io.Reader) error
	Delete(context.Context, string) error
	Path(string) string
	URL(string) string
}

type Repository interface {
	List(context.Context) ([]Media, error)
	Get(context.Context, string) (Media, error)
	Create(context.Context, Media) (Media, error)
	Delete(context.Context, string) error
}

type Service struct {
	repository Repository
	storage    Storage
}

func NewService(repository Repository, storage Storage) *Service {
	return &Service{repository: repository, storage: storage}
}

func (s *Service) List(ctx context.Context) ([]Media, error) {
	if s == nil || s.repository == nil || s.storage == nil {
		return nil, ErrInvalidMedia
	}
	items, err := s.repository.List(ctx)
	if err != nil {
		return nil, err
	}
	for i := range items {
		items[i].URL = s.storage.URL(items[i].Path)
	}
	return items, nil
}

func (s *Service) Get(ctx context.Context, id string) (Media, error) {
	if s == nil || s.repository == nil || s.storage == nil || strings.TrimSpace(id) == "" {
		return Media{}, ErrInvalidMedia
	}
	item, err := s.repository.Get(ctx, strings.TrimSpace(id))
	if err != nil {
		return Media{}, err
	}
	item.URL = s.storage.URL(item.Path)
	return item, nil
}

func (s *Service) Upload(ctx context.Context, input UploadInput) (Media, error) {
	if s == nil || s.repository == nil || s.storage == nil || input.Content == nil {
		return Media{}, ErrInvalidMedia
	}
	if input.Size < 0 || input.Size > MaxUploadSize {
		return Media{}, ErrMediaTooLarge
	}
	name := strings.TrimSpace(input.OriginalName)
	if name == "" || filepath.Base(name) != name || strings.Contains(name, "..") || strings.TrimSpace(input.MIMEType) == "" {
		return Media{}, ErrInvalidMedia
	}
	extension := strings.ToLower(filepath.Ext(name))
	if len(extension) > 16 {
		return Media{}, ErrInvalidMedia
	}
	id := uuid.NewString()
	path := id + extension
	if err := s.storage.Put(ctx, path, input.Content); err != nil {
		return Media{}, err
	}
	item := Media{
		ID: id, Disk: "local", Path: path, OriginalName: name,
		MIMEType: strings.TrimSpace(input.MIMEType), Size: input.Size,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
	created, err := s.repository.Create(ctx, item)
	if err != nil {
		_ = s.storage.Delete(ctx, path)
		return Media{}, err
	}
	created.URL = s.storage.URL(created.Path)
	return created, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	item, err := s.Get(ctx, id)
	if err != nil {
		return err
	}
	if err := s.storage.Delete(ctx, item.Path); err != nil {
		return err
	}
	return s.repository.Delete(ctx, item.ID)
}

type MemoryRepository struct {
	mu   sync.RWMutex
	rows map[string]Media
}

func NewMemoryRepository(initial ...Media) *MemoryRepository {
	rows := make(map[string]Media, len(initial))
	for _, item := range initial {
		if strings.TrimSpace(item.ID) != "" {
			rows[item.ID] = item
		}
	}
	return &MemoryRepository{rows: rows}
}

func (r *MemoryRepository) List(_ context.Context) ([]Media, error) {
	r.mu.RLock()
	items := make([]Media, 0, len(r.rows))
	for _, item := range r.rows {
		items = append(items, item)
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

func (r *MemoryRepository) Get(_ context.Context, id string) (Media, error) {
	r.mu.RLock()
	item, ok := r.rows[id]
	r.mu.RUnlock()
	if !ok {
		return Media{}, ErrNotFound
	}
	return item, nil
}

func (r *MemoryRepository) Create(_ context.Context, item Media) (Media, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.rows[item.ID]; exists {
		return Media{}, ErrInvalidMedia
	}
	r.rows[item.ID] = item
	return item, nil
}

func (r *MemoryRepository) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.rows[id]; !ok {
		return ErrNotFound
	}
	delete(r.rows, id)
	return nil
}
