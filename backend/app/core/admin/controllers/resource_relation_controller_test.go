package controllers

import "testing"

func TestNormalizeRelationOptionsQuery(t *testing.T) {
	query := normalizeRelationOptionsQuery("  design  ", "7", "3", "25")
	if query.Search != "design" || query.Selected != "7" || query.Page != 3 || query.PerPage != 25 {
		t.Fatalf("query = %+v, want normalized search and pagination", query)
	}
}

func TestNormalizeRelationOptionsQueryBoundsInvalidValues(t *testing.T) {
	query := normalizeRelationOptionsQuery("", "", "0", "1000")
	if query.Page != 1 || query.PerPage != 100 {
		t.Fatalf("query = %+v, want page 1 and per_page 100", query)
	}
}
