package controllers

import (
	"strings"

	"github.com/goravel/framework/contracts/database/orm"
	"github.com/goravel/framework/contracts/http"

	adminactions "goravel/app/core/admin/actions"
	"goravel/app/core/resource"
	"goravel/app/facades"
)

const maxQuerySelectionIDs = 1000

func resolveSelectionIDs(ctx http.Context, manifest resource.Manifest, selection adminactions.Selection) ([]int64, error) {
	if selection.Mode == "ids" {
		return selection.IDs, nil
	}
	q, err := applySelectionQuery(ctx, facades.Orm().Query().Table(manifest.Table), manifest, selection.Query)
	if err != nil {
		return nil, err
	}
	if len(selection.ExcludeIDs) > 0 {
		values := make([]any, 0, len(selection.ExcludeIDs))
		for _, id := range selection.ExcludeIDs {
			values = append(values, id)
		}
		q = q.WhereNotIn("id", values)
	}
	var ids []int64
	if err := q.OrderBy("id", "asc").Limit(maxQuerySelectionIDs+1).Pluck("id", &ids); err != nil {
		return nil, err
	}
	if len(ids) > maxQuerySelectionIDs {
		return nil, adminactions.ErrTooManyIDs
	}
	return ids, nil
}

func applySelectionQuery(ctx http.Context, query orm.Query, manifest resource.Manifest, filters map[string]string) (orm.Query, error) {
	fieldPolicies, err := resourceFieldPoliciesFor(ctx, manifest, "view")
	if err != nil {
		return query, err
	}
	query, err = applyResourceScope(ctx, query, manifest, "view")
	if err != nil {
		return query, err
	}
	for key, value := range filters {
		value = strings.TrimSpace(value)
		if value == "" || value == "all" {
			continue
		}
		if key == "search" {
			query = applyResourceSearch(query, value, fieldNames(resourceSearchFieldsWithPolicies(manifest, fieldPolicies))...)
			continue
		}
		if key == "status" && manifest.Name == "users" {
			if normalizeUserStatusFilter(value) == "" {
				return query, adminactions.ErrSelectionContract
			}
			query = query.Where("status = ?", value)
			continue
		}
		allowed := false
		for _, field := range resourceFilterFieldsWithPolicies(manifest, fieldPolicies) {
			if field.Name == key {
				query, allowed = applyResourceFilter(query, field, value)
				break
			}
		}
		if !allowed {
			return query, adminactions.ErrSelectionContract
		}
	}
	return query, nil
}
