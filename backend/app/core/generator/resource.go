package generator

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v3"
)

type ResourceOptions struct {
	RootDir  string
	Module   string
	Manifest ResourceManifest
}

type ResourceManifest struct {
	ID          string            `yaml:"id" json:"id"`
	Table       string            `yaml:"table" json:"table"`
	Label       string            `yaml:"label" json:"label"`
	PrimaryKey  string            `yaml:"primary_key" json:"primary_key"`
	SoftDelete  bool              `yaml:"soft_delete,omitempty" json:"soft_delete,omitempty"`
	Audit       bool              `yaml:"audit,omitempty" json:"audit,omitempty"`
	Permissions map[string]string `yaml:"permissions,omitempty" json:"permissions,omitempty"`
	Fields      []ManifestField   `yaml:"fields" json:"fields"`
}

type ManifestField struct {
	Name       string            `yaml:"name" json:"name"`
	Type       string            `yaml:"type" json:"type"`
	Label      string            `yaml:"label" json:"label"`
	Required   bool              `yaml:"required,omitempty" json:"required,omitempty"`
	Searchable bool              `yaml:"searchable,omitempty" json:"searchable,omitempty"`
	Sortable   bool              `yaml:"sortable,omitempty" json:"sortable,omitempty"`
	Options    []ManifestOption  `yaml:"options,omitempty" json:"options,omitempty"`
	Relation   *ManifestRelation `yaml:"relation,omitempty" json:"relation,omitempty"`
}

type ManifestOption struct {
	Value string `yaml:"value" json:"value"`
	Label string `yaml:"label" json:"label"`
}

type ManifestRelation struct {
	Resource     string `yaml:"resource" json:"resource"`
	ForeignKey   string `yaml:"foreign_key" json:"foreign_key"`
	DisplayField string `yaml:"display_field" json:"display_field"`
}

type TableSchema struct {
	Name      string
	Columns   []ColumnSchema
	Relations []RelationSchema
}

type ColumnSchema struct {
	Name       string
	DataType   string
	ColumnType string
	Nullable   bool
	Key        string
	Extra      string
	Comment    string
}

type RelationSchema struct {
	Column           string
	ReferencedTable  string
	ReferencedColumn string
}

var sqlIdentifierPattern = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

func ManifestFromTableSchema(schema TableSchema) ResourceManifest {
	manifest := ResourceManifest{
		ID:         schema.Name,
		Table:      schema.Name,
		Label:      titleWords(schema.Name),
		PrimaryKey: "id",
		Permissions: map[string]string{
			"list":   schema.Name + ".view",
			"get":    schema.Name + ".view",
			"create": schema.Name + ".create",
			"update": schema.Name + ".update",
			"delete": schema.Name + ".delete",
		},
	}
	relations := make(map[string]RelationSchema, len(schema.Relations))
	for _, relation := range schema.Relations {
		relations[relation.Column] = relation
	}
	for _, column := range schema.Columns {
		if column.Name == "deleted_at" {
			manifest.SoftDelete = true
			continue
		}
		if column.Name == "created_at" || column.Name == "updated_at" {
			continue
		}
		if column.Key == "PRI" {
			manifest.PrimaryKey = column.Name
		}
		field := ManifestField{
			Name:       column.Name,
			Type:       manifestFieldType(column),
			Label:      column.Comment,
			Required:   !column.Nullable && column.Extra != "auto_increment",
			Searchable: column.DataType == "char" || column.DataType == "varchar" || column.DataType == "text",
			Sortable:   true,
		}
		if field.Label == "" {
			field.Label = titleWords(column.Name)
		}
		if field.Type == "select" {
			field.Options = enumOptions(column.ColumnType)
		}
		if relation, ok := relations[column.Name]; ok {
			field.Relation = &ManifestRelation{
				Resource:     relation.ReferencedTable,
				ForeignKey:   relation.ReferencedColumn,
				DisplayField: "name",
			}
		}
		manifest.Fields = append(manifest.Fields, field)
	}
	return manifest
}

