package query

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"goravel/app/core/shared/pagination"
)

type FilterQuery struct {
	Values map[string]string `json:"values,omitempty"`
}

type SortQuery struct {
	Field string `json:"field,omitempty"`
	Desc  bool   `json:"desc,omitempty"`
}

type SearchQuery struct {
	Term string `json:"term,omitempty"`
}

type ListQuery struct {
	Pagination pagination.Query `json:"pagination"`
	Filter     FilterQuery      `json:"filter"`
	Sort       SortQuery        `json:"sort"`
	Search     SearchQuery      `json:"search"`
}

func Parse(values url.Values) (ListQuery, error) {
	page, err := positiveInt(values.Get("page"), 1)
	if err != nil {
		return ListQuery{}, fmt.Errorf("page: %w", err)
	}
	perPage, err := positiveInt(values.Get("per_page"), 20)
	if err != nil {
		return ListQuery{}, fmt.Errorf("per_page: %w", err)
	}
	if perPage > 100 {
		return ListQuery{}, fmt.Errorf("per_page must be at most 100")
	}

	filter := make(map[string]string)
	for key, value := range values {
		if strings.HasPrefix(key, "filter[") && strings.HasSuffix(key, "]") && len(value) > 0 {
			filter[strings.TrimSuffix(strings.TrimPrefix(key, "filter["), "]")] = value[0]
		}
	}

	desc := values.Get("sort_dir") == "desc"
	return ListQuery{
		Pagination: pagination.Query{Page: page, PerPage: perPage},
		Filter:     FilterQuery{Values: filter},
		Sort:       SortQuery{Field: values.Get("sort"), Desc: desc},
		Search:     SearchQuery{Term: strings.TrimSpace(values.Get("search"))},
	}, nil
}

func positiveInt(value string, fallback int) (int, error) {
	if value == "" {
		return fallback, nil
	}
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed < 1 {
		return 0, fmt.Errorf("must be a positive integer")
	}
	return parsed, nil
}
