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
	if err := os.MkdirAll(filepath.Join(root, "admin", "src", "core", "extensions"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "backend", "routes", "web.go"), []byte("package routes\n\nimport (\n)\n\nfunc Web() {\n\tauthController := controllers.NewAuthController()\n}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "backend", "go.mod"), []byte("module goravel\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "admin", "src", "core", "extensions", "runtime.ts"), []byte("export function registerBuiltinExtensions(): void {\n  registerModule(moduleDefinition)\n}\n"), 0o644); err != nil {
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
	if !strings.Contains(string(module), "route: '/admin/resources/products'") {
		t.Fatal("generated resource navigation was not registered in module.ts")
	}
	backendResource, err := os.ReadFile(filepath.Join(root, "modules/catalog/backend/resources/products.go"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(backendResource), "NewGormResourceRepository") || !strings.Contains(string(backendResource), "database *gorm.DB") {
		t.Fatal("generated backend resource does not use the GORM repository")
	}
	routes, err := os.ReadFile(filepath.Join(root, "backend", "routes", "web.go"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(routes), "catalogModule.RegisterRoutes(authController, resourceDatabase.DB)") {
		t.Fatal("generated backend module was not registered in application routes")
	}
	if strings.Index(string(routes), "resourceDatabase, databaseErr :=") > strings.Index(string(routes), "catalogModule.RegisterRoutes") {
		t.Fatal("generated database initialization must precede module route registration")
	}
	runtime, err := os.ReadFile(filepath.Join(root, "admin", "src", "core", "extensions", "runtime.ts"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(runtime), "catalogModuleDefinition") || !strings.Contains(string(runtime), "registerModule(catalogModuleDefinition)") {
		t.Fatal("generated frontend module was not registered in extension runtime")
	}
}

func TestGenerateResourceRegistersMultipleResourcesInOneBackendModule(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "backend", "routes"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, "admin", "src", "core", "extensions"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "backend", "routes", "web.go"), []byte("package routes\n\nimport (\n)\n\nfunc Web() {\n\tauthController := controllers.NewAuthController()\n}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "backend", "go.mod"), []byte("module goravel\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "admin", "src", "core", "extensions", "runtime.ts"), []byte("export function registerBuiltinExtensions(): void {\n  registerModule(moduleDefinition)\n}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := GenerateModule(ModuleOptions{RootDir: root, Name: "catalog"}); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"products", "categories"} {
		if err := GenerateResource(ResourceOptions{RootDir: root, Module: "catalog", Manifest: ResourceManifest{ID: name, Table: name, Label: name, PrimaryKey: "id", Fields: []ManifestField{{Name: "id", Type: "text", Label: "ID"}}}}); err != nil {
			t.Fatal(err)
		}
	}
	module, err := os.ReadFile(filepath.Join(root, "modules", "catalog", "backend", "module.go"))
	if err != nil {
		t.Fatal(err)
	}
	text := string(module)
	for _, marker := range []string{"products \"go-vue-admin-module/catalog/resources/products\"", "categories \"go-vue-admin-module/catalog/resources/categories\"", "products.RegisterRoutes(auth, database)", "categories.RegisterRoutes(auth, database)"} {
		if !strings.Contains(text, marker) {
			t.Errorf("backend module missing %q: %s", marker, text)
		}
	}
	if strings.Count(text, "func RegisterRoutes(") != 1 {
		t.Fatalf("expected one module route registrar, got %d", strings.Count(text, "func RegisterRoutes("))
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
