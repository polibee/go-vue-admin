package controllers

import (
	"testing"

	"goravel/app/core/resource"
)

func TestValidateGeneratedValuesEnforcesRequiredAndSelectFields(t *testing.T) {
	manifest := resource.Manifest{Fields: []resource.Field{
		{Name: "title", Type: "text", Required: true},
		{Name: "status", Type: "select", Required: true, Options: []resource.Option{{Value: "active"}, {Value: "disabled"}}},
		{Name: "enabled", Type: "boolean"},
	}}

	if _, err := validateGeneratedValues(map[string]any{"title": "Post", "status": "active"}, manifest); err != nil {
		t.Fatalf("valid values rejected: %v", err)
	}
	if _, err := validateGeneratedValues(map[string]any{"title": "Post", "status": "unknown"}, manifest); err == nil {
		t.Fatal("invalid select option accepted")
	}
	if _, err := validateGeneratedValues(map[string]any{"status": "active"}, manifest); err == nil {
		t.Fatal("missing required field accepted")
	}
}