func GenerateResource(options ResourceOptions) error {
	if !moduleNamePattern.MatchString(options.Module) {
		return fmt.Errorf("module name must match %s: %q", moduleNamePattern.String(), options.Module)
	}
	if !moduleNamePattern.MatchString(options.Manifest.ID) {
		return fmt.Errorf("resource id must match %s: %q", moduleNamePattern.String(), options.Manifest.ID)
	}
	if !sqlIdentifierPattern.MatchString(options.Manifest.Table) {
		return fmt.Errorf("table must be a safe SQL identifier: %q", options.Manifest.Table)
	}
	rootDir := options.RootDir
	if rootDir == "" {
		rootDir = "."
	}
	moduleDir := filepath.Join(rootDir, "modules", options.Module)
	if _, err := os.Stat(filepath.Join(moduleDir, "admin", "module.ts")); err != nil {
		return fmt.Errorf("module admin entry is missing: %w", err)
	}
	resourceDir := filepath.Join(moduleDir, "resources")
	adminDir := filepath.Join(moduleDir, "admin", "resources")
	backendDir := filepath.Join(moduleDir, "backend", "resources")
	manifestPath := filepath.Join(resourceDir, options.Manifest.ID+".yaml")
	adminPath := filepath.Join(adminDir, options.Manifest.ID+".ts")
	backendPath := filepath.Join(backendDir, options.Manifest.ID+".go")
	for _, path := range []string{manifestPath, adminPath, backendPath} {
		if _, err := os.Stat(path); err == nil {
			return fmt.Errorf("resource already exists: %s", options.Manifest.ID)
		} else if !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("inspect %s: %w", path, err)
		}
	}
	manifestBytes, err := yaml.Marshal(options.Manifest)
	if err != nil {
		return fmt.Errorf("marshal resource manifest: %w", err)
	}
	files := map[string][]byte{
		manifestPath: manifestBytes,
		adminPath:    []byte(renderAdminResource(options.Manifest)),
		backendPath:  []byte(renderBackendResource(options.Manifest)),
	}
	for path, content := range files {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return fmt.Errorf("create %s: %w", filepath.Dir(path), err)
		}
		if err := os.WriteFile(path, content, 0o644); err != nil {
			return fmt.Errorf("write %s: %w", path, err)
		}
	}
	if err := registerResourceInModule(filepath.Join(moduleDir, "admin", "module.ts"), options.Manifest); err != nil {
		return err
	}
	if err := registerResourceInBackendModule(filepath.Join(moduleDir, "backend", "module.go"), options.Module, options.Manifest.ID); err != nil {
		return err
	}
	if err := registerResourceInApplicationRoutes(rootDir, options.Module); err != nil {
		return err
	}
	return registerModuleInFrontendRuntime(rootDir, options.Module)
}

func renderAdminResource(manifest ResourceManifest) string {
	typeName := pascalName(manifest.ID)
	var b strings.Builder
	fmt.Fprintln(&b, "import { defineResource } from '@/resource-engine/core/ResourceDefinition'")
	fmt.Fprintln(&b)
	fmt.Fprintf(&b, "export interface %sResourceRecord {\n", typeName)
	for _, field := range manifest.Fields {
		fmt.Fprintf(&b, "  %s: %s\n", field.Name, tsType(field.Type))
	}
	fmt.Fprintln(&b, "}")
	fmt.Fprintln(&b)
	fmt.Fprintf(&b, "export const %sResource = defineResource<%sResourceRecord>({\n", manifest.ID, typeName)
	fmt.Fprintf(&b, "  name: '%s',\n  label: '%s',\n  endpoint: '/api/resources/%s',\n", manifest.ID, manifest.Label, manifest.ID)
	fmt.Fprintln(&b, "  permissions: {")
	for _, action := range []string{"list", "get", "create", "update", "delete"} {
		if permission := manifest.Permissions[action]; permission != "" {
			fmt.Fprintf(&b, "    %s: '%s',\n", action, permission)
		}
	}
	if permission := manifest.Permissions["delete"]; permission != "" {
		fmt.Fprintf(&b, "    bulkDelete: '%s',\n", permission)
	}
	fmt.Fprintln(&b, "  },\n  columns: [")
	for _, field := range manifest.Fields {
		fmt.Fprintf(&b, "    { key: '%s', label: '%s', sortable: %t },\n", field.Name, field.Label, field.Sortable)
	}
	fmt.Fprintln(&b, "  ],\n  fields: [")
	for _, field := range manifest.Fields {
		fmt.Fprintf(&b, "    { name: '%s', label: '%s'", field.Name, field.Label)
		if field.Type != "text" {
			fmt.Fprintf(&b, ", type: '%s'", field.Type)
		}
		if field.Required {
			b.WriteString(", required: true")
		}
		if len(field.Options) > 0 {
			b.WriteString(", options: [")
			for _, option := range field.Options {
				fmt.Fprintf(&b, "{ label: '%s', value: '%s' }, ", option.Label, option.Value)
			}
			b.WriteString("]")
		}
		b.WriteString(" },\n")
	}
	fmt.Fprintln(&b, "  ],\n})")
	return b.String()
}

