package generator

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestManifestFromTableSchemaInfersFieldsAndEnumOptions(t *testing.T) {
	manifest := ManifestFromTableSchema(TableSchema{
		Name: "products",
		Columns: []ColumnSchema{
			{Name: "id", DataType: "bigint", ColumnType: "bigint", Nullable: false, Key: "PRI", Extra: "auto_increment"},
			{Name: "name", DataType: "varchar", ColumnType: "varchar(191)", Nullable: false, Comment: "产品名称"},
			{Name: "status", DataType: "enum", ColumnType: "enum('draft','active')", Nullable: false},
			{Name: "deleted_at", DataType: "datetime", ColumnType: "datetime", Nullable: true},
		},
	})

	if manifest.ID != "products" || manifest.Table != "products" || manifest.PrimaryKey != "id" {
		t.Fatalf("unexpected identity: %#v", manifest)
	}
	if !manifest.SoftDelete {
		t.Fatal("expected deleted_at to enable soft delete")
	}
	status := manifest.field("status")
	if status == nil || status.Type != "select" || len(status.Options) != 2 {
		t.Fatalf("unexpected status field: %#v", status)
	}
	if name := manifest.field("name"); name == nil || name.Label != "产品名称" || !name.Required {
		t.Fatalf("unexpected name field: %#v", name)
	}
}

func TestGenerateResourceRendersManifestAndModuleFiles(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "backend", "routes"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "backend", "routes", "web.go"), []byte("package routes\n\nimport (\n)\n\nfunc Web() {\n\tauthController := controllers.NewAuthController()\n}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "backend", "go.mod"), []byte("module goravel\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := GenerateModule(ModuleOptions{RootDir: root, Name: "catalog"}); err != nil {
		t.Fatal(err)
	}
	manifest := ManifestFromTableSchema(TableSchema{
		Name:    "products",
		Columns: []ColumnSchema{{Name: "id", DataType: "varchar", ColumnType: "varchar(64)", Nullable: false, Key: "PRI"}, {Name: "name", DataType: "varchar", ColumnType: "varchar(191)", Nullable: false}},
	})
	if err := GenerateResource(ResourceOptions{RootDir: root, Module: "catalog", Manifest: manifest}); err != nil {
		t.Fatal(err)
	}

	for _, file := range []string{
		"modules/catalog/resources/products.yaml",
		"modules/catalog/backend/resources/products.go",
		"modules/catalog/admin/resources/products.ts",
	} {
		if _, err := os.Stat(filepath.Join(root, file)); err != nil {
			t.Errorf("generated file %s: %v", file, err)
		}
	}
	module, err := os.ReadFile(filepath.Join(root, "modules/catalog/admin/module.ts"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(module), "productsResource") {
		t.Fatal("generated resource was not registered in module.ts")
	}
}

func TestMySQLIntrospectorRequiresSafeTableName(t *testing.T) {
	if _, err := NewMySQLIntrospector(context.Background(), "root", "root", "127.0.0.1:3306", "db").Inspect(context.Background(), "products;drop"); err == nil {
		t.Fatal("expected unsafe table name to be rejected")
	}
}

func (manifest ResourceManifest) field(name string) *ManifestField {
	for index := range manifest.Fields {
		if manifest.Fields[index].Name == name {
			return &manifest.Fields[index]
		}
	}
	return nil
}
