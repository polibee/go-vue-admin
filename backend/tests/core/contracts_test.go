package core

import (
	"net/url"
	"testing"

	apierrors "goravel/app/core/shared/errors"
	"goravel/app/core/shared/pagination"
	"goravel/app/core/shared/query"
	"goravel/app/core/shared/response"
)

func TestListQueryParsesTheSharedProtocol(t *testing.T) {
	values := url.Values{
		"page":         {"2"},
		"per_page":     {"25"},
		"search":       {"  admin  "},
		"sort":         {"created_at"},
		"sort_dir":     {"desc"},
		"filter[name]": {"alice"},
	}

	parsed, err := query.Parse(values)
	if err != nil {
		t.Fatalf("parse query: %v", err)
	}
	if parsed.Pagination.Page != 2 || parsed.Pagination.PerPage != 25 {
		t.Fatalf("unexpected pagination: %+v", parsed.Pagination)
	}
	if parsed.Search.Term != "admin" || parsed.Sort.Field != "created_at" || !parsed.Sort.Desc {
		t.Fatalf("unexpected search/sort: %+v", parsed)
	}
	if parsed.Filter.Values["name"] != "alice" {
		t.Fatalf("unexpected filters: %+v", parsed.Filter.Values)
	}
}

func TestListQueryRejectsInvalidPagination(t *testing.T) {
	if _, err := query.Parse(url.Values{"per_page": {"101"}}); err == nil {
		t.Fatal("expected per_page validation error")
	}
}

func TestPaginationMetaCalculatesTotalPages(t *testing.T) {
	meta := pagination.NewMeta(pagination.Query{Page: 2, PerPage: 20}, 41)
	if meta.TotalPages != 3 || meta.Page != 2 || meta.Total != 41 {
		t.Fatalf("unexpected pagination meta: %+v", meta)
	}
}

func TestResponseContractsKeepStableEnvelopeKeys(t *testing.T) {
	success := response.Success(map[string]string{"status": "ok"}, response.Meta{RequestID: "req-1"})
	if success.Meta.RequestID != "req-1" || success.Data == nil {
		t.Fatalf("unexpected success envelope: %+v", success)
	}

	failure := apierrors.InvalidQuery(testError("bad page"))
	if failure.Error.Code != "INVALID_QUERY" || failure.Error.Details == nil {
		t.Fatalf("unexpected error envelope: %+v", failure)
	}
}

type testError string

func (e testError) Error() string { return string(e) }
