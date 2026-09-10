package resource

import (
	"context"
	"testing"
)

func TestMySQLCoreResourceRepositoryRejectsUnsafeDefinition(t *testing.T) {
	if _, err := NewMySQLCoreResourceRepository(MySQLCoreResourceOptions{Table: "products;drop", PrimaryKey: "id"}); err == nil {
		t.Fatal("expected unsafe table name to be rejected")
	}
}

func TestMySQLCoreResourceRepositoryFiltersInputToManifestFields(t *testing.T) {
	query := &fakeCoreResourceQuery{}
	repository, err := newMySQLCoreResourceRepositoryWithQuery(MySQLCoreResourceOptions{
		Table:      "products",
		PrimaryKey: "id",
		Fields:     []string{"id", "name", "status"},
	}, func(context.Context) coreResourceQuery { return query })
	if err != nil {
		t.Fatal(err)
	}
	if _, err := repository.Create(context.Background(), CoreResourceRecord{
		"id": "p-1", "name": "Demo", "status": "draft", "secret": "must-not-persist",
	}); err == nil {
		t.Fatal("expected fake query to stop after insert/get boundary")
	}
	if _, ok := query.inserted["secret"]; ok {
		t.Fatal("unlisted field was passed to insert")
	}
}

type fakeCoreResourceQuery struct {
	inserted map[string]any
}

func (q *fakeCoreResourceQuery) Where(any, ...any) coreResourceQuery         { return q }
func (q *fakeCoreResourceQuery) WhereAny([]string, ...any) coreResourceQuery { return q }
func (q *fakeCoreResourceQuery) WhereNull(string) coreResourceQuery          { return q }
func (q *fakeCoreResourceQuery) WhereIn(string, []any) coreResourceQuery     { return q }
func (q *fakeCoreResourceQuery) OrderBy(string, ...string) coreResourceQuery { return q }
func (q *fakeCoreResourceQuery) OrderByDesc(string) coreResourceQuery        { return q }
func (q *fakeCoreResourceQuery) Paginate(int, int, any, *int64) error        { return nil }
func (q *fakeCoreResourceQuery) Get(any) error                               { return nil }
func (q *fakeCoreResourceQuery) Insert(value any) error {
	q.inserted = value.(CoreResourceRecord)
	return nil
}
func (q *fakeCoreResourceQuery) Update(any) error { return nil }
func (q *fakeCoreResourceQuery) Delete() error    { return nil }
