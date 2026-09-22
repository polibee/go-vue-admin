package actions

import (
	"errors"
	"sort"

	"github.com/goravel/framework/contracts/http"
)

var (
	ErrEmptyIDs          = errors.New("action IDs are required")
	ErrInvalidID         = errors.New("action ID must be positive")
	ErrTooManyIDs        = errors.New("too many action IDs")
	ErrDuplicateHandler  = errors.New("action handler kind already registered")
	ErrHandlerNotFound   = errors.New("action handler not found")
	ErrActionNotBatch    = errors.New("action does not support batch execution")
	ErrPayloadContract   = errors.New("action payload contract is invalid")
	ErrSelectionContract = errors.New("selection contract is invalid")
)

const MaxIDs = 100

type Request struct {
	Action  string
	IDs     []int64
	Payload map[string]any
}

// Selection describes the records targeted by a batch operation. IDs are
// explicit and query mode is re-evaluated by the server under the caller's
// current resource scope.
type Selection struct {
	Mode       string            `json:"mode"`
	IDs        []int64           `json:"ids,omitempty"`
	Query      map[string]string `json:"query,omitempty"`
	ExcludeIDs []int64           `json:"exclude_ids,omitempty"`
}

func (s Selection) Validate() error {
	switch s.Mode {
	case "ids":
		if len(s.IDs) == 0 || len(s.Query) != 0 {
			return ErrSelectionContract
		}
	case "query":
		if len(s.Query) == 0 && len(s.IDs) == 0 {
			return ErrSelectionContract
		}
	default:
		return ErrSelectionContract
	}
	return nil
}

type Failure struct {
	ID   int64  `json:"id"`
	Code string `json:"code"`
}

type Result struct {
	Action    string    `json:"action"`
	Requested int       `json:"requested"`
	Succeeded int       `json:"succeeded"`
	Failed    int       `json:"failed"`
	Skipped   int       `json:"skipped"`
	Failures  []Failure `json:"failures,omitempty"`
	Skips     []Failure `json:"skips,omitempty"`
}

type Handler interface {
	Kind() string
	Payload() string
	Execute(ctx http.Context, request Request) (Result, error)
}

type Registry struct {
	handlers map[string]Handler
}

func NewRegistry() *Registry { return &Registry{handlers: make(map[string]Handler)} }

func (r *Registry) Register(handler Handler) error {
	if handler == nil || handler.Kind() == "" {
		return ErrHandlerNotFound
	}
	if _, exists := r.handlers[handler.Kind()]; exists {
		return ErrDuplicateHandler
	}
	r.handlers[handler.Kind()] = handler
	return nil
}

func (r *Registry) Find(kind string) (Handler, error) {
	handler, ok := r.handlers[kind]
	if !ok {
		return nil, ErrHandlerNotFound
	}
	return handler, nil
}

func NormalizeRequest(request Request) (Request, error) {
	if len(request.IDs) == 0 {
		return Request{}, ErrEmptyIDs
	}
	seen := make(map[int64]struct{}, len(request.IDs))
	ids := make([]int64, 0, len(request.IDs))
	for _, id := range request.IDs {
		if id < 1 {
			return Request{}, ErrInvalidID
		}
		if _, exists := seen[id]; exists {
			continue
		}
		seen[id] = struct{}{}
		ids = append(ids, id)
	}
	if len(ids) > MaxIDs {
		return Request{}, ErrTooManyIDs
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
	request.IDs = ids
	return request, nil
}
