package routes

import "testing"

func TestAPISurfacePrefixes(t *testing.T) {
	if adminAPIBase != "/api/v1/admin" {
		t.Fatalf("admin API base = %q", adminAPIBase)
	}
	if appAPIBase != "/api/v1/app" {
		t.Fatalf("app API base = %q", appAPIBase)
	}
}