func renderBackendResource(manifest ResourceManifest) string {
	packageName := strings.ReplaceAll(manifest.ID, "-", "_")
	var b strings.Builder
	fmt.Fprintf(&b, "package %s\n\n", packageName)
	fmt.Fprintln(&b, "import (\n\t\"goravel/app/core/resource\"\n\t\"goravel/app/http/controllers\"\n\t\"gorm.io/gorm\"\n)")
	fmt.Fprintln(&b, "// Code generated by admin-gen from Resource Manifest.")
	fmt.Fprintf(&b, "const (\n\tID = \"%s\"\n\tTable = \"%s\"\n)\n\n", manifest.ID, manifest.Table)
	fmt.Fprint(&b, "type Field struct { Name string; Type string; Label string; Required bool }\n\n")
	fmt.Fprintln(&b, "var Fields = []Field{")
	for _, field := range manifest.Fields {
		fmt.Fprintf(&b, "\t{Name: \"%s\", Type: \"%s\", Label: \"%s\", Required: %t},\n", field.Name, field.Type, field.Label, field.Required)
	}
	fmt.Fprintln(&b, "}")
	fmt.Fprintln(&b)
	fmt.Fprintln(&b, "func RegisterRoutes(auth *controllers.AuthController, database *gorm.DB) error {")
	fmt.Fprintln(&b, "\trepository, err := resource.NewGormResourceRepository(resource.GormResourceOptions{")
	fmt.Fprintln(&b, "\t\tDB: database,")
	fmt.Fprintln(&b, "\t\tSchema: resource.ResourceSchema{")
	fmt.Fprintf(&b, "\t\tTable: %q, PrimaryKey: %q,\n", manifest.Table, manifest.PrimaryKey)
	fmt.Fprint(&b, "\t\tFields: []string{")
	for _, field := range manifest.Fields {
		fmt.Fprintf(&b, "%q, ", field.Name)
	}
	fmt.Fprintln(&b, "},")
	fmt.Fprint(&b, "\t\tSearchable: []string{")
	for _, field := range manifest.Fields {
		if field.Searchable {
			fmt.Fprintf(&b, "%q, ", field.Name)
		}
	}
	fmt.Fprintln(&b, "},")
	fmt.Fprint(&b, "\t\tSortable: []string{")
	for _, field := range manifest.Fields {
		if field.Sortable {
			fmt.Fprintf(&b, "%q, ", field.Name)
		}
	}
	fmt.Fprintln(&b, "},")
	fmt.Fprintf(&b, "\t\tSoftDelete: %t,\n", manifest.SoftDelete)
	fmt.Fprintln(&b, "\t\t},")
	fmt.Fprintln(&b, "\t})")
	fmt.Fprintln(&b, "\tif err != nil { return err }")
	fmt.Fprintf(&b, "\tcontroller := controllers.NewCoreResourceController(auth, resource.NewCoreResourceService(repository), %q)\n", manifest.ID)
	fmt.Fprintf(&b, "\tcontrollers.RegisterCoreResourceRoutes(%q, controller)\n", manifest.ID)
	fmt.Fprintln(&b, "\treturn nil\n}")
	return b.String()
}

func registerResourceInModule(path string, manifest ResourceManifest) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	text := string(content)
	resourceID := manifest.ID
	importLine := fmt.Sprintf("import { %sResource } from './resources/%s'", resourceID, resourceID)
	if !strings.Contains(text, importLine) {
		text = importLine + "\n" + text
	}
	if strings.Contains(text, "resources: []") {
		text = strings.Replace(text, "resources: []", "resources: ["+resourceID+"Resource]", 1)
	} else if !strings.Contains(text, resourceID+"Resource") {
		text = strings.Replace(text, "resources: [", "resources: ["+resourceID+"Resource, ", 1)
	} else if !strings.Contains(text, "resources:") {
		text = strings.Replace(text, "defineAdminModule({", "defineAdminModule({\n  resources: ["+resourceID+"Resource],", 1)
	}
	navigation := fmt.Sprintf("{ id: '%s', label: '%s', route: '/admin/resources/%s', permission: '%s' }", resourceID, manifest.Label, resourceID, manifest.Permissions["list"])
	if strings.Contains(text, "navigation: []") {
		text = strings.Replace(text, "navigation: []", "navigation: ["+navigation+"]", 1)
	} else if !strings.Contains(text, "id: '"+resourceID+"'") {
		text = strings.Replace(text, "navigation: [", "navigation: ["+navigation+", ", 1)
	}
	return os.WriteFile(path, []byte(text), 0o644)
}

