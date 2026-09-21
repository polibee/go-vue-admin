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
