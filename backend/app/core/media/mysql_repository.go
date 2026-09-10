package media

import (
	"context"
	"time"

	"github.com/goravel/framework/contracts/database/db"
	"github.com/goravel/framework/facades"
)

const mediaTable = "media"

type mediaQuery interface {
	Where(any, ...any) mediaQuery
	OrderBy(string, ...string) mediaQuery
	Get(any) error
	Insert(any) error
	Delete() error
}

type mediaQueryFactory func(context.Context) mediaQuery

type MySQLRepository struct{ query mediaQueryFactory }

func NewMySQLRepository() Repository {
	return newMySQLRepositoryWithQuery(func(ctx context.Context) mediaQuery {
		return &goravelMediaQuery{query: facades.DB().WithContext(ctx).Table(mediaTable)}
	})
}

func newMySQLRepositoryWithQuery(factory mediaQueryFactory) *MySQLRepository {
	return &MySQLRepository{query: factory}
}

func (r *MySQLRepository) List(ctx context.Context) ([]Media, error) {
	var rows []mediaRow
	if err := r.query(ctx).OrderBy("created_at", "desc").Get(&rows); err != nil {
		return nil, err
	}
	items := make([]Media, 0, len(rows))
	for _, row := range rows {
		items = append(items, row.media())
	}
	return items, nil
}

func (r *MySQLRepository) Get(ctx context.Context, id string) (Media, error) {
	var rows []mediaRow
	if err := r.query(ctx).Where("id", id).Get(&rows); err != nil {
		return Media{}, err
	}
	if len(rows) == 0 {
		return Media{}, ErrNotFound
	}
	return rows[0].media(), nil
}

func (r *MySQLRepository) Create(ctx context.Context, item Media) (Media, error) {
	row := map[string]any{
		"id": item.ID, "disk": item.Disk, "path": item.Path,
		"original_name": item.OriginalName, "mime_type": item.MIMEType,
		"size": item.Size, "created_at": item.CreatedAt,
	}
	if err := r.query(ctx).Insert(row); err != nil {
		return Media{}, err
	}
	return r.Get(ctx, item.ID)
}

func (r *MySQLRepository) Delete(ctx context.Context, id string) error {
	if _, err := r.Get(ctx, id); err != nil {
		return err
	}
	return r.query(ctx).Where("id", id).Delete()
}

type mediaRow struct {
	ID           string    `mapstructure:"id"`
	Disk         string    `mapstructure:"disk"`
	Path         string    `mapstructure:"path"`
	OriginalName string    `mapstructure:"original_name"`
	MIMEType     string    `mapstructure:"mime_type"`
	Size         int64     `mapstructure:"size"`
	CreatedAt    time.Time `mapstructure:"created_at"`
}

func (r mediaRow) media() Media {
	createdAt := ""
	if !r.CreatedAt.IsZero() {
		createdAt = r.CreatedAt.UTC().Format(time.RFC3339)
	}
	return Media{ID: r.ID, Disk: r.Disk, Path: r.Path, OriginalName: r.OriginalName, MIMEType: r.MIMEType, Size: r.Size, CreatedAt: createdAt}
}

type goravelMediaQuery struct{ query db.Query }

func (q *goravelMediaQuery) Where(value any, args ...any) mediaQuery {
	return &goravelMediaQuery{query: q.query.Where(value, args...)}
}

func (q *goravelMediaQuery) OrderBy(column string, directions ...string) mediaQuery {
	return &goravelMediaQuery{query: q.query.OrderBy(column, directions...)}
}

func (q *goravelMediaQuery) Get(dest any) error { return q.query.Get(dest) }

func (q *goravelMediaQuery) Insert(data any) error {
	_, err := q.query.Insert(data)
	return err
}

func (q *goravelMediaQuery) Delete() error {
	_, err := q.query.Delete()
	return err
}