func registerResourceInBackendModule(path, module, resourceID string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	text := string(content)
	packageName := strings.ReplaceAll(resourceID, "-", "_")
	resourceImport := fmt.Sprintf("\t%s \"go-vue-admin-module/%s/resources/%s\"", packageName, module, resourceID)
	if !strings.Contains(text, resourceImport) {
		if strings.Contains(text, "import platformmodule \"goravel/app/core/module\"") {
			text = strings.Replace(text, "import platformmodule \"goravel/app/core/module\"", "import (\n\tplatformmodule \"goravel/app/core/module\"\n\t\"goravel/app/http/controllers\"\n"+resourceImport+"\n)", 1)
		} else if index := strings.Index(text, "\n)"); index >= 0 {
			text = text[:index] + "\n" + resourceImport + text[index:]
		}
	}
	if !strings.Contains(text, "\"gorm.io/gorm\"") {
		text = strings.Replace(text, "\t\"goravel/app/http/controllers\"\n", "\t\"goravel/app/http/controllers\"\n\t\"gorm.io/gorm\"\n", 1)
	}
	function := "func RegisterRoutes(auth *controllers.AuthController, database *gorm.DB) error {"
	if !strings.Contains(text, function) {
		text += fmt.Sprintf("\n\n%s\n\tif err := %s.RegisterRoutes(auth, database); err != nil { return err }\n\treturn nil\n}\n", function, packageName)
	} else {
		call := fmt.Sprintf("\tif err := %s.RegisterRoutes(auth, database); err != nil { return err }", packageName)
		if !strings.Contains(text, call) {
			if strings.Contains(text, "\treturn nil\n}") {
				text = strings.Replace(text, "\treturn nil\n}", call+"\n\treturn nil\n}", 1)
			} else {
				oldCall := fmt.Sprintf("\treturn %s.RegisterRoutes(auth, database)", packageName)
				text = strings.Replace(text, oldCall, call+"\n\treturn nil", 1)
			}
		}
	}
	return os.WriteFile(path, []byte(text), 0o644)
}

func registerResourceInApplicationRoutes(rootDir, module string) error {
	routesPath := filepath.Join(rootDir, "backend", "routes", "web.go")
	content, err := os.ReadFile(routesPath)
	if err != nil {
		return fmt.Errorf("read application routes: %w", err)
	}
	text := string(content)
	alias := strings.ReplaceAll(module, "-", "") + "Module"
	importLine := fmt.Sprintf("\t%s \"go-vue-admin-module/%s\"", alias, module)
	if !strings.Contains(text, importLine) {
		text = strings.Replace(text, "import (", "import (\n"+importLine, 1)
	}
	databaseSetup := "\tresourceDatabase := resource.ApplicationResourceDatabase()\n\tif resourceDatabase == nil { panic(resource.ApplicationResourceDatabaseError()) }"
	if !strings.Contains(text, "resourceDatabase := resource.ApplicationResourceDatabase()") {
		text = strings.Replace(text, "\tauthController := controllers.NewAuthController()", "\tauthController := controllers.NewAuthController()\n"+databaseSetup, 1)
	}
	call := fmt.Sprintf("\tif err := %s.RegisterRoutes(authController, resourceDatabase.DB); err != nil {\n\t\tpanic(err)\n\t}", alias)
	if !strings.Contains(text, call) {
		text = strings.Replace(text, databaseSetup, databaseSetup+"\n"+call, 1)
	}
	if err := os.WriteFile(routesPath, []byte(text), 0o644); err != nil {
		return err
	}
	goModPath := filepath.Join(rootDir, "backend", "go.mod")
	goMod, err := os.ReadFile(goModPath)
	if err != nil {
		return fmt.Errorf("read backend go.mod: %w", err)
	}
	modulePath := fmt.Sprintf("go-vue-admin-module/%s", module)
	goModText := string(goMod)
	if !strings.Contains(goModText, modulePath) {
		goModText += fmt.Sprintf("\nrequire %s v0.0.0\nreplace %s => ../modules/%s/backend\n", modulePath, modulePath, module)
		if err := os.WriteFile(goModPath, []byte(goModText), 0o644); err != nil {
			return err
		}
	}
	return nil
}

