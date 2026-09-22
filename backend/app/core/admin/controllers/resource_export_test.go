package controllers

import (
	"testing"

	"goravel/app/core/resource"
)

func TestCSVCellEscapesRFC4180Characters(t *testing.T) {
	got, err := renderResourceCSV([]resource.Column{{Name: "name", Label: "Name"}}, []map[string]any{{"name": "a,b\"c\nd"}})
	if err != nil {
		t.Fatalf("renderResourceCSV() error = %v", err)
	}
	if want := "Name\n\"a,b\"\"c\nd\"\n"; got != want {
		t.Fatalf("renderResourceCSV() = %q, want %q", got, want)
	}
}

func TestExportColumnsUseManifestColumnsAndIDOnly(t *testing.T) {
	manifest := resource.Manifest{
		Columns: []resource.Column{
			{Name: "name", Label: "Name"},
			{Name: "password", Label: "Password"},
			{Name: "status", Label: "Status"},
		},
	}
	columns := exportColumns(manifest)
	if got, want := len(columns), 3; got != want {
		t.Fatalf("export columns = %d, want %d", got, want)
	}
	if columns[0].Name != "id" || columns[1].Name != "name" || columns[2].Name != "status" {
		t.Fatalf("unexpected export columns: %#v", columns)
	}
}

func TestExportColumnsExcludeSensitiveAndUnreadableManifestFields(t *testing.T) {
	manifest := resource.Manifest{
		Fields: []resource.Field{
			{Name: "name", Type: "text", Visible: true, Readable: true, Writable: true, PolicyConfigured: true},
			{Name: "token", Type: "text", Visible: true, Readable: true, Writable: true, Sensitive: true, PolicyConfigured: true},
			{Name: "internal", Type: "text", Visible: false, Readable: false, Writable: false, PolicyConfigured: true},
		},
		Columns: []resource.Column{
			{Name: "name", Label: "Name"},
			{Name: "token", Label: "Token"},
			{Name: "internal", Label: "Internal"},
		},
	}
	columns := exportColumns(manifest)
	if got := joinColumnNames(columns); got != "id,name" {
		t.Fatalf("restricted export columns = %q, want %q", got, "id,name")
	}
}

func joinColumnNames(columns []resource.Column) string {
	result := ""
	for index, column := range columns {
		if index > 0 {
			result += ","
		}
		result += column.Name
	}
	return result
}
