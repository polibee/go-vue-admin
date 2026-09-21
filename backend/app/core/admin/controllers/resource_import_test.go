package controllers

import (
	"strings"
	"testing"

	"goravel/app/core/resource"
)

func TestParseResourceCSVRejectsUnknownAndSensitiveColumns(t *testing.T) {
	manifest := resource.Manifest{
		Fields:  []resource.Field{{Name: "name", Type: "text", Required: true}, {Name: "password", Type: "password"}},
		Columns: []resource.Column{{Name: "name", Label: "Name"}},
	}
	result := parseResourceCSV(strings.NewReader("name,password,unknown\nAlice,secret,x\n"), manifest)
	if !containsString(result.Errors, "unsupported column: password") || !containsString(result.Errors, "unsupported column: unknown") {
		t.Fatalf("unexpected parse errors: %#v", result.Errors)
	}
}

func TestParseResourceCSVReportsMissingRequiredValue(t *testing.T) {
	manifest := resource.Manifest{
		Fields:  []resource.Field{{Name: "name", Type: "text", Required: true}},
		Columns: []resource.Column{{Name: "name", Label: "Name"}},
	}
	result := parseResourceCSV(strings.NewReader("name\n \n"), manifest)
	if len(result.RowErrors) != 1 || result.RowErrors[0].Row != 2 {
		t.Fatalf("unexpected row errors: %#v", result.RowErrors)
	}
}

func containsString(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}