func registerModuleInFrontendRuntime(rootDir, module string) error {
	runtimePath := filepath.Join(rootDir, "admin", "src", "core", "extensions", "runtime.ts")
	content, err := os.ReadFile(runtimePath)
	if err != nil {
		return fmt.Errorf("read frontend extension runtime: %w", err)
	}
	text := string(content)
	alias := strings.ReplaceAll(module, "-", "") + "ModuleDefinition"
	importLine := fmt.Sprintf("import { moduleDefinition as %s } from '../../../../modules/%s/admin/module'", alias, module)
	if !strings.Contains(text, importLine) {
		text = importLine + "\n" + text
	}
	registration := fmt.Sprintf("  registerModule(%s)", alias)
	if !strings.Contains(text, registration) {
		text = strings.Replace(text, "  registerModule(moduleDefinition)\n", "  registerModule(moduleDefinition)\n"+registration+"\n", 1)
	}
	return os.WriteFile(runtimePath, []byte(text), 0o644)
}

func manifestFieldType(column ColumnSchema) string {
	dataType := strings.ToLower(column.DataType)
	switch dataType {
	case "enum", "set":
		return "select"
	case "bool", "boolean":
		return "checkbox"
	case "tinyint":
		if strings.EqualFold(column.ColumnType, "tinyint(1)") {
			return "checkbox"
		}
	case "int", "integer", "bigint", "smallint", "mediumint", "decimal", "numeric", "float", "double":
		return "number"
	case "date":
		return "date"
	case "datetime", "timestamp", "time":
		return "datetime"
	case "text", "mediumtext", "longtext", "blob", "mediumblob", "longblob":
		return "textarea"
	}
	return "text"
}

func enumOptions(columnType string) []ManifestOption {
	start := strings.Index(strings.ToLower(columnType), "enum(")
	if start < 0 || !strings.HasSuffix(columnType, ")") {
		return nil
	}
	valueText := columnType[start+5 : len(columnType)-1]
	parts := strings.Split(valueText, "','")
	options := make([]ManifestOption, 0, len(parts))
	for _, part := range parts {
		value := strings.Trim(part, " '")
		if value != "" {
			options = append(options, ManifestOption{Value: value, Label: titleWords(strings.ReplaceAll(value, "_", "-"))})
		}
	}
	return options
}

func tsType(fieldType string) string {
	switch fieldType {
	case "number":
		return "number"
	case "checkbox", "switch":
		return "boolean"
	default:
		return "string"
	}
}

func pascalName(value string) string {
	return strings.ReplaceAll(titleWords(strings.ReplaceAll(value, "_", "-")), " ", "")
}

type MySQLIntrospector struct {
	db *sql.DB
}

func NewMySQLIntrospector(_ context.Context, user, password, address, database string) *MySQLIntrospector {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true&charset=utf8mb4", user, password, address, database)
	db, _ := sql.Open("mysql", dsn)
	return &MySQLIntrospector{db: db}
}

func (i *MySQLIntrospector) Close() error {
	if i == nil || i.db == nil {
		return nil
	}
	return i.db.Close()
}

func (i *MySQLIntrospector) Inspect(ctx context.Context, table string) (TableSchema, error) {
	if !sqlIdentifierPattern.MatchString(table) {
		return TableSchema{}, fmt.Errorf("unsafe table name %q", table)
	}
	if i == nil || i.db == nil {
		return TableSchema{}, errors.New("mysql database is not configured")
	}
	rows, err := i.db.QueryContext(ctx, `SELECT COLUMN_NAME, DATA_TYPE, COLUMN_TYPE, IS_NULLABLE, COLUMN_KEY, EXTRA, COLUMN_COMMENT FROM information_schema.columns WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION`, table)
	if err != nil {
		return TableSchema{}, err
	}
	defer rows.Close()
	schema := TableSchema{Name: table}
	for rows.Next() {
		var column ColumnSchema
		var nullable, comment string
		if err := rows.Scan(&column.Name, &column.DataType, &column.ColumnType, &nullable, &column.Key, &column.Extra, &comment); err != nil {
			return TableSchema{}, err
		}
		column.Nullable = nullable == "YES"
		column.Comment = comment
		schema.Columns = append(schema.Columns, column)
	}
	if err := rows.Err(); err != nil {
		return TableSchema{}, err
	}
	relations, err := i.db.QueryContext(ctx, `SELECT COLUMN_NAME, REFERENCED_TABLE_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? AND REFERENCED_TABLE_NAME IS NOT NULL`, table)
	if err != nil {
		return TableSchema{}, err
	}
	defer relations.Close()
	for relations.Next() {
		var relation RelationSchema
		if err := relations.Scan(&relation.Column, &relation.ReferencedTable, &relation.ReferencedColumn); err != nil {
			return TableSchema{}, err
		}
		schema.Relations = append(schema.Relations, relation)
	}
	return schema, relations.Err()
}

func sortFields(fields []ManifestField) {
	sort.SliceStable(fields, func(i, j int) bool { return fields[i].Name < fields[j].Name })
}
