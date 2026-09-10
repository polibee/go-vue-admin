package audit

import (
	"context"
	"encoding/json"
	"time"

	"github.com/goravel/framework/contracts/database/db"
	"github.com/goravel/framework/facades"
)

const auditTable = "audit_entries"

type auditQuery interface {
	Where(any, ...any) auditQuery
	OrderBy(string, ...string) auditQuery
	Get(any) error
	Insert(any) error
}

type auditQueryFactory func(context.Context) auditQuery

type MySQLRepository struct{ query auditQueryFactory }

func NewMySQLRepository() Repository {
	return newMySQLRepositoryWithQuery(func(ctx context.Context) auditQuery {
		return &goravelAuditQuery{query: facades.DB().WithContext(ctx).Table(auditTable)}
	})
}

func newMySQLRepositoryWithQuery(factory auditQueryFactory) *MySQLRepository {
	return &MySQLRepository{query: factory}
}

func (r *MySQLRepository) List(ctx context.Context) ([]Entry, error) {
	var rows []auditRow
	if err := r.query(ctx).OrderBy("created_at", "desc").Get(&rows); err != nil {
		return nil, err
	}
	entries := make([]Entry, 0, len(rows))
	for _, row := range rows {
		entry, err := row.entry()
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

func (r *MySQLRepository) Get(ctx context.Context, id string) (Entry, error) {
	var rows []auditRow
	if err := r.query(ctx).Where("id", id).Get(&rows); err != nil {
		return Entry{}, err
	}
	if len(rows) == 0 {
		return Entry{}, ErrNotFound
	}
	return rows[0].entry()
}

func (r *MySQLRepository) Create(ctx context.Context, entry Entry) (Entry, error) {
	before, err := json.Marshal(entry.Before)
	if err != nil {
		return Entry{}, err
	}
	after, err := json.Marshal(entry.After)
	if err != nil {
		return Entry{}, err
	}
	row := map[string]any{
		"id": entry.ID, "actor_id": entry.ActorID, "actor_email": entry.ActorEmail,
		"action": entry.Action, "resource_type": entry.ResourceType, "resource_id": entry.ResourceID,
		"before_state": string(before), "after_state": string(after), "ip": entry.IP,
		"user_agent": entry.UserAgent, "created_at": entry.CreatedAt,
	}
	if err := r.query(ctx).Insert(row); err != nil {
		return Entry{}, err
	}
	return r.Get(ctx, entry.ID)
}

type auditRow struct {
	ID           string    `mapstructure:"id"`
	ActorID      string    `mapstructure:"actor_id"`
	ActorEmail   string    `mapstructure:"actor_email"`
	Action       string    `mapstructure:"action"`
	ResourceType string    `mapstructure:"resource_type"`
	ResourceID   string    `mapstructure:"resource_id"`
	BeforeState  string    `mapstructure:"before_state"`
	AfterState   string    `mapstructure:"after_state"`
	IP           string    `mapstructure:"ip"`
	UserAgent    string    `mapstructure:"user_agent"`
	CreatedAt    time.Time `mapstructure:"created_at"`
}

func (r auditRow) entry() (Entry, error) {
	var before, after any
	if r.BeforeState != "" {
		if err := json.Unmarshal([]byte(r.BeforeState), &before); err != nil {
			return Entry{}, err
		}
	}
	if r.AfterState != "" {
		if err := json.Unmarshal([]byte(r.AfterState), &after); err != nil {
			return Entry{}, err
		}
	}
	createdAt := ""
	if !r.CreatedAt.IsZero() {
		createdAt = r.CreatedAt.UTC().Format(time.RFC3339)
	}
	return Entry{ID: r.ID, ActorID: r.ActorID, ActorEmail: r.ActorEmail, Action: r.Action, ResourceType: r.ResourceType, ResourceID: r.ResourceID, Before: before, After: after, IP: r.IP, UserAgent: r.UserAgent, CreatedAt: createdAt}, nil
}

type goravelAuditQuery struct{ query db.Query }

func (q *goravelAuditQuery) Where(value any, args ...any) auditQuery {
	return &goravelAuditQuery{query: q.query.Where(value, args...)}
}

func (q *goravelAuditQuery) OrderBy(column string, directions ...string) auditQuery {
	return &goravelAuditQuery{query: q.query.OrderBy(column, directions...)}
}

func (q *goravelAuditQuery) Get(dest any) error { return q.query.Get(dest) }

func (q *goravelAuditQuery) Insert(data any) error {
	_, err := q.query.Insert(data)
	return err
}
