package setting

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/goravel/framework/contracts/database/db"
	"github.com/goravel/framework/facades"
)

const settingsTable = "settings"

type settingQuery interface {
	Where(any, ...any) settingQuery
	OrderBy(string, ...string) settingQuery
	Get(any) error
	Insert(any) error
	Update(any) error
	Delete() error
}

type settingQueryFactory func(context.Context) settingQuery

type MySQLRepository struct{ query settingQueryFactory }

func NewMySQLRepository() Repository {
	return newMySQLRepositoryWithQuery(func(ctx context.Context) settingQuery {
		return &goravelSettingQuery{query: facades.DB().WithContext(ctx).Table(settingsTable)}
	})
}

func newMySQLRepositoryWithQuery(factory settingQueryFactory) *MySQLRepository {
	return &MySQLRepository{query: factory}
}

func (r *MySQLRepository) List(ctx context.Context, namespace string) ([]Setting, error) {
	query := r.query(ctx).OrderBy("namespace").OrderBy("key")
	if namespace != "" {
		query = query.Where("namespace", namespace)
	}
	var rows []settingRow
	if err := query.Get(&rows); err != nil {
		return nil, err
	}
	items := make([]Setting, 0, len(rows))
	for _, row := range rows {
		item, err := row.setting()
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *MySQLRepository) Get(ctx context.Context, namespace, key string) (Setting, error) {
	var rows []settingRow
	if err := r.query(ctx).Where("namespace", namespace).Where("key", key).Get(&rows); err != nil {
		return Setting{}, err
	}
	if len(rows) == 0 {
		return Setting{}, ErrNotFound
	}
	return rows[0].setting()
}

func (r *MySQLRepository) Upsert(ctx context.Context, input Setting) (Setting, error) {
	value, err := encodeValue(input.Value)
	if err != nil {
		return Setting{}, err
	}
	row := map[string]any{
		"namespace": input.Namespace, "key": input.Key, "value": value,
		"value_type": input.ValueType, "description": input.Description,
	}
	if _, err := r.Get(ctx, input.Namespace, input.Key); err == ErrNotFound {
		if err := r.query(ctx).Insert(row); err != nil {
			return Setting{}, err
		}
	} else if err != nil {
		return Setting{}, err
	} else if err := r.query(ctx).Where("namespace", input.Namespace).Where("key", input.Key).Update(row); err != nil {
		return Setting{}, err
	}
	return r.Get(ctx, input.Namespace, input.Key)
}

func (r *MySQLRepository) Delete(ctx context.Context, namespace, key string) error {
	if _, err := r.Get(ctx, namespace, key); err != nil {
		return err
	}
	return r.query(ctx).Where("namespace", namespace).Where("key", key).Delete()
}

type settingRow struct {
	Namespace   string    `mapstructure:"namespace"`
	Key         string    `mapstructure:"key"`
	Value       string    `mapstructure:"value"`
	ValueType   ValueType `mapstructure:"value_type"`
	Description string    `mapstructure:"description"`
	UpdatedAt   time.Time `mapstructure:"updated_at"`
}

func (r settingRow) setting() (Setting, error) {
	value, err := decodeValue(r.Value, r.ValueType)
	if err != nil {
		return Setting{}, err
	}
	updatedAt := ""
	if !r.UpdatedAt.IsZero() {
		updatedAt = r.UpdatedAt.UTC().Format(time.RFC3339)
	}
	return Setting{Namespace: r.Namespace, Key: r.Key, Value: value, ValueType: r.ValueType, Description: r.Description, UpdatedAt: updatedAt}, nil
}

func encodeValue(value any) (string, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return "", fmt.Errorf("encode setting value: %w", err)
	}
	return string(data), nil
}

func decodeValue(value string, valueType ValueType) (any, error) {
	if valueType == ValueTypeString {
		var result string
		if err := json.Unmarshal([]byte(value), &result); err == nil {
			return result, nil
		}
		return value, nil
	}
	var result any
	if err := json.Unmarshal([]byte(value), &result); err != nil {
		return nil, fmt.Errorf("decode setting value: %w", err)
	}
	if valueType == ValueTypeInteger {
		if number, ok := result.(float64); ok {
			return int64(number), nil
		}
	}
	return result, nil
}

type goravelSettingQuery struct{ query db.Query }

func (q *goravelSettingQuery) Where(value any, args ...any) settingQuery {
	return &goravelSettingQuery{query: q.query.Where(value, args...)}
}

func (q *goravelSettingQuery) OrderBy(column string, directions ...string) settingQuery {
	return &goravelSettingQuery{query: q.query.OrderBy(column, directions...)}
}

func (q *goravelSettingQuery) Get(dest any) error { return q.query.Get(dest) }

func (q *goravelSettingQuery) Insert(data any) error {
	_, err := q.query.Insert(data)
	return err
}

func (q *goravelSettingQuery) Update(data any) error {
	_, err := q.query.Update(data)
	return err
}

func (q *goravelSettingQuery) Delete() error {
	_, err := q.query.Delete()
	return err
}
